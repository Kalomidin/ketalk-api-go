package repository

import (
	"context"
	"ketalk-api/common"

	"github.com/google/uuid"
)

type Geofence struct {
	ID   uuid.UUID `gorm:"type:uuid;primary_key;"`
	Name string    `gorm:"type:varchar(255);"`
	Geom []byte    `gorm:"type:geometry(Polygon,4326)"`
	common.CreatedDeleted
}

type GeofenceRepository interface {
	GetGeofenceByID(ctx context.Context, id uuid.UUID) (*Geofence, error)
	FindGeofeceByLocation(ctx context.Context, location common.Location) (*Geofence, error)
	Migrate() error
}
