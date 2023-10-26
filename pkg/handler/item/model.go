package item_handler

import "github.com/gin-gonic/gin"

type ItemHandler interface {
	GetItems(ctx *gin.Context) (*GetItemsResponse, error)
	CreateItem(ctx *gin.Context, req CreateItemRequest) (*CreateItemResponse, error)
	UploadItemImages(ctx *gin.Context, r UploadItemImagesRequest) (*UploadItemImagesResponse, error)
	GetItem(ctx *gin.Context) (*Item, error)
	GetFavoriteItems(ctx *gin.Context) (*GetFavoriteItemsResponse, error)
	FavoriteItem(ctx *gin.Context, req FavoriteItemRequest) (*FavoriteItemResponse, error)
	GetUserItems(ctx *gin.Context) (*GetUserItemsResponse, error)
	GetPurchasedItems(ctx *gin.Context) (*GetPurchasedItemsResponse, error)
	UpdateItem(ctx *gin.Context, req UpdateItemRequest) (*UpdateItemResponse, error)
}
