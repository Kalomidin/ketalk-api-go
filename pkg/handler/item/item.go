package item_handler

import (
	"context"
	"fmt"
	"ketalk-api/common"
	item_manager "ketalk-api/pkg/manager/item"

	"github.com/gin-gonic/gin"
)

type HttpHandler struct {
	handler    ItemHandler
	middleware common.Middleware
}

type handler struct {
	manager item_manager.ItemManager
}

func NewHandler(manager item_manager.ItemManager) ItemHandler {
	return &handler{
		manager,
	}
}

func NewHttpHandler(ctx context.Context, h ItemHandler, middleware common.Middleware) *HttpHandler {
	return &HttpHandler{
		h,
		middleware,
	}
}

func (c *HttpHandler) Init(ctx context.Context, router *gin.Engine) {
	routes := map[string]map[string]common.HandlerFunc{
		"POST": {
			"":              c.middleware.HandlerWithAuth(c.CreateItem),
			"/:id/favorite": c.middleware.HandlerWithAuth(c.FavoriteItem),
			"/:id/purchase": c.middleware.HandlerWithAuth(c.CreatePurchase),
		},
		"PUT": {
			"/image/upload":              c.middleware.HandlerWithAuth(c.UploadItemImages),
			"/:id":                       c.middleware.HandlerWithAuth(c.UpdateItem),
			"/:id/incrementConversation": c.middleware.HandlerWithAuth(c.IncrementConversationCount),
		},
		"GET": {
			"/all/:geofenceId": c.middleware.HandlerWithAuth(c.GetItems),
			"/:id":             c.middleware.HandlerWithAuth(c.GetItem),
			"/favorite":        c.middleware.HandlerWithAuth(c.GetFavoriteItems),
			"/purchase":        c.middleware.HandlerWithAuth(c.GetPurchasedItems),
			"/user":            c.middleware.HandlerWithAuth(c.GetUserItems),
			"/karats":          c.GetAllKarats,
			"/categories":      c.GetAllCategories,
			"/:id/similar":     c.GetSimilarItems,
			"/:id/buyer":       c.middleware.HandlerWithAuth(c.GetItemBuyers),
		},
	}
	for method, route := range routes {
		for r, h := range route {
			router.Handle(method, fmt.Sprintf("/item%s", r), common.GenericHandler(h))
		}
	}
	fmt.Println("initialized item handler")
}
