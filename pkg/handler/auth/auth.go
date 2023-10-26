package auth_handler

import (
	"context"
	"fmt"
	"ketalk-api/common"
	manager "ketalk-api/pkg/manager/auth"

	"github.com/gin-gonic/gin"
)

type HttpHandler struct {
	handler    AuthHandler
	middleware common.Middleware
}

type handler struct {
	service manager.AuthManager
}

func NewHandler(service manager.AuthManager) AuthHandler {
	return &handler{
		service,
	}
}

func NewHttpHandler(ctx context.Context, h AuthHandler, middleware common.Middleware) *HttpHandler {
	return &HttpHandler{
		h,
		middleware,
	}
}

func (c *HttpHandler) Init(ctx context.Context, router *gin.Engine) {
	routes := map[string]map[string]common.HandlerFunc{
		"POST": {
			"/signup-or-login": c.SignupOrLogin,
		},
	}
	for method, route := range routes {
		for r, h := range route {
			router.Handle(method, fmt.Sprintf("/auth%s", r), common.GenericHandler(h))
		}
	}
	fmt.Println("initialized auth handler")
}
