// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package biz

//go:generate mockgen -self_package=github.com/marmotedu/goserver/internal/goserver/service/v1 -destination mock_service.go -package v1 github.com/marmotedu/goserver/internal/goserver/service/v1 Service,UserSrv,PostSrv

import (
	"github.com/marmotedu/goserver/internal/goserver/biz/post"
	"github.com/marmotedu/goserver/internal/goserver/biz/user"
	"github.com/marmotedu/goserver/internal/goserver/store"
)

// BizFactory defines functions used to return resource interface.
type BizFactory interface {
	Users() user.UserBiz
	Posts() post.PostBiz
}

var _ BizFactory = (*biz)(nil)

// The implementation of BizFactory.
type biz struct {
	ds store.IStore
}

// NewBiz returns BizFactory interface.
func NewBiz(ds store.IStore) *biz {
	return &biz{ds: ds}
}

func (b *biz) Users() user.UserBiz {
	return user.New(b.ds)
}

func (b *biz) Posts() post.PostBiz {
	return post.New(b.ds)
}
