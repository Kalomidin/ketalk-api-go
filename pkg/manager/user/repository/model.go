package repository

import (
	"ketalk-api/common"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"primaryKey;default:gen_random_uuid()"`
	Username string
	Image    *string
	Email    string
	Password *string
	common.CreatedUpdatedDeleted
}
