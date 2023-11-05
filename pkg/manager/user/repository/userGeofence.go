package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userGeofenceRepository struct {
	db *gorm.DB
}

func NewUserGeofenceRepository(db *gorm.DB) UserGeofenceRepository {
	return &userGeofenceRepository{
		db,
	}
}

func (r *userGeofenceRepository) GetUserGeofence(ctx context.Context, userID uuid.UUID) (*UserGeofence, error) {
	var userGeofence UserGeofence
	resp := r.db.Where("user_id = ?", userID).First(&userGeofence)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return &userGeofence, nil
}

func (r *userGeofenceRepository) Create(ctx context.Context, userGeofence *UserGeofence) error {
	resp := r.db.Create(userGeofence)
	if resp.Error != nil {
		return resp.Error
	}
	return nil
}

func (r *userGeofenceRepository) Migrate() error {
	return r.db.AutoMigrate(&UserGeofence{})
}
