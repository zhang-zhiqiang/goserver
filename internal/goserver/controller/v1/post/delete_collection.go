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

// DeleteCollection batch delete posts by multiple post ids.
func (pc *PostController) DeleteCollection(c *gin.Context) {
	log.L(c).Info("Batch delete post function called")

	postIDs := c.QueryArray("postID")
	if err := pc.b.Posts().DeleteCollection(c, c.GetString(known.XUsernameKey), postIDs); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, nil)
}
