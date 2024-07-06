package config

import (
	"ketalk-api/jwt"
	conn_redis "ketalk-api/pkg/manager/conversation/redis"
	"ketalk-api/pkg/manager/conversation/ws"
	"ketalk-api/pkg/provider/google"
	"ketalk-api/storage"
)

type Config struct {
	Google           google.Config                  `yaml:"google"`
	Auth             jwt.Config                     `yaml:"auth"`
	AzureBlobStorage storage.AzureBlobStorageConfig `yaml:"azure"`
	R2Storage        storage.R2CloudFlareConfig     `yaml:"r2"`
	DB               Postgres                       `yaml:"db"`
	Redis            conn_redis.Config              `yaml:"redis"`
	WebSocketServer  ws.Config                      `yaml:"ws"`
}
