// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package store

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/marmotedu/goserver/internal/pkg/model"
)

// UserStore defines the user storage interface.
type UserStore interface {
	Create(ctx context.Context, user *model.UserM) error
	Update(ctx context.Context, user *model.UserM) error
	Delete(ctx context.Context, username string) error
	Get(ctx context.Context, username string) (*model.UserM, error)
	List(ctx context.Context, offset, limit int) (int64, []*model.UserM, error)
}

// The implementation of UserStore interface.
type users struct {
	db *gorm.DB
}

// Make sure that users implements the UserStore interface.
var _ UserStore = (*users)(nil)

func newUsers(db *gorm.DB) *users {
	return &users{db}
}

// Create creates a new user account.
func (u *users) Create(ctx context.Context, user *model.UserM) error {
	return u.db.Create(&user).Error
}

// Update updates an user account information.
func (u *users) Update(ctx context.Context, user *model.UserM) error {
	return u.db.Save(user).Error
}

// Delete deletes the user by the user identifier.
func (u *users) Delete(ctx context.Context, username string) error {
	err := u.db.Where("username = ?", username).Delete(&model.UserM{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return nil
}

// Get return an user by the user identifier.
func (u *users) Get(ctx context.Context, username string) (*model.UserM, error) {
	var user model.UserM
	if err := u.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// List list limited user records according to the given offset and limit.
func (u *users) List(ctx context.Context, offset, limit int) (count int64, ret []*model.UserM, err error) {
	err = u.db.Offset(offset).Limit(limit).Order("id desc").Find(&ret).
		Offset(-1).
		Limit(-1).
		Count(&count).
		Error

	return
}
