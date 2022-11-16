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

// PostStore defines the post storage interface.
type PostStore interface {
	Create(ctx context.Context, post *model.PostM) error
	Update(ctx context.Context, post *model.PostM) error
	Delete(ctx context.Context, username string, postIDs []string) error
	Get(ctx context.Context, username, postID string) (*model.PostM, error)
	List(ctx context.Context, username string, offset, limit int) (int64, []*model.PostM, error)
}

// The implementation of PostStore interface.
type posts struct {
	db *gorm.DB
}

// Make sure that posts implements the PostStore interface.
var _ PostStore = (*posts)(nil)

func newPosts(db *gorm.DB) *posts {
	return &posts{db}
}

// Create creates a new post account.
func (u *posts) Create(ctx context.Context, post *model.PostM) error {
	return u.db.Create(&post).Error
}

// Update updates an post account information.
func (u *posts) Update(ctx context.Context, post *model.PostM) error {
	return u.db.Save(post).Error
}

// Delete deletes the post by the post identifier.
func (u *posts) Delete(ctx context.Context, username string, postIDs []string) error {
	err := u.db.Where("username = ? and postID in (?)", username, postIDs).Delete(&model.PostM{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return nil
}

// Get return an post by the post identifier.
func (u *posts) Get(ctx context.Context, username, postID string) (*model.PostM, error) {
	var post model.PostM
	if err := u.db.Where("username = ? and postID = ?", username, postID).First(&post).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

// List list limited post records according to the given offset and limit.
func (u *posts) List(ctx context.Context, username string, offset, limit int) (count int64, ret []*model.PostM, err error) {
	err = u.db.Where("username = ?", username).Offset(offset).Limit(limit).Order("id desc").Find(&ret).
		Offset(-1).
		Limit(-1).
		Count(&count).
		Error

	return
}
