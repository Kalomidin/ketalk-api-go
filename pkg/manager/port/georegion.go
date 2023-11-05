package port

import (
	"context"
	"ketalk-api/common"

	"github.com/google/uuid"
)

type Geofence struct {
	ID   uuid.UUID
	Name string
}

type GeofencePort interface {
	GetGeofenceByLocation(ctx context.Context, location common.Location) (*Geofence, error)
	GetGeofenceById(ctx context.Context, id uuid.UUID) (*Geofence, error)
}
