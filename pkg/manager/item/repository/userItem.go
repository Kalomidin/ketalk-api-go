package repository

import (
	"context"
	"fmt"
	"ketalk-api/common"
	"ketalk-api/pkg/config"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userItemRepository struct {
	*gorm.DB
	dbConfig config.Postgres
}

func NewUserItemRepository(db *gorm.DB, dbConfig config.Postgres) UserItemRepository {
	return &userItemRepository{
		db,
		dbConfig,
	}
}

func (r *userItemRepository) Insert(ctx context.Context, userItem *UserItem) error {
	res := r.Create(userItem)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected != 1 {
		return common.ErrMoreThanOneRowUpdated
	}
	return nil
}

func (r *userItemRepository) Update(ctx context.Context, userItem *UserItem) error {
	res := r.Model(userItem).Where("id = ?", userItem.ID).Updates(map[string]interface{}{
		"is_favorite":  userItem.IsFavorite,
		"is_purchased": userItem.IsPurchased,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected != 1 {
		return common.ErrMoreThanOneRowUpdated
	}
	return nil
}

func (r *userItemRepository) GetUserFavoriteItems(ctx context.Context, userID uuid.UUID) ([]Item, error) {
	var userItems []Item = make([]Item, 0)
	resp := r.Model(&Item{}).
		InnerJoins(fmt.Sprintf("INNER JOIN %s.%s on user_item.item_id = item.id", r.dbConfig.GetSchema(), "user_item")).
		Where("user_item.user_id = ? and is_favorite = ?", userID, true).
		Find(&userItems)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return userItems, nil
}

func (r *userItemRepository) GetPurchasedItems(ctx context.Context, userID uuid.UUID) ([]Item, error) {
	var userItems []Item = make([]Item, 0)
	resp := r.Model(&Item{}).
		InnerJoins(fmt.Sprintf("INNER JOIN %s.%s on user_item.item_id = item.id", r.dbConfig.GetSchema(), "user_item")).
		Where("user_item.user_id = ? and is_purchased = ?", userID, true).
		Find(&userItems)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return userItems, nil

}

func (r *userItemRepository) GetUserItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID) (*UserItem, error) {
	var userItem UserItem
	resp := r.Model(&UserItem{}).
		Where("user_id = ? and item_id = ?", userID, itemID).
		First(&userItem)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return &userItem, nil
}

func (r *userItemRepository) PurchaseItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID) error {
	res := r.Model(&UserItem{}).
		Where("user_id = ? and item_id = ?", userID, itemID).
		Update("is_purchased", true)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected != 1 {
		return common.ErrMoreThanOneRowUpdated
	}
	return nil
}

func (r *userItemRepository) GetItemBuyer(ctx context.Context, itemID uuid.UUID) ([]UserItem, error) {
	var userItems []UserItem = make([]UserItem, 0)
	resp := r.Model(&UserItem{}).
		Where("item_id = ? and is_purchased = ?", itemID, true).
		Find(&userItems)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return userItems, nil
}

func (r *userItemRepository) Migrate() error {
	return r.AutoMigrate(&UserItem{})
}
