package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log/slog"

	"vincadrn.com/santuy/internal/model"
	"vincadrn.com/santuy/internal/repository"
)

type AccountService struct {
	repo repository.AccountRepository
}

func NewAccountService(repo repository.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

func (svc *AccountService) GetUserDetails(ctx context.Context, email string) (*model.User, error) {
	return svc.repo.GetUserByEmail(ctx, email)
}

func (svc *AccountService) SaveUser(ctx context.Context, user *model.User) error {
	return svc.repo.SaveUser(ctx, user)
}

func (svc *AccountService) ListGroupsByUser(ctx context.Context, user *model.User) (*[]model.GroupRole, error) {
	return svc.repo.ListGroupsByUser(ctx, user)
}

func (svc *AccountService) CreateGroupInvitation(ctx context.Context, group *model.Group, user *model.User) error {
	token, err := generateToken(16)
	if err != nil {
		slog.Error("cannot generate token: %v", err)
		return err
	}

	return svc.repo.CreateGroupInvite(ctx, group, token, user)
}

// User contains email, group is what will be assigned, groupRoles is the provided list of group the user is in
func (svc *AccountService) AssignUserToGroup(ctx context.Context, user *model.User, token string) (string, error) {
	return svc.repo.SetUserToGroup(ctx, user, token)
}

func generateToken(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
