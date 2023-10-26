package main

import (
	"fmt"
	"ketalk-api/pkg/config"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server ServerConfig  `yaml:"server"`
	Config config.Config `yaml:"config"`
}

type ServerConfig struct {
	Port int `yaml:"port" env:"PORT" env-default:"8080"`
}

func (cfg *Config) Load() error {
	dir := "configs"
	files := []string{
		"defaults.yaml",
	}
	for _, file := range files {
		if err := cleanenv.ReadConfig(fmt.Sprintf("./%s/%s", dir, file), cfg); err == nil {
			return nil
		}
	}
	return cleanenv.ReadEnv(cfg)
}
