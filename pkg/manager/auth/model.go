package auth_manager

import (
	"context"

	"github.com/google/uuid"
)

type GoogleToken struct {
	IdToken      string
	AccessToken  string
	RefreshToken string
}

type SignUpDetails struct {
	UserName string
	Email    string
	Password string
}

type ProviderToken struct {
	GoogleToken *GoogleToken
}

type SignupOrLoginRequest struct {
	ProviderToken *ProviderToken
	DeviceID      string
	DeviceOS      string

	SignUpDetails *SignUpDetails
}

type SignupOrLoginResponse struct {
	Id           uuid.UUID
	UserName     string
	Email        string
	Image        *string
	AuthToken    string
	RefreshToken string
}

type AuthManager interface {
	SignupOrLogin(ctx context.Context, req SignupOrLoginRequest) (*SignupOrLoginResponse, error)
}
