package postgres

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type ConfigPostgres interface {
	GetDatabase() string
	GetSchema() string
	GetUserName() string
	GetPassword() string
	GetAddress() string
	GetHost() string
	GetPort() string
	GetAlias() string
	GetMaxWaitForConnection() time.Duration
	GetMaxConns() int
	SetDatabase(string)
	SetSchema(string)
	GetSSLMode() string
}

func InitDB(ctx context.Context, cfg ConfigPostgres) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s  dbname=%s sslmode=%s",
		cfg.GetHost(), cfg.GetPort(), cfg.GetUserName(), cfg.GetPassword(), cfg.GetDatabase(), cfg.GetSSLMode())

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   fmt.Sprintf("%s.", cfg.GetSchema()),
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.GetMaxConns())
	sqlDB.SetConnMaxIdleTime(cfg.GetMaxWaitForConnection())

	return db, nil
}
