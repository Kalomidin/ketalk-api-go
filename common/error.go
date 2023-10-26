package common

import (
	"fmt"

	"gorm.io/gorm"
)

var ErrMoreThanOneRowUpdated = fmt.Errorf("more than one row updated")
var ErrInvalidInput = fmt.Errorf("invalid input")
var ErrRecordNotFound = gorm.ErrRecordNotFound
