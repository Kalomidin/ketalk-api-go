package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, userId uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	MigrateUser() error
}

type repository struct {
	*gorm.DB
}

func NewRepository(ctx context.Context, db *gorm.DB) Repository {
	return &repository{
		db,
	}
}

func (r *repository) CreateUser(ctx context.Context, user *User) error {
	res := r.Create(user)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected != 1 {
		return fmt.Errorf("more than one row updated")
	}
	return nil
}

func (r *repository) GetUser(ctx context.Context, userId uuid.UUID) (*User, error) {
	var user User
	resp := r.DB.Where("id = ?", userId).First(&user)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return &user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	resp := r.DB.Where("email = ?", email).First(&user)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return &user, nil
}

func (r *repository) UpdateUser(ctx context.Context, user *User) error {
	res := r.Updates(user)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *repository) MigrateUser() error {
	return r.AutoMigrate(&User{})
}
