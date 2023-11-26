package common

import (
	"context"
	"fmt"
	"ketalk-api/jwt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

func GetUserId(ctx context.Context) (uuid.UUID, error) {
	tokenData, err := jwt.GetJWTToken(ctx)
	if err != nil {
		return uuid.UUID{}, err
	}
	return uuid.Parse(tokenData.Subject)
}

func GetLocation(req *http.Request) (*Location, error) {
	latitude := req.Header.Get("latitude")
	longitude := req.Header.Get("longitude")
	if latitude == "" || longitude == "" {
		return nil, fmt.Errorf("latitude or longitude is empty")
	}

	lat, err := strconv.ParseFloat(latitude, 64)
	if err != nil {
		return nil, err
	}
	lon, err := strconv.ParseFloat(longitude, 64)
	if err != nil {
		return nil, err
	}
	return &Location{
		Latitude:  lat,
		Longitude: lon,
	}, nil
}
