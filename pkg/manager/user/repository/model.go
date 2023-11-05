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

type UserGeofence struct {
	ID         uuid.UUID `gorm:"primaryKey;default:gen_random_uuid()"`
	UserID     uuid.UUID `gorm:"type:uuid;"`
	GeofenceID uuid.UUID
	common.CreatedDeleted
}

type UserGeofenceRepository interface {
	GetUserGeofence(ctx context.Context, userID uuid.UUID) (*UserGeofence, error)
	Create(ctx context.Context, userGeofence *UserGeofence) error
	Migrate() error
}
