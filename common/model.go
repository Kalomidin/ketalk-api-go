package common

import (
	"time"

	"gorm.io/gorm"
)

type CreatedUpdated struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreatedDeleted struct {
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type CreatedUpdatedDeleted struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
