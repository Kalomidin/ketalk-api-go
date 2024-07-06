package config

import (
	"time"
)

type Postgres struct {
	Alias string `yaml:"alias"`
	// Database defaults to service name when value is ""
	Host                 string        `yaml:"pghost"      env:"PG_HOST"      env-default:"localhost"`
	Port                 string        `yaml:"port"        env:"PG_PORT"      env-default:"5432"`
	Database             string        `yaml:"database"    env:"PG_DATABASE"  env-default:"ketalk-api"`
	Username             string        `yaml:"username"    env:"PG_USERNAME"  env-default:"root"`
	Password             string        `yaml:"password"    env:"PG_PASSWORD"  env-default:"password"`
	MaxConns             int           `yaml:"maxConns"    env:"PG_MAX_CONNS" env-default:"10"`
	MaxWaitForConnection time.Duration `yaml:"waitForConn" env:"PG_MAX_WAIT"  env-default:"5s"`
	SSLMode              string        `yaml:"sslMode" env:"PG_SSL_MODE" env-default:"disable"`
	// Schema defaults to service env when value is ""
	Schema string `yaml:"schema" env:"USER_PG_SCHEMA" env-default:"ketalk"`
}

func (cfg *Postgres) GetAlias() string {
	return cfg.Alias
}

func (cfg *Postgres) GetHost() string {
	return cfg.Host
}

func (cfg *Postgres) GetPort() string {
	return cfg.Port
}

func (cfg *Postgres) GetDatabase() string {
	return cfg.Database
}

func (cfg *Postgres) GetUserName() string {
	return cfg.Username
}

func (cfg *Postgres) GetPassword() string {
	return cfg.Password
}

func (cfg *Postgres) GetMaxConns() int {
	return cfg.MaxConns
}

func (cfg *Postgres) GetMaxWaitForConnection() time.Duration {
	return cfg.MaxWaitForConnection
}

func (cfg *Postgres) GetSchema() string {
	return cfg.Schema
}

func (cfg *Postgres) GetSSLMode() string {
	return cfg.SSLMode
}

func (cfg *Postgres) SetDatabase(db string) {
	cfg.Database = db
}

func (cfg *Postgres) SetSchema(schema string) {
	cfg.Schema = schema
}
