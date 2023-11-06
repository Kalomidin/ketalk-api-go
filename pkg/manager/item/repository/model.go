package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"ketalk-api/common"

	"github.com/google/uuid"
)

type Item struct {
	ID            uuid.UUID `gorm:"primaryKey;default:gen_random_uuid()"`
	Title         string
	Description   string
	Price         uint32
	Negotiable    bool
	OwnerID       uuid.UUID
	ItemStatus    string
	IsHidden      bool
	FavoriteCount uint32
	MessageCount  uint32
	SeenCount     uint32
	Size          float32
	Weight        float32
	KaratID       uuid.UUID
	CategoryID    uuid.UUID
	GeofenceID    string
	common.CreatedUpdatedDeleted
}

type ItemRepository interface {
	AddItem(ctx context.Context, item *Item) error
	Update(ctx context.Context, item *Item) error
	GetItems(ctx context.Context, GeofenceID string, userID uuid.UUID) ([]Item, error)
	GetUserItems(ctx context.Context, userID uuid.UUID) ([]Item, error)
	GetItem(ctx context.Context, itemId uuid.UUID) (*Item, error)
	IncrementFavoriteCount(ctx context.Context, itemId uuid.UUID) error
	DecrementFavoriteCount(ctx context.Context, itemId uuid.UUID) error
	IncrementMessageCount(ctx context.Context, itemId uuid.UUID) error
	DecrementMessageCount(ctx context.Context, itemId uuid.UUID) error
	GetLimitedUserItems(ctx context.Context, userID uuid.UUID, limit int) ([]Item, error)
	GetLimitedItemsByCategoryOrKarat(ctx context.Context, userIDToExlude uuid.UUID, categoryID uuid.UUID, karatID uuid.UUID, limit int) ([]Item, error)
	SearchItems(ctx context.Context, keyword string, priceRange []uint32, sizeRange []float32, karatIds []uuid.UUID, categoryIds []uuid.UUID) ([]Item, error)
	Migrate() error
}

type ItemImage struct {
	ID              uuid.UUID `gorm:"primaryKey;default:gen_random_uuid()"`
	Key             string
	ItemID          uuid.UUID
	IsCover         bool
	UploadedToCloud bool
	common.CreatedUpdatedDeleted
}

type UserItem struct {
	ID          uuid.UUID `gorm:"primaryKey;default:gen_random_uuid()"`
	UserID      uuid.UUID
	ItemID      uuid.UUID
	IsFavorite  bool
	IsPurchased bool
	common.CreatedUpdatedDeleted
}

type ItemImageRepository interface {
	GetItemImages(ctx context.Context, itemID uuid.UUID) ([]ItemImage, error)
	AddItemImages(ctx context.Context, itemID uuid.UUID, images []ItemImage) error
	GetItemThumbnail(ctx context.Context, itemID uuid.UUID) (ItemImage, error)
	UpdateItemImagesToUploaded(ctx context.Context, itemID uuid.UUID, imageIds []uuid.UUID) error
	UpdateItemImage(ctx context.Context, itemID uuid.UUID, imageId uuid.UUID, isCover bool) error
	DeleteItemImages(ctx context.Context, itemID uuid.UUID, imageIds []uuid.UUID) error
	Migrate() error
}

type UserItemRepository interface {
	Insert(ctx context.Context, userItem *UserItem) error
	Update(ctx context.Context, userItem *UserItem) error
	GetUserItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID) (*UserItem, error)
	GetUserFavoriteItems(ctx context.Context, userID uuid.UUID) ([]Item, error)
	PurchaseItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID) error
	GetPurchasedItems(ctx context.Context, userID uuid.UUID) ([]Item, error)
	GetItemBuyer(ctx context.Context, itemID uuid.UUID) ([]UserItem, error)
	Migrate() error
}

type Karat struct {
	ID      uuid.UUID `gorm:"primaryKey;default:gen_random_uuid()"`
	Name    string
	Locales KaratLocales `gorm:"type:json"`
	common.CreatedDeleted
}

type KaratLocales map[string]KaratLocale

type KaratLocale struct {
	Description string
	Name        string
}

func (kl *KaratLocales) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	err := json.Unmarshal(bytes, kl)
	if err != nil {
		return err
	}
	return nil
}

type KaratRepository interface {
	GetAllKarats(ctx context.Context) ([]Karat, error)
	Migrate() error
}

type Category struct {
	ID      uuid.UUID `gorm:"primaryKey;default:gen_random_uuid()"`
	Name    string
	Locales CategoryLocales `gorm:"type:json"`
	common.CreatedDeleted
}

type CategoryLocales map[string]CategoryLocale

type CategoryLocale struct {
	Description string
	Name        string
}

func (kl *CategoryLocales) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("Failed to unmarshal JSONB value: %v", value)
	}

	err := json.Unmarshal(bytes, kl)
	if err != nil {
		return err
	}
	return nil
}

type CategoryRepository interface {
	GetAllCategories(ctx context.Context) ([]Category, error)
	Migrate() error
}
