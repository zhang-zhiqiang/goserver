// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package post

import (
	"github.com/marmotedu/goserver/internal/goserver/biz"
	"github.com/marmotedu/goserver/internal/goserver/store"
)

// PostController create a post handler used to handle request for post resource.
type PostController struct {
	b biz.BizFactory
}

// New creates a post handler.
func New(ds store.IStore) *PostController {
	return &PostController{b: biz.NewBiz(ds)}
}
