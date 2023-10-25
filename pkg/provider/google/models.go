package google

import (
	"ketalk-api/pkg/provider/model"
)

type GoogleUserDetails struct {
	UserEmail          string `json:"email"`
	GoogleId           string `json:"id"`
	VerifiedEmail      bool   `json:"verified_email"`
	GivenName          string `json:"given_name"`
	UserProfilePicture string `json:"picture"`
	UserLocale         string `json:"locale"`
}

func (u *GoogleUserDetails) To(token model.Token) *model.ProviderUserDetails {
	return &model.ProviderUserDetails{
		ExternalID:   u.GoogleId,
		FirstName:    u.GivenName,
		Email:        u.UserEmail,
		ProviderName: model.ProviderNameGoogle,
		Token:        token,
		Image:        u.UserProfilePicture,
	}
}

type GoogleIdTokenUserInformation struct {
	ExternalID    string `json:"id"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	VerifiedEmail string `json:"email"`
}

func (u *GoogleIdTokenUserInformation) To() *model.IdTokenUserDetails {
	return &model.IdTokenUserDetails{
		ExternalID:    u.ExternalID,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		VerifiedEmail: u.VerifiedEmail,
	}
}
