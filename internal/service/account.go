package service

import (
	"context"

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

// User contains email, group is what will be assigned, groupRoles is the provided list of group the user is in
func (svc *AccountService) AssignUserToGroup(ctx context.Context, user *model.User, group *model.Group) error {
	return svc.repo.SetUserToGroup(ctx, user, group)
}
