package repository

import (
	"context"
	"ketalk-api/common"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"primaryKey;default:gen_random_uuid()"`
	Username string
	Image    *string
	Email    string
	Password *string
	common.CreatedUpdatedDeleted
}

type Repository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, userId uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	MigrateUser() error
}
