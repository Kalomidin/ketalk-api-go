package user_manager

import (
	"context"
	"fmt"
	"ketalk-api/pkg/manager/user/repository"
	"ketalk-api/storage"

	"github.com/google/uuid"
)

type userManager struct {
	repository       repository.Repository
	azureBlobStorage storage.AzureBlobStorage
}

func NewUserManager(repository repository.Repository, azureBlobStorage storage.AzureBlobStorage) UserManager {
	return &userManager{
		repository,
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
		imageUrl := m.azureBlobStorage.GetFrontDoorUrl(*user.Image)
		url = &imageUrl
	}
	return &User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Image:    url,
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
	return &User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Image:    user.Image,
	}, nil
}

func (m *userManager) GetPresignedUrl(ctx context.Context, req GetPresignedUrlRequest) (*GetPresignedUrlResponse, error) {
	blob := fmt.Sprintf("%s/%s", req.UserID, uuid.New())
	url, err := m.azureBlobStorage.GeneratePresignedUrlToUpload(blob)
	if err != nil {
		return nil, err
	}
	return &GetPresignedUrlResponse{
		Url:       url,
		ImageName: blob,
	}, nil
}
