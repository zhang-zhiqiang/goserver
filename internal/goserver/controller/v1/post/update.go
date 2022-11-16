// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package post

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"

	"github.com/marmotedu/goserver/internal/pkg/known"
	"github.com/marmotedu/goserver/internal/pkg/log"
	v1 "github.com/marmotedu/goserver/pkg/api/goserver/v1"
	"github.com/marmotedu/goserver/pkg/core"
	"github.com/marmotedu/goserver/pkg/errno"
)

// Update update a post info by the post identifier.
func (pc *PostController) Update(c *gin.Context) {
	log.L(c).Infow("Update post function called")

	var r v1.UpdatePostRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)

		return
	}

	if _, err := govalidator.ValidateStruct(r); err != nil {
		core.WriteResponse(c, errno.ErrValidation.SetMessage(err.Error()), nil)

		return
	}

	if err := pc.b.Posts().Update(c, c.GetString(known.XUsernameKey), c.Param("postID"), &r); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, nil)
}
