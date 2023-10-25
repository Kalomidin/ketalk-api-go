package google

import (
	"context"
	"ketalk-api/pkg/provider/model"

	"golang.org/x/oauth2"
)

func (g *googleClient) UpdateAccessToken(ctx context.Context, token model.Token) (*oauth2.Token, error) {
	config := oauth2.Config{
		ClientID:     g.cfg.ID,
		ClientSecret: g.cfg.Secret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
	}
	// Create a new token source using the refresh token
	tokenSource := config.TokenSource(ctx, &oauth2.Token{
		RefreshToken: token.RefreshToken,
	})

	// Request a new access token
	return tokenSource.Token()
}
