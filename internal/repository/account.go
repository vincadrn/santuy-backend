package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"strconv"
	"time"

	"vincadrn.com/santuy/internal/model"
)

type AccountRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	SaveUser(ctx context.Context, user *model.User) error
	ListGroupsByUser(ctx context.Context, user *model.User) (*[]model.GroupRole, error)
	CreateGroupInvite(ctx context.Context, group *model.Group, token string, user *model.User) error
	SetUserToGroup(ctx context.Context, user *model.User, token string) (string, error)
}

type accountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	rows := r.db.QueryRowContext(
		ctx,
		`SELECT id, name, email FROM account.user_account WHERE email = $1;`,
		email,
	)

	user := &model.User{}
	err := rows.Scan(&user.Id, &user.Name, &user.Email)

	if err != nil {
		slog.Error(err.Error())
		return nil, nil
	}

	return user, nil
}

// TODO: Need to consider whether this `ON CONFLICT`
// thing should be there or not. Consider
// to error out if email already exists.
func (r *accountRepository) SaveUser(ctx context.Context, user *model.User) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO account.user_account (name, email) VALUES ($1, $2) ON CONFLICT (email) DO NOTHING`,
		user.Name, user.Email,
	)

	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}

func (r *accountRepository) ListGroupsByUser(ctx context.Context, user *model.User) (*[]model.GroupRole, error) {
	slog.Info("Attempting to list groups by user", "user", user.Email)
	tx, err := r.db.Begin()
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err != nil {
		slog.Error("Cannot start tx for listing groups by user", "user", user.Email)
		slog.Error(err.Error())

		return nil, err
	}

	groups := &[]model.GroupRole{}
	rows, err := tx.QueryContext(
		ctx,
		`SELECT vg.id, vg.name, uvg.role
		FROM travel.vacation_group vg
		JOIN travel.user_vacation_group uvg
			ON uvg.group_id = vg.id
		JOIN account.user_account ua
			ON uvg.user_id = ua.id
		WHERE ua.email = $1
		`,
		user.Email,
	)
	if err != nil {
		slog.Error("Cannot list groups", "user", user.Email)
		slog.Error(err.Error())

		return nil, err
	}

	for rows.Next() {
		if rows.Err() != nil {
			slog.Error("Cannot find next rows when listing groups", "user", user.Email, "group length", len(*groups))
			slog.Error(err.Error())

			return nil, err
		}

		group := model.GroupRole{}
		err = rows.Scan(&group.Id, &group.Name, &group.Role)
		if err != nil {
			slog.Error("Failed when try to scan db rows", "user", user.Email)
			slog.Error(err.Error())

			return nil, err
		}

		*groups = append(*groups, group)
	}

	err = tx.Commit()
	if err != nil {
		slog.Error("Cannot commit tx when listing groups")
		slog.Error(err.Error())

		return nil, err
	}

	return groups, nil
}

func (r *accountRepository) CreateGroupInvite(ctx context.Context, group *model.Group, token string, user *model.User) error {
	slog.Info("creating group invite", "group", group.Name, "user", user.Email)

	_, err := r.db.ExecContext(
		ctx,
		// TODO: Maybe not-hard-coded 1) expiration time, 2) max use?
		`INSERT INTO account.group_invite (group_id, token, expires_at, max_uses, created_by, updated_by)
		VALUES ($1, $2, NOW() + INTERVAL '24 hours', 5, $3, $3)`,
		group.Id, token, user.Id,
	)

	if err != nil {
		slog.Error("cannot create group invite", "group", group.Name, "user", user.Email)
		slog.Error(err.Error())
	}

	return nil
}

func (r *accountRepository) SetUserToGroup(ctx context.Context, user *model.User, token string) (string, error) {
	slog.Info("attempting to set user to group", "user", user.Email)
	tx, err := r.db.Begin()
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err != nil {
		slog.Error("cannot start tx to set user to group", "user", user.Email)
		slog.Error(err.Error())

		return "", err
	}

	// Check for group invite
	var inviteId int
	var groupId string
	var expiresAt time.Time
	var maxUses, usedCount int
	err = tx.QueryRowContext(
		ctx,
		`SELECT id, group_id, expires_at, max_uses, used_count
		FROM account.group_invite
		WHERE token = $1
		FOR UPDATE
		`,
		token,
	).Scan(&inviteId, &groupId, &expiresAt, &maxUses, &usedCount)

	if err != nil {
		slog.Error("cannot get invite from token", "error", err)
		return "", errors.New("invalid token")
	}

	now := time.Now()
	if now.After(expiresAt) {
		slog.Error("expired token", "now", now, "expires", expiresAt)
		return "", errors.New("token expired")
	}

	if usedCount >= maxUses {
		slog.Error("token usage exceeded", "usage", usedCount, "max", maxUses)
		return "", errors.New("token usage exceeded")
	}

	// Add user to group
	defaultRole := "member"

	numericUserId, err := strconv.Atoi(user.Id)
	if err != nil {
		slog.Error("user id is invalid", "user", user.Email, "userId", user.Id)
		slog.Error(err.Error())

		return "", err
	}

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO travel.user_vacation_group
		(user_id, group_id, role)
		VALUES
		($1, $2, $3)
		ON CONFLICT DO NOTHING
		`,
		numericUserId, groupId, defaultRole,
	)
	if err != nil {
		slog.Error("cannot set user to group", "user", user.Email, "group", groupId)
		slog.Error(err.Error())

		return "", err
	}

	// Increment token usage of the group invite
	_, err = tx.ExecContext(
		ctx,
		`UPDATE account.group_invite SET used_count = used_count + 1
		WHERE id = $1`,
		inviteId,
	)
	if err != nil {
		slog.Error("cannot increment token usage")
		slog.Error(err.Error())

		return "", err
	}

	// Finally commit
	err = tx.Commit()

	if err != nil {
		slog.Error("Cannot commit tx when setting user to group")
		slog.Error(err.Error())

		return "", err
	}

	return groupId, nil
}
