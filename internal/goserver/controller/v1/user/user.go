// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"github.com/marmotedu/goserver/internal/goserver/biz"
	"github.com/marmotedu/goserver/internal/goserver/store"
)

// UserController create a user handler used to handle request for user resource.
type UserController struct {
	b biz.BizFactory
}

// New creates a user handler.
func New(ds store.IStore) *UserController {
	return &UserController{b: biz.NewBiz(ds)}
}
