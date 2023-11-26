package auth_handler

import (
	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	SignupOrLogin(ctx *gin.Context, req SignupOrLoginRequest) (*SignupOrLoginResponse, error)
	RefreshAccessToken(ctx *gin.Context, req RefreshAccessTokenRequest) (*RefreshAccessTokenResponse, error)
	Logout(ctx *gin.Context, logoutReq LogoutRequest) error
}
