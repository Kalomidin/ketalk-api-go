package user_manager

import (
	"context"
	"fmt"
	"ketalk-api/pkg/manager/port"
	"ketalk-api/pkg/manager/user/repository"
	"ketalk-api/storage"

	"github.com/google/uuid"
)

type userManager struct {
	repository             repository.Repository
	userGeofenceRepository repository.UserGeofenceRepository
	geofencePort           port.GeofencePort
	azureBlobStorage       storage.Storage
}

func NewUserManager(repository repository.Repository, userGeofenceRepository repository.UserGeofenceRepository, geofencePort port.GeofencePort, azureBlobStorage storage.Storage) UserManager {
	return &userManager{
		repository,
		userGeofenceRepository,
		geofencePort,
		azureBlobStorage,
	}
}

func (m *userManager) GetUser(ctx context.Context, userID uuid.UUID) (*User, error) {
	user, err := m.repository.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	var url *string
	if user.Image != nil {
		image := m.azureBlobStorage.GetUserImage(*user.Image)
		url = &image
	}
	userGeofence, err := m.userGeofenceRepository.GetUserGeofence(ctx, userID)
	if err != nil {
		return nil, err
	}
	geofence, err := m.geofencePort.GetGeofenceById(ctx, userGeofence.GeofenceID)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Image:    url,
		Geofence: Geofence{
			ID:   geofence.ID,
			Name: geofence.Name,
		},
	}, nil
}

func (m *userManager) Update(ctx context.Context, req UpdateUserRequest) (*User, error) {
	if req.Name == nil && req.Image == nil {
		return nil, fmt.Errorf("empty update request is not allowed")
	}
	user, err := m.repository.GetUser(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		user.Username = *req.Name
	}
	if req.Image != nil {
		user.Image = req.Image
	}
	if err := m.repository.UpdateUser(ctx, user); err != nil {
		return nil, err
	}
	userGeofence, err := m.userGeofenceRepository.GetUserGeofence(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	geofence, err := m.geofencePort.GetGeofenceById(ctx, userGeofence.GeofenceID)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Image:    user.Image,
		Geofence: Geofence{
			ID:   geofence.ID,
			Name: geofence.Name,
		},
	}, nil
}

func (m *userManager) GetPresignedUrl(ctx context.Context, req GetPresignedUrlRequest) (*GetPresignedUrlResponse, error) {
	blob := fmt.Sprintf("%s/%s", req.UserID, uuid.New())
	url, err := m.azureBlobStorage.GeneratePresignedUrlToUpload(ctx, blob, storage.ContainerProfiles)
	if err != nil {
		return nil, err
	}
	return &GetPresignedUrlResponse{
		Url:       url,
		ImageName: blob,
	}, nil
}
