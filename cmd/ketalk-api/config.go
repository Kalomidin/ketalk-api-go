package main

import (
	"fmt"
	"ketalk-api/pkg/config"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
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
		"dev.yaml",
		"defaults.yaml",
	}
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return err
	}
	for _, file := range files {
		if err := cleanenv.ReadConfig(fmt.Sprintf("./%s/%s", dir, file), cfg); err == nil {
			return nil
		}
	}
	return nil
}
