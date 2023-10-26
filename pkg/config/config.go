package config

import (
	"ketalk-api/jwt"
	"ketalk-api/pkg/provider/google"
	"ketalk-api/storage"
)

type Config struct {
	Google           google.Config                  `yaml:"google"`
	Auth             jwt.Config                     `yaml:"auth"`
	AzureBlobStorage storage.AzureBlobStorageConfig `yaml:"azure"`
	DB               Postgres                       `yaml:"db"`
}
