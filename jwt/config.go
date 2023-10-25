package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Config struct {
	Issuer                     string        `yaml:"issuer" env:"JWT_ISSUER"`
	Key                        string        `yaml:"key" env:"JWT_KEY"`
	KeyID                      string        `yaml:"keyID" env:"JWT_KEYID"`
	ValidDuration              time.Duration `yaml:"validDuration" env:"JWT_TOKEN_VALID_DURATION" default:"10h"`
	RefreshTokenExpiryDuration time.Duration `yaml:"refreshTokenExpiryDuration" env:"JWT_REFRESH_TOKEN_EXPIRY_DURATION" default:"100h"`
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
