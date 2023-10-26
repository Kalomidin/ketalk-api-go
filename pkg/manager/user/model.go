package user_manager

import (
	"context"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	Username string
	Email    string
	Image    *string
}

type UpdateUserRequest struct {
	UserID uuid.UUID
	Name   *string
	Image  *string
}

type GetPresignedUrlRequest struct {
	UserID uuid.UUID
}

type GetPresignedUrlResponse struct {
	Url       string
	ImageName string
}

type UserManager interface {
	GetUser(ctx context.Context, userID uuid.UUID) (*User, error)
	Update(ctx context.Context, req UpdateUserRequest) (*User, error)
	GetPresignedUrl(ctx context.Context, req GetPresignedUrlRequest) (*GetPresignedUrlResponse, error)
}
