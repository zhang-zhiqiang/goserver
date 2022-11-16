// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package post

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/marmotedu/goserver/internal/pkg/known"
	"github.com/marmotedu/goserver/internal/pkg/log"
	v1 "github.com/marmotedu/goserver/pkg/api/goserver/v1"
	"github.com/marmotedu/goserver/pkg/core"
	"github.com/marmotedu/goserver/pkg/errno"
)

// List list the posts in the storage.
func (pc *PostController) List(c *gin.Context) {
	log.L(c).Info("List post function called.")
	fmt.Println("33333333333333333333333", c.GetString(known.XRequestIDKey))

	var r v1.ListPostRequest
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)

		return
	}

	resp, err := pc.b.Posts().List(c, c.GetString(known.XUsernameKey), r.Offset, r.Limit)
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, resp)
}
