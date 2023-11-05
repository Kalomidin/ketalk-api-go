package item_manager

import (
	"context"
	"fmt"
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
	KaratID     uuid.UUID
	CategoryID  uuid.UUID
	GeofenceID  string
	Images      []string
	Thumbnail   string
	Location    common.Location
}

type AddItemResponse struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	PresignedUrls []ImageUploadUrlWithName
}

type ImageUploadUrlWithName struct {
	ID        uuid.UUID
	SignedUrl string
	Name      string
}

type ItemImage struct {
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
	Images         []ItemImage
	CreatedAt      time.Time
	Negotiable     bool
	IsHidden       bool
	IsUserFavorite bool
	KaratID        uuid.UUID
	CategoryID     uuid.UUID
	Weight         float32
	Size           float32
}

type ItemOwner struct {
	ID       uuid.UUID
	Name     string
	Avatar   *string
	Geofence Geofence
}

type Geofence struct {
	ID   uuid.UUID
	Name string
}

type GetItemsRequest struct {
	GeofenceID string
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

type UpdateItemRequest struct {
	UserID      uuid.UUID
	ItemID      uuid.UUID
	IsHidden    *bool
	ItemStatus  *ItemStatus
	Title       *string
	Description *string
	Price       *uint32
	Negotiable  *bool
	Size        *float32
	Weight      *float32
	KaratID     *uuid.UUID
	CategoryID  *uuid.UUID
	Images      []UpdatedItemImage
}

type UpdatedItemImage struct {
	ID      *uuid.UUID
	Name    *string
	IsCover bool
}

type UpdateItemResponse struct {
	NewImagesPresignedUrls []ImageUploadUrlWithName
}

type IncrementConversationCountRequest struct {
	ItemID uuid.UUID
	UserID uuid.UUID
}

type Karat struct {
	ID      uuid.UUID
	Name    string
	Locales map[string]KaratLocale
}

type KaratLocale struct {
	Description string
	Name        string
}

type Category struct {
	ID      uuid.UUID
	Name    string
	Locales map[string]CategoryLocale
}

type CategoryLocale struct {
	Description string
	Name        string
}

type GetItemResponse struct {
	Item
}

type GetSimilarItemsRequest struct {
	ItemID uuid.UUID
	UserID uuid.UUID
}

type GetSimilarItemsResponse struct {
	SuggestedItems []ItemBlock
	OtherUserItems []ItemBlock
}

type ItemBuyer struct {
	ID             uuid.UUID
	Name           string
	Avatar         *string
	LastMessagedAt time.Time
}

type GetItemBuyersRequest struct {
	ItemID uuid.UUID
}

type CreatePurchaseRequest struct {
	ItemID  uuid.UUID
	BuyerID uuid.UUID
}

type CreatePurchaseResponse struct {
	ItemID  uuid.UUID
	BuyerID uuid.UUID
}

type ItemManager interface {
	AddItem(ctx context.Context, item AddItemRequest) (*AddItemResponse, error)
	UploadItemImages(ctx context.Context, req UploadItemImagesRequest) (*UploadItemImagesResponse, error)
	GetItems(ctx context.Context, req GetItemsRequest) ([]ItemBlock, error)
	GetItem(ctx context.Context, req GetItemRequest) (*GetItemResponse, error)
	GetFavoriteItems(ctx context.Context, req GetFavoriteItemsRequest) ([]ItemBlock, error)
	GetUserItems(ctx context.Context, req GetUserItemsRequest) ([]ItemBlock, error)
	FavoriteItem(ctx context.Context, req FavoriteItemRequest) (*FavoriteItemResponse, error)
	IncrementConversationCount(ctx context.Context, req IncrementConversationCountRequest) error
	GetPurchasedItems(ctx context.Context, req GetPurchasedItemsRequest) ([]ItemBlock, error)
	UpdateItem(ctx context.Context, req UpdateItemRequest) (*UpdateItemResponse, error)
	GetAllKarats(ctx context.Context) ([]Karat, error)
	GetAllCategories(ctx context.Context) ([]Category, error)
	GetSimilarItems(ctx context.Context, req GetSimilarItemsRequest) (*GetSimilarItemsResponse, error)
	GetItemBuyers(ctx context.Context, req GetItemBuyersRequest) ([]ItemBuyer, error)
	CreatePurchase(ctx context.Context, req CreatePurchaseRequest) (*CreatePurchaseResponse, error)
}

type ItemStatus string

const (
	ItemStatusActive   ItemStatus = "Active"
	ItemStatusSold     ItemStatus = "Sold"
	ItemStatusReserved ItemStatus = "Reserved"
)

var ErrInvalidItemStatus = fmt.Errorf("invalid item status")

func ParseItemStatus(itemStatus string) (*ItemStatus, error) {
	switch ItemStatus(itemStatus) {
	case ItemStatusActive, ItemStatusReserved, ItemStatusSold:
		itemStatus := ItemStatus(itemStatus)
		return &itemStatus, nil
	default:
		return nil, ErrInvalidItemStatus
	}
}
