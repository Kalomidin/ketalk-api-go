package common

import (
	"context"
	"ketalk-api/jwt"

	"github.com/google/uuid"
)

func GetUserId(ctx context.Context) (uuid.UUID, error) {
	tokenData, err := jwt.GetJWTToken(ctx)
	if err != nil {
		return uuid.UUID{}, err
	}
	return uuid.Parse(tokenData.Subject)
}
