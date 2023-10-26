package repository

import (
	"context"
	"ketalk-api/common"
	"time"

	"github.com/google/uuid"
)

type DeviceRefreshToken struct {
	RefreshToken         string
	DeviceID             string
	DeviceOS             string
	UserId               uuid.UUID
	RefreshTokenExpiryAt time.Time
	common.CreatedDeleted
}

type Repository interface {
	AddRefreshToken(ctx context.Context, deviceRefreshToken *DeviceRefreshToken) error
	DeleteRefreshToken(ctx context.Context, deviceId string) error
	GetRefreshToken(ctx context.Context, refreshToken string) (*DeviceRefreshToken, error)
	Migrate() error
}
