// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"context"
	"errors"
	"regexp"
	"sync"

	"github.com/jinzhu/copier"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"

	"github.com/marmotedu/goserver/internal/goserver/store"
	"github.com/marmotedu/goserver/internal/pkg/log"
	"github.com/marmotedu/goserver/internal/pkg/model"
	v1 "github.com/marmotedu/goserver/pkg/api/goserver/v1"
	"github.com/marmotedu/goserver/pkg/auth"
	"github.com/marmotedu/goserver/pkg/errno"
	"github.com/marmotedu/goserver/pkg/token"
)

// UserBiz defines functions used to handle user request.
type UserBiz interface {
	ChangePassword(ctx context.Context, username string, r *v1.ChangePasswordRequest) error
	Login(ctx context.Context, r *v1.LoginRequest) (*v1.LoginResponse, error)
	Create(ctx context.Context, r *v1.CreateUserRequest) error
	Update(ctx context.Context, username string, r *v1.UpdateUserRequest) error
	Delete(ctx context.Context, username string) error
	Get(ctx context.Context, username string) (*v1.GetUserResponse, error)
	List(ctx context.Context, offset, limit int) (*v1.ListUserResponse, error)
}

// The implementation of UserBiz interface.
type userBiz struct {
	ds store.IStore
}

// Make sure that userBiz implements the UserBiz interface.
// We can find this problem in the compile stage with the following assignment statement.
var _ UserBiz = (*userBiz)(nil)

func New(ds store.IStore) *userBiz {
	return &userBiz{ds: ds}
}

// ChangePassword is the implementation of the `ChangePassword` method in UserBiz interface.
func (b *userBiz) ChangePassword(ctx context.Context, username string, r *v1.ChangePasswordRequest) error {
	userM, err := b.ds.Users().Get(ctx, username)
	if err != nil {
		return err
	}

	if err := auth.Compare(userM.Password, r.OldPassword); err != nil {
		return errno.ErrPasswordIncorrect
	}

	userM.Password, _ = auth.Encrypt(r.NewPassword)
	if err := b.ds.Users().Update(ctx, userM); err != nil {
		return err
	}

	return nil
}

// Login is the implementation of the `Login` method in UserBiz interface.
func (b *userBiz) Login(ctx context.Context, r *v1.LoginRequest) (*v1.LoginResponse, error) {
	// Get the user information by the login username.
	user, err := b.ds.Users().Get(ctx, r.Username)
	if err != nil {
		return nil, errno.ErrUserNotFound
	}

	// Compare the login password with the user password.
	if err := auth.Compare(user.Password, r.Password); err != nil {
		return nil, errno.ErrPasswordIncorrect
	}

	// Sign the json web token.
	t, err := token.Sign(r.Username)
	if err != nil {
		return nil, errno.ErrToken
	}

	return &v1.LoginResponse{Token: t}, nil
}

// Create is the implementation of the `Create` method in UserBiz interface.
func (b *userBiz) Create(ctx context.Context, r *v1.CreateUserRequest) error {
	var userM model.UserM
	_ = copier.Copy(&userM, r)

	if err := b.ds.Users().Create(ctx, &userM); err != nil {
		if match, _ := regexp.MatchString("Duplicate entry '.*' for key 'username'", err.Error()); match {
			return errno.ErrUserAlreadyExist
		}

		return err
	}

	return nil
}

// Delete is the implementation of the `Delete` method in UserBiz interface.
func (b *userBiz) Delete(ctx context.Context, username string) error {
	if err := b.ds.Users().Delete(ctx, username); err != nil {
		return err
	}

	return nil
}

// Get is the implementation of the `Get` method in UserBiz interface.
func (b *userBiz) Get(ctx context.Context, username string) (*v1.GetUserResponse, error) {
	user, err := b.ds.Users().Get(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrUserNotFound
		}

		return nil, err
	}

	var resp v1.GetUserResponse
	_ = copier.Copy(&resp, user)

	resp.CreatedAt = user.CreatedAt.Format("2006-01-02 15:04:05")
	resp.UpdatedAt = user.UpdatedAt.Format("2006-01-02 15:04:05")

	return &resp, nil
}

// Update is the implementation of the `Update` method in UserBiz interface.
func (b *userBiz) Update(ctx context.Context, username string, user *v1.UpdateUserRequest) error {
	userM, err := b.ds.Users().Get(ctx, username)
	if err != nil {
		return err
	}

	if user.Email != nil {
		userM.Email = *user.Email
	}

	if user.Nickname != nil {
		userM.Nickname = *user.Nickname
	}

	if user.Phone != nil {
		userM.Phone = *user.Phone
	}

	if err := b.ds.Users().Update(ctx, userM); err != nil {
		return err
	}

	return nil
}

// List is the implementation of the `List` method in UserBiz interface.
func (b *userBiz) List(ctx context.Context, offset, limit int) (*v1.ListUserResponse, error) {
	count, list, err := b.ds.Users().List(ctx, offset, limit)
	if err != nil {
		log.L(ctx).Errorf("list users from storage failed: %s", err.Error())

		return nil, err
	}

	var m sync.Map
	eg, ctx := errgroup.WithContext(ctx)
	// Improve query efficiency in parallel
	for _, item := range list {
		user := item
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				shadowedPass, err := auth.Shadow(user.Password)
				if err != nil {
					log.L(ctx).Errorw("Failed to shadow password", "err", err)
					return err
				}

				m.Store(user.ID, &v1.UserInfo{
					Username:  user.Username,
					Nickname:  user.Nickname,
					Password:  shadowedPass,
					Email:     user.Email,
					Phone:     user.Email,
					CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
					UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
				})

				return nil
			}
		})
	}

	if err := eg.Wait(); err != nil {
		log.L(ctx).Errorw("Failed to wait all function calls returned", "err", err)
		return nil, err
	}

	users := make([]*v1.UserInfo, 0, len(list))
	for _, item := range list {
		user, _ := m.Load(item.ID)
		users = append(users, user.(*v1.UserInfo))
	}

	log.L(ctx).Debugw("Get uses from backend storage", "count", len(users))

	return &v1.ListUserResponse{TotalCount: count, Users: users}, nil
}
