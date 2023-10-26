package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type itemRepository struct {
	*gorm.DB
}

func NewItemRepository(ctx context.Context, db *gorm.DB) ItemRepository {
	return &itemRepository{
		db,
	}
}

func (r *itemRepository) AddItem(ctx context.Context, item *Item) error {
	res := r.Create(item)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected != 1 {
		return fmt.Errorf("more than one row updated")
	}
	return nil
}

func (r *itemRepository) GetItems(ctx context.Context, GeofenceID uint32, userID uuid.UUID) ([]Item, error) {
	var items []Item = make([]Item, 0)
	resp := r.Where("geofence_id = ? AND owner_id != ?", GeofenceID, userID).Find(&items)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return items, nil
}

func (r *itemRepository) GetItem(ctx context.Context, itemID uuid.UUID) (*Item, error) {
	var item Item
	resp := r.Where("id = ?", itemID).First(&item)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return &item, nil
}

func (r *itemRepository) GetUserItems(ctx context.Context, userID uuid.UUID) ([]Item, error) {
	var items []Item = make([]Item, 0)
	resp := r.Where("owner_id = ?", userID).Find(&items)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return items, nil
}

func (r *itemRepository) Migrate() error {
	return r.AutoMigrate(&Item{})
}
