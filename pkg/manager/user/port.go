package user_manager

import (
	"context"
	"ketalk-api/pkg/manager/port"
	"ketalk-api/pkg/manager/user/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userPort struct {
	userRepository repository.Repository
}

func NewUserPort(userRepository repository.Repository) port.UserPort {
	return &userPort{
		userRepository,
	}
}

func (p *userPort) CreateOrGetUser(ctx context.Context, req port.CreateOrGetUserRequest) (*port.User, error) {
	// get user by email
	user, err := p.userRepository.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return &port.User{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Password: user.Password,
			Image:    user.Image,
		}, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// else create user
	user = &repository.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Image:    req.Image,
	}
	if err = p.userRepository.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	return &port.User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Image:    user.Image,
	}, nil
}

func (p *userPort) GetUser(ctx context.Context, userId uuid.UUID) (*port.User, error) {
	user, err := p.userRepository.GetUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	return &port.User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Image:    user.Image,
	}, nil
}
