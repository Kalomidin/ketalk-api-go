package user_manager

import (
	"context"
	"ketalk-api/pkg/manager/user/repository"

	"github.com/google/uuid"
)

type userManager struct {
	repository repository.Repository
}

func NewUserManager(repository repository.Repository) UserManager {
	return &userManager{
		repository,
	}
}

func (m *userManager) GetUser(ctx context.Context, userID uuid.UUID) (*User, error) {
	user, err := m.repository.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Image:    user.Image,
		Password: user.Password,
	}, nil
}
