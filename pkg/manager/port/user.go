package port

import (
	"context"

	"github.com/google/uuid"
)

type CreateOrGetUserRequest struct {
	Username   string
	Email      string
	Image      *string
	Password   *string
	GeofenceID uuid.UUID
}

type User struct {
	ID         uuid.UUID
	Username   string
	Email      string
	Image      *string
	Password   *string
	GeofenceID uuid.UUID
}

type UserPort interface {
	CreateOrGetUser(ctx context.Context, req CreateOrGetUserRequest) (*User, error)
	GetUser(ctx context.Context, userId uuid.UUID) (*User, error)
}
