package item_manager

import (
	"context"
	"ketalk-api/common"
	"time"

	"github.com/google/uuid"
)

type AddItemRequest struct {
	Title       string
	Description string
	Price       uint32
	Size        float32
	Weight      float32
	OwnerID     uuid.UUID
	Negotiable  bool
	KaratID     uint32
	CategoryID  uint32
	GeofenceID  uint32
	Images      []string
	Thumbnail   string
}

type AddItemResponse struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	PresignedUrls []SignedUrlWithImageID
}

type SignedUrlWithImageID struct {
	ID        uuid.UUID
	SignedUrl string
	Name      string
}

type ItemBlock struct {
	ID            uuid.UUID
	Title         string
	Description   string
	Price         uint32
	OwnerID       uuid.UUID
	FavoriteCount uint32
	MessageCount  uint32
	SeenCount     uint32
	ItemStatus    ItemStatus
	IsHidden      bool
	Thumbnail     string
	CreatedAt     time.Time
}

type Item struct {
	ID             uuid.UUID
	Title          string
	Description    string
	Price          uint32
	Owner          ItemOwner
	FavoriteCount  uint32
	MessageCount   uint32
	SeenCount      uint32
	ItemStatus     ItemStatus
	Thumbnail      string
	Images         []string
	CreatedAt      time.Time
	Negotiable     bool
	IsHidden       bool
	IsUserFavorite bool
}

type ItemOwner struct {
	ID       uuid.UUID
	Name     string
	Avatar   *string
	Location *common.Location
}

type GetItemsRequest struct {
	GeofenceID uint32
	UserID     uuid.UUID
}

type UploadItemImagesRequest struct {
	ItemID   uuid.UUID
	ImageIds []uuid.UUID
}

type UploadItemImagesResponse struct {
}

type GetItemRequest struct {
	ItemID uuid.UUID
	UserID uuid.UUID
}

type GetFavoriteItemsRequest struct {
	UserID uuid.UUID
}
type FavoriteItemRequest struct {
	UserID     uuid.UUID
	ItemID     uuid.UUID
	IsFavorite bool
}

type FavoriteItemResponse struct {
}

type GetUserItemsRequest struct {
	UserID uuid.UUID
}

type GetPurchasedItemsRequest struct {
	UserID uuid.UUID
}

type ItemManager interface {
	AddItem(ctx context.Context, item AddItemRequest) (*AddItemResponse, error)
	UploadItemImages(ctx context.Context, req UploadItemImagesRequest) (*UploadItemImagesResponse, error)
	GetItems(ctx context.Context, req GetItemsRequest) ([]ItemBlock, error)
	GetItem(ctx context.Context, req GetItemRequest) (*Item, error)
	GetFavoriteItems(ctx context.Context, req GetFavoriteItemsRequest) ([]ItemBlock, error)
	GetUserItems(ctx context.Context, req GetUserItemsRequest) ([]ItemBlock, error)
	FavoriteItem(ctx context.Context, req FavoriteItemRequest) (*FavoriteItemResponse, error)
	GetPurchasedItems(ctx context.Context, req GetPurchasedItemsRequest) ([]ItemBlock, error)
}

type ItemStatus string

const (
	ItemStatusActive   ItemStatus = "Active"
	ItemStatusSold     ItemStatus = "Sold"
	ItemStatusReserved ItemStatus = "Reserved"
)
