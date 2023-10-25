package google

import (
	"context"
	"ketalk-api/pkg/provider/model"
)

func (g *googleClient) QueryUserDetails(ctx context.Context, token model.Token) (*model.ProviderUserDetails, error) {
	var result GoogleUserDetails
	err := g.Get(ctx, &token.AccessToken, "https://www.googleapis.com/oauth2/v2/userinfo", &result)
	if err != nil {
		return nil, err
	}
	return result.To(token), nil
}
