package repository

import (
	"context"
	"database/sql"
	"log/slog"
	"strconv"

	"vincadrn.com/santuy/internal/model"
)

type AccountRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	SaveUser(ctx context.Context, user *model.User) error
	ListGroupsByUser(ctx context.Context, user *model.User) (*[]model.GroupRole, error)
	SetUserToGroup(ctx context.Context, user *model.User, group *model.Group) error
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
		`SELECT id_user, nama, email FROM user_table WHERE email = $1;`,
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
		`INSERT INTO user_table (nama, email) VALUES ($1, $2) ON CONFLICT (email) DO UPDATE SET email = $2;`,
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
	defer tx.Rollback()

	if err != nil {
		slog.Error("Cannot start tx for listing groups by user", "user", user.Email)
		slog.Error(err.Error())

		return nil, err
	}

	groups := &[]model.GroupRole{}
	rows, err := tx.QueryContext(
		ctx,
		`SELECT gt.id_group, gt.nama, gu.role
		FROM group_table gt
		JOIN groupuser gu
			ON gu.id_group = gt.id_group
		JOIN user_table ut
			ON gu.id_user = ut.id_user
		WHERE ut.email = $1
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

	tx.Commit()

	return groups, nil
}

func (r *accountRepository) SetUserToGroup(ctx context.Context, user *model.User, group *model.Group) error {
	slog.Info("Attempting to set user to group", "user", user.Email, "group", group.Id)
	tx, err := r.db.Begin()
	defer tx.Rollback()

	if err != nil {
		slog.Error("Cannot start tx to set user to group", "user", user.Email, "group", group.Id)
		slog.Error(err.Error())

		return err
	}

	defaultRole := "member"

	numericUserId, err := strconv.Atoi(user.Id)
	if err != nil {
		slog.Error("user id is invalid", "user", user.Email, "userId", user.Id, "group", group.Id)
		slog.Error(err.Error())

		return err
	}

	res, err := r.db.ExecContext(
		ctx,
		`INSERT INTO groupuser
		(id_user, id_group, role)
		VALUES
		($1, $2, $3)
		ON CONFLICT DO NOTHING
		`,
		numericUserId, group.Id, defaultRole,
	)
	if err != nil {
		slog.Error("Cannot set user to group", "user", user.Email, "group", group.Id)
		slog.Error(err.Error())

		return err
	}

	affectedRows, _ := res.RowsAffected()
	slog.Info("Set user to group successful", "user", user.Email, "group", group.Id, "affected_rows", affectedRows)

	tx.Commit()

	return nil
}
