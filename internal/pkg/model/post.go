package model

import (
	"time"

	validator "gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"

	"github.com/marmotedu/goserver/pkg/util/id"
)

// PostM represents a blog for a user.
type PostM struct {
	ID        int64     `gorm:"column:id;primary_key" json:"id"`
	Username  string    `json:"username,omitempty" gorm:"column:username;not null"`
	PostID    string    `json:"postID,omitempty" gorm:"column:postID;not null"`
	Title     string    `json:"title" gorm:"column:title;not null" binding:"required" validate:"min=1,max=256"`
	Content   string    `json:"content" gorm:"column:content" binding:"required" validate:"min=1,max=10240"`
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

// TableName maps to mysql table name.
func (p *PostM) TableName() string {
	return "post"
}

// BeforeCreate run before create database record.
func (p *PostM) BeforeCreate(tx *gorm.DB) error {
	p.PostID = "post-" + id.GenShortId()

	return nil
}

// Validate the fields.
func (p *PostM) Validate() error {
	validate := validator.New()

	return validate.Struct(p)
}
