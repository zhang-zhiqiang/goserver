package model

import (
	"time"

	"gorm.io/gorm"

	"github.com/marmotedu/goserver/pkg/auth"
)

// UserM represents a registered user.
type UserM struct {
	ID        int64          `gorm:"column:id;primary_key" json:"id"`
	Username  string         `json:"username" gorm:"column:username;not null" binding:"required" validate:"min=1,max=32"`
	Password  string         `json:"password,omitempty" gorm:"column:password;not null" binding:"required" validate:"min=5,max=128"`
	Nickname  string         `json:"nickname" gorm:"column:nickname" binding:"required" validate:"required,min=1,max=30"`
	Email     string         `json:"email" gorm:"column:email" binding:"required" validate:"required,email,min=1,max=100"`
	Phone     string         `json:"phone" gorm:"column:phone" binding:"required" validate:"required,phone,min=1,max=16"`
	CreatedAt time.Time      `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deletedAt" json:"deletedAt"`
}

// TableName maps to mysql table name.
func (u *UserM) TableName() string {
	return "user"
}

// BeforeCreate run before create database record.
func (u *UserM) BeforeCreate(tx *gorm.DB) (err error) {
	// Encrypt the user password.
	u.Password, err = auth.Encrypt(u.Password)
	if err != nil {
		return err
	}

	return nil
}
