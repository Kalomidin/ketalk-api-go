package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type itemImageRepository struct {
	*gorm.DB
}

func NewItemImageRepository(ctx context.Context, db *gorm.DB) ItemImageRepository {
	return &itemImageRepository{
		db,
	}
}

func (r *itemImageRepository) AddItemImages(ctx context.Context, itemID uuid.UUID, images []ItemImage) error {
	res := r.CreateInBatches(&images, len(images))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected != int64(len(images)) {
		return fmt.Errorf("unexpected number of rows affected")
	}
	return nil
}

func (r *itemImageRepository) GetItemImages(ctx context.Context, itemID uuid.UUID) ([]ItemImage, error) {
	var images []ItemImage = make([]ItemImage, 0)
	resp := r.Where("item_id = ?", itemID).Find(&images)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return images, nil
}

func (r *itemImageRepository) GetItemThumbnail(ctx context.Context, itemID uuid.UUID) (ItemImage, error) {
	var image ItemImage
	resp := r.Where("item_id = ? AND is_cover = ?", itemID, true).First(&image)
	if resp.Error != nil {
		return ItemImage{}, resp.Error
	}
	return image, nil
}

func (r *itemImageRepository) UpdateItemImagesToUploaded(ctx context.Context, itemID uuid.UUID, imageKeys []uuid.UUID) error {
	res := r.Model(&ItemImage{}).Where("item_id = ? AND id IN ?", itemID, imageKeys).Update("uploaded_to_cloud", true)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected != int64(len(imageKeys)) {
		return fmt.Errorf("unexpected number of rows affected")
	}
	return nil
}

func (r *itemImageRepository) Migrate() error {
	return r.AutoMigrate(&ItemImage{})
}
