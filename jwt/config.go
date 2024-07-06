package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Config struct {
	Issuer                     string        `yaml:"issuer" env:"AUTH_JWT_ISSUER"`
	Key                        string        `yaml:"key" env:"AUTH_JWT_KEY"`
	KeyID                      string        `yaml:"keyID" env:"AUTH_JWT_KEYID"`
	ValidDuration              time.Duration `yaml:"validDuration" env:"AUTH_JWT_TOKEN_VALID_DURATION" env-default:"10h"`
	RefreshTokenExpiryDuration time.Duration `yaml:"refreshTokenExpiryDuration" env:"AUTH_JWT_REFRESH_TOKEN_EXPIRY_DURATION" env-default:"100h"`
}

var JWTMethod = jwt.SigningMethodHS256

func (c *Config) SetDefaults() {
	if c.Issuer == "" {
		c.Issuer = "magicx-ai"
	}
	if c.ValidDuration <= 0 {
		c.ValidDuration = time.Minute * 10
	}
}
