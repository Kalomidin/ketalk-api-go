package item_handler

import "github.com/gin-gonic/gin"

type ItemHandler interface {
	GetItems(ctx *gin.Context) (*GetItemsResponse, error)
	CreateItem(ctx *gin.Context, req CreateItemRequest) (*CreateItemResponse, error)
	UploadItemImages(ctx *gin.Context, r UploadItemImagesRequest) (*UploadItemImagesResponse, error)
}
