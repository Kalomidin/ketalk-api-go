package provider

import (
	"context"
	"fmt"
	"ketalk-api/pkg/provider/model"

	"golang.org/x/oauth2"
)

type ProviderClient interface {
	QueryUserDetails(ctx context.Context, externalIdentifier *model.ProviderToken) (*model.ProviderUserDetails, error)
	VerifyIDToken(ctx context.Context, externalIdentifier *model.ProviderToken) (*model.ProviderUserDetails, error)
	UpdateAccessToken(ctx context.Context, externalIdentifier *model.ProviderToken) (*oauth2.Token, error)
}

type providerClient struct {
	googleClient model.ProviderClient
}

func NewProviderClient(googleClient model.ProviderClient) ProviderClient {
	return &providerClient{
		googleClient: googleClient,
	}
}

func (p *providerClient) QueryUserDetails(ctx context.Context, providerToken *model.ProviderToken) (*model.ProviderUserDetails, error) {
	if providerToken.GoogleToken != nil {
		return p.googleClient.QueryUserDetails(ctx, providerToken.GoogleToken.Token)
	}
	return nil, fmt.Errorf("unsupported external identifier")
}

func (p *providerClient) VerifyIDToken(ctx context.Context, providerToken *model.ProviderToken) (*model.ProviderUserDetails, error) {
	if providerToken != nil && providerToken.GoogleToken != nil {
		_, err := p.googleClient.VerifyIDToken(ctx, providerToken.GoogleToken.IdToken)
		if err != nil {
			return nil, err
		}
		// TODO: Maybe not necessary to query user details here
		// We are currently querying for signup/signin. We can query only for signup
		// Moreover, verify id token already returns some of the user details
		return p.QueryUserDetails(ctx, providerToken)
	}
	return nil, fmt.Errorf("unsupported external identifier")
}

func (p *providerClient) UpdateAccessToken(ctx context.Context, providerToken *model.ProviderToken) (*oauth2.Token, error) {
	if providerToken != nil && providerToken.GoogleToken != nil {
		return p.googleClient.UpdateAccessToken(ctx, providerToken.GoogleToken.Token)
	}
	return nil, fmt.Errorf("unsupported external identifier")
}
