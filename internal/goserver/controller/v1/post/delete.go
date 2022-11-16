// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package post

import (
	"github.com/gin-gonic/gin"

	"github.com/marmotedu/goserver/internal/pkg/known"
	"github.com/marmotedu/goserver/internal/pkg/log"
	"github.com/marmotedu/goserver/pkg/core"
)

// Delete delete an post by the post identifier.
func (pc *PostController) Delete(c *gin.Context) {
	log.L(c).Info("Delete post function called")

	if err := pc.b.Posts().Delete(c, c.GetString(known.XUsernameKey), c.Param("postID")); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, nil)
}
