package georegion_manager

import (
	"context"
	"ketalk-api/common"
	"ketalk-api/pkg/manager/georegion/repository"
	"ketalk-api/pkg/manager/port"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var GlobalGeofence = repository.Geofence{
	ID:   uuid.UUID{},
	Name: "Globe Geofence",
}

type geofencePort struct {
	geofenceRepository repository.GeofenceRepository
}

func NewGeofencePort(geofenceRepository repository.GeofenceRepository) port.GeofencePort {
	return &geofencePort{
		geofenceRepository,
	}
}

func (p *geofencePort) GetGeofenceByLocation(ctx context.Context, location common.Location) (*port.Geofence, error) {
	geofence, err := p.geofenceRepository.FindGeofeceByLocation(ctx, location)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &port.Geofence{
				ID:   GlobalGeofence.ID,
				Name: GlobalGeofence.Name,
			}, nil
		}
		return nil, err
	}

	return &port.Geofence{
		ID:   geofence.ID,
		Name: geofence.Name,
	}, nil
}

func (p *geofencePort) GetGeofenceById(ctx context.Context, id uuid.UUID) (*port.Geofence, error) {
	emptyUuid := uuid.UUID{}
	if id == emptyUuid {
		return &port.Geofence{
			ID:   GlobalGeofence.ID,
			Name: GlobalGeofence.Name,
		}, nil
	}

	geofence, err := p.geofenceRepository.GetGeofenceByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &port.Geofence{
		ID:   geofence.ID,
		Name: geofence.Name,
	}, nil
}
