package google

import (
	"context"
	"encoding/json"
	"fmt"
	"ketalk-api/pkg/provider/model"
	"strings"

	"golang.org/x/exp/slices"
	"google.golang.org/api/idtoken"
)

var issuers = []string{
	"https://accounts.google.com",
	"accounts.google.com",
}

type Claims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
}

func (v *googleClient) VerifyIDToken(ctx context.Context, idToken string) (*model.IdTokenUserDetails, error) {
	var claims Claims

	payload, err := idtoken.Validate(ctx, idToken, v.cfg.Audience)
	if err != nil {
		return nil, err
	}

	if !slices.Contains(issuers, payload.Issuer) {
		return nil, fmt.Errorf("invalid issuer: %v", payload.Issuer)
	}

	claimsInBytes, err := json.Marshal(payload.Claims)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal claims: %+v", err)
	}

	err = json.Unmarshal(claimsInBytes, &claims)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal claims: %w", err)
	}
	var info GoogleIdTokenUserInformation

	info.ExternalID = payload.Subject
	if info.ExternalID == "" {
		return nil, fmt.Errorf("empty google id")
	}
	info.FirstName = claims.FirstName
	info.LastName = claims.LastName
	if !claims.EmailVerified {
		return nil, fmt.Errorf("unverified email")
	}
	info.VerifiedEmail = strings.TrimSpace(strings.ToLower(claims.Email))
	return info.To(), nil
}
