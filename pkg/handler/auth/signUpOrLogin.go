package auth_handler

import (
	"ketalk-api/common"
	auth_manager "ketalk-api/pkg/manager/auth"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GoogleToken struct {
	IdToken      string `json:"idToken"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type ProviderToken struct {
	GoogleToken *GoogleToken `json:"googleToken"`
}

type SignUpDetails struct {
	UserName string `json:"userName"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupOrLoginRequest struct {
	ProviderToken *ProviderToken `json:"providerToken"`
	DeviceID      string         `json:"deviceId"`
	DeviceOS      string         `json:"deviceOs"`

	// TODO: will be removed once SSO integrated
	SignUpDetails *SignUpDetails  `json:"signUpDetails"`
	Location      common.Location `json:"location"`
}

type SignupOrLoginResponse struct {
	Id           uuid.UUID `json:"id"`
	UserName     string    `json:"userName"`
	Email        string    `json:"email"`
	Image        *string   `json:"image"`
	AuthToken    string    `json:"authToken"`
	RefreshToken string    `json:"refreshToken"`
}

func (h *HttpHandler) SignupOrLogin(ctx *gin.Context, r *http.Request) (interface{}, error) {
	var req SignupOrLoginRequest
	if err := ctx.BindJSON(&req); err != nil {
		return nil, err
	}
	resp, err := h.handler.SignupOrLogin(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (h *handler) SignupOrLogin(ctx *gin.Context, req SignupOrLoginRequest) (*SignupOrLoginResponse, error) {
	manReq := auth_manager.SignupOrLoginRequest{
		DeviceID: req.DeviceID,
		DeviceOS: req.DeviceOS,
		Location: req.Location,
	}
	if req.ProviderToken != nil {
		manReq.ProviderToken = &auth_manager.ProviderToken{}
		if req.ProviderToken.GoogleToken != nil {
			manReq.ProviderToken.GoogleToken = &auth_manager.GoogleToken{
				IdToken:      req.ProviderToken.GoogleToken.IdToken,
				AccessToken:  req.ProviderToken.GoogleToken.AccessToken,
				RefreshToken: req.ProviderToken.GoogleToken.RefreshToken,
			}
		}
	}
	if req.SignUpDetails != nil {
		manReq.SignUpDetails = &auth_manager.SignUpDetails{
			UserName: req.SignUpDetails.UserName,
			Email:    req.SignUpDetails.Email,
			Password: req.SignUpDetails.Password,
		}
	}

	manResp, err := h.service.SignupOrLogin(ctx, manReq)
	if err != nil {
		return nil, err
	}
	return &SignupOrLoginResponse{
		Id:           manResp.Id,
		UserName:     manResp.UserName,
		Email:        manResp.Email,
		Image:        manResp.Image,
		AuthToken:    manResp.AuthToken,
		RefreshToken: manResp.RefreshToken,
	}, nil
}
