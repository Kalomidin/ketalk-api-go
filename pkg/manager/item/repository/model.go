package repository

import (
	"context"
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
	KaratID       uint32
	CategoryID    uint32
	GeofenceID    uint32
	common.CreatedUpdatedDeleted
}

type ItemRepository interface {
	AddItem(ctx context.Context, item *Item) error
	Update(ctx context.Context, item *Item) error
	GetItems(ctx context.Context, GeofenceID uint32, userID uuid.UUID) ([]Item, error)
	GetUserItems(ctx context.Context, userID uuid.UUID) ([]Item, error)
	GetItem(ctx context.Context, itemId uuid.UUID) (*Item, error)
	IncrementFavoriteCount(ctx context.Context, itemId uuid.UUID) error
	DecrementFavoriteCount(ctx context.Context, itemId uuid.UUID) error
	IncrementMessageCount(ctx context.Context, itemId uuid.UUID) error
	DecrementMessageCount(ctx context.Context, itemId uuid.UUID) error
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
	Migrate() error
}

type UserItemRepository interface {
	Insert(ctx context.Context, userItem *UserItem) error
	Update(ctx context.Context, userItem *UserItem) error
	GetUserItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID) (*UserItem, error)
	GetUserFavoriteItems(ctx context.Context, userID uuid.UUID) ([]Item, error)
	GetPurchasedItems(ctx context.Context, userID uuid.UUID) ([]Item, error)
	Migrate() error
}
