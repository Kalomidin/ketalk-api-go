package jwt

import (
	"context"
	"encoding/json"
	"time"

	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type ContextKey int

const (
	JWTTokenRawContextKey ContextKey = iota
	JWTTokenContextKey
)

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"username,omitempty"`
}

func GetJWTToken(ctx context.Context) (*JWTClaims, error) {
	jwtToken := ctx.Value(JWTTokenContextKey)
	if token, ok := jwtToken.(JWTClaims); ok {
		return &token, nil
	}
	fmt.Printf("token is not valid2, jwtToken: %+v\n", jwtToken)
	return nil, fmt.Errorf("invalid token")
}

func IssueToken(cfg Config, userID uuid.UUID, additionalDetails map[string]interface{}) (string, error) {
	subject := fmt.Sprintf("%+v", userID)

	jwtToken := jwt.RegisteredClaims{
		Issuer:    cfg.Issuer,
		Subject:   subject,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.ValidDuration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	jwtTokenBytes, err := json.Marshal(jwtToken)
	if err != nil {
		return "", fmt.Errorf("failed to marshal jwt token: %+v", err)
	}
	dataValues := make(map[string]interface{}, len(additionalDetails)+7)
	if err = json.Unmarshal(jwtTokenBytes, &dataValues); err != nil {
		return "", fmt.Errorf("failed to unmarshal jwt token: %+v", err)
	}
	for detail, value := range additionalDetails {
		dataValues[detail] = value
	}
	if cfg.Key == "" || cfg.KeyID == "" {
		return "", fmt.Errorf("should set secret key")
	}
	token := jwt.NewWithClaims(JWTMethod, jwt.MapClaims(dataValues))
	token.Header["kid"] = cfg.KeyID
	return token.SignedString([]byte(cfg.Key))
}

func ParseJWTToken(ctx context.Context, cfg Config, tokenString string) (context.Context, error) {
	claims := JWTClaims{}
	jwtToken, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != JWTMethod.Name {
			return nil, fmt.Errorf("invalid method")
		}
		if keyID, ok := token.Header["kid"].(string); ok && len(keyID) > 0 && cfg.KeyID == keyID {
			return []byte(cfg.Key), nil
		}
		return nil, fmt.Errorf("invalid key id")
	})
	if err != nil {
		return ctx, fmt.Errorf("failed to verify token: %+v", err)
	}
	if !jwtToken.Valid {
		return ctx, fmt.Errorf("invalid token")
	}
	return context.WithValue(context.WithValue(ctx, JWTTokenRawContextKey, tokenString), JWTTokenContextKey, claims), nil
}
