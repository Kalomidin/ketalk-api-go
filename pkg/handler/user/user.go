package user_handler

import (
	"context"
	"fmt"
	"ketalk-api/common"
	manager "ketalk-api/pkg/manager/user"

	"github.com/gin-gonic/gin"
)

type HttpHandler struct {
	handler    UserHandler
	middleware common.Middleware
}

type handler struct {
	service manager.UserManager
}

func NewHandler(service manager.UserManager) UserHandler {
	return &handler{
		service,
	}
}

func NewHttpHandler(ctx context.Context, h UserHandler, middleware common.Middleware) *HttpHandler {
	return &HttpHandler{
		h,
		middleware,
	}
}

func (c *HttpHandler) Init(ctx context.Context, router *gin.Engine) {
	routes := map[string]map[string]common.HandlerFunc{
		"GET": {
			"/user": c.middleware.HandlerWithAuth(c.GetUser),
		},
	}
	for method, route := range routes {
		for r, h := range route {
			router.Handle(method, r, common.GenericHandler(h))
		}
	}
	fmt.Println("initialized user handler")
}
