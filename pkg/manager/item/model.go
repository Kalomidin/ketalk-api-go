package item_manager

import (
	"context"
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

type Item struct {
	ID            uuid.UUID
	Title         string
	Description   string
	Price         uint32
	OwnerID       uuid.UUID
	FavoriteCount uint32
	MessageCount  uint32
	SeenCount     uint32
	ItemStatus    ItemStatus
	Thumbnail     string
	CreatedAt     time.Time
}

type GetItemsRequest struct {
	GeofenceID uint32
}

type UploadItemImagesRequest struct {
	ItemID   uuid.UUID
	ImageIds []uuid.UUID
}

type UploadItemImagesResponse struct {
}

type ItemManager interface {
	AddItem(ctx context.Context, item AddItemRequest) (*AddItemResponse, error)
	UploadItemImages(ctx context.Context, r UploadItemImagesRequest) (*UploadItemImagesResponse, error)
	GetItems(ctx context.Context, req GetItemsRequest) ([]Item, error)
}

type ItemStatus string

const (
	ItemStatusActive   ItemStatus = "Active"
	ItemStatusSold     ItemStatus = "Sold"
	ItemStatusReserved ItemStatus = "Reserved"
)
