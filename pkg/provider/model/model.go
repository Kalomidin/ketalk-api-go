package model

import (
	"context"

	"golang.org/x/oauth2"
)

type ProviderName string

const (
	ProviderNameGoogle  ProviderName = "google"
	ProviderNameOutlook ProviderName = "outlook"
)

type IdTokenUserDetails struct {
	ExternalID    string
	FirstName     string
	LastName      string
	VerifiedEmail string
}

type ProviderUserDetails struct {
	Email        string
	FirstName    string
	ExternalID   string
	Image        string
	Locale       string
	ProviderName ProviderName
	Token        Token
}

type Token struct {
	AccessToken  string
	RefreshToken string
}

type GoogleToken struct {
	IdToken string
	Token   Token
}

type ProviderToken struct {
	GoogleToken *GoogleToken
}

type ProviderClient interface {
	QueryUserDetails(ctx context.Context, token Token) (*ProviderUserDetails, error)
	VerifyIDToken(ctx context.Context, idToken string) (*IdTokenUserDetails, error)
	UpdateAccessToken(ctx context.Context, token Token) (*oauth2.Token, error)
}
