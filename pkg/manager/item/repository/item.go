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

func (r *itemRepository) Update(ctx context.Context, item *Item) error {
	res := r.Model(item).Where("id = ?", item.ID).Updates(map[string]interface{}{
		"title":       item.Title,
		"description": item.Description,
		"price":       item.Price,
		"negotiable":  item.Negotiable,
		"item_status": item.ItemStatus,
		"is_hidden":   item.IsHidden,
		"size":        item.Size,
		"weight":      item.Weight,
		"karat_id":    item.KaratID,
		"category_id": item.CategoryID,
		"geofence_id": item.GeofenceID,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected != 1 {
		return fmt.Errorf("more than one row updated")
	}
	return nil
}

func (r *itemRepository) GetItems(ctx context.Context, GeofenceID uuid.UUID, userID uuid.UUID) ([]Item, error) {
	var items []Item = make([]Item, 0)
	resp := r.Where("geofence_id = ? AND owner_id != ? and is_hidden = false", GeofenceID, userID).Find(&items)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return items, nil
}

func (r *itemRepository) SearchItems(ctx context.Context, keyword string, priceRange []uint32, sizeRange []float32, karatIds []uuid.UUID, categoryIds []uuid.UUID) ([]Item, error) {
	var items []Item = make([]Item, 0)
	query := r.Where("price BETWEEN ? AND ? AND size BETWEEN ? AND ?", priceRange[0], priceRange[1], sizeRange[0], sizeRange[1])
	if keyword != "" {
		query = query.Where("title LIKE ?", fmt.Sprintf("%%%s%%", keyword))
	}
	if len(karatIds) > 0 {
		query = query.Where("karat_id IN ?", karatIds)
	}
	if len(categoryIds) > 0 {
		query = query.Where("category_id IN ?", categoryIds)
	}
	resp := query.Find(&items)
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

func (r *itemRepository) GetLimitedUserItems(ctx context.Context, userID uuid.UUID, limit int) ([]Item, error) {
	var items []Item = make([]Item, 0)
	resp := r.Where("owner_id = ?", userID).Limit(limit).Find(&items)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return items, nil
}

func (r *itemRepository) GetLimitedItemsByCategoryOrKarat(ctx context.Context, userIDToExlude uuid.UUID, categoryID uuid.UUID, karatID uuid.UUID, limit int) ([]Item, error) {
	var items []Item = make([]Item, 0)
	resp := r.Where("(category_id = ? OR karat_id = ?) AND owner_id != ?", categoryID, karatID, userIDToExlude).Limit(limit).Find(&items)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return items, nil
}

func (r *itemRepository) IncrementFavoriteCount(ctx context.Context, itemID uuid.UUID) error {
	return r.Model(&Item{}).Where("id = ?", itemID).Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error
}

func (r *itemRepository) DecrementFavoriteCount(ctx context.Context, itemId uuid.UUID) error {
	return r.Model(&Item{}).Where("id = ? AND favorite_count > 0", itemId).Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error
}

func (r *itemRepository) IncrementMessageCount(ctx context.Context, itemId uuid.UUID) error {
	return r.Model(&Item{}).Where("id = ?", itemId).Update("message_count", gorm.Expr("message_count + ?", 1)).Error
}

func (r *itemRepository) DecrementMessageCount(ctx context.Context, itemId uuid.UUID) error {
	return r.Model(&Item{}).Where("id = ? AND message_count > 0", itemId).Update("message_count", gorm.Expr("message_count - ?", 1)).Error
}

func (r *itemRepository) DeleteItem(ctx context.Context, itemId uuid.UUID) error {
	resp := r.Delete(&Item{
		ID: itemId,
	}, itemId)
	if resp.RowsAffected == 0 {
		return fmt.Errorf("item not found")
	}
	return resp.Error
}

func (r *itemRepository) Migrate() error {
	return r.AutoMigrate(&Item{})
}
