package repository

import (
	"context"
	"fmt"
	"ketalk-api/common"
	"log"
	"time"

	"gorm.io/gorm"
)

type repository struct {
	*gorm.DB
}

func NewRepository(ctx context.Context, db *gorm.DB) Repository {
	return &repository{
		db,
	}
}

func (r *repository) AddRefreshToken(ctx context.Context, deviceRefreshToken *DeviceRefreshToken) error {
	res := r.Create(deviceRefreshToken)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected != 1 {
		return common.ErrMoreThanOneRowUpdated
	}
	return nil
}

func (r *repository) DeleteRefreshToken(ctx context.Context, deviceId string) error {
	var d DeviceRefreshToken
	resp := r.Model(&d).Where("device_id = ? and deleted_at is null", deviceId).Update("deleted_at", time.Now())
	if resp.Error != nil {
		return resp.Error
	}
	if resp.RowsAffected > 1 {
		fmt.Println("more than one refresh token is being deleted")
	} else {
		log.Printf("could not find any active refresh token for given device id %s\n", deviceId)
	}
	return nil
}

func (r *repository) GetRefreshToken(ctx context.Context, refreshToken string) (*DeviceRefreshToken, error) {
	var d DeviceRefreshToken
	resp := r.Model(&d).Where("refresh_token = ? and deleted_at is null", refreshToken).First(&d)

	if resp.Error != nil {
		return nil, resp.Error
	}
	return &d, nil
}

func (r *repository) Migrate() error {
	return r.AutoMigrate(&DeviceRefreshToken{})
}
