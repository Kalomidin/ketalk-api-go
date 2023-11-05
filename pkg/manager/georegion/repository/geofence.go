package repository

import (
	"context"
	"fmt"
	"ketalk-api/common"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type geofenceRepository struct {
	*gorm.DB
}

func NewGeofenceRepository(db *gorm.DB) GeofenceRepository {
	return &geofenceRepository{
		db,
	}
}

func (r *geofenceRepository) GetGeofenceByID(ctx context.Context, id uuid.UUID) (*Geofence, error) {
	var geofence Geofence
	resp := r.Where("id = ?", id).First(&geofence)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return &geofence, nil
}

func (r *geofenceRepository) FindGeofeceByLocation(ctx context.Context, location common.Location) (*Geofence, error) {
	var geofence Geofence
	query := fmt.Sprintf("ST_Contains(geom, ST_GeomFromText('POINT(%f %f)', 4326))", location.Longitude, location.Latitude)
	resp := r.Where(query).First(&geofence)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return &geofence, nil
}

func (r *geofenceRepository) Migrate() error {
	return r.AutoMigrate(&Geofence{})
	// return nil
}
