// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"github.com/gin-gonic/gin"

	"github.com/marmotedu/goserver/internal/pkg/log"
	"github.com/marmotedu/goserver/pkg/core"
)

// Delete delete an user by the user identifier. Only administrator can call this function.
func (uc *UserController) Delete(c *gin.Context) {
	log.L(c).Info("Delete user function called")

	if err := uc.b.Users().Delete(c, c.Param("name")); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, nil)
}
