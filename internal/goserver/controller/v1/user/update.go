// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"

	"github.com/marmotedu/goserver/internal/pkg/log"
	v1 "github.com/marmotedu/goserver/pkg/api/goserver/v1"
	"github.com/marmotedu/goserver/pkg/core"
	"github.com/marmotedu/goserver/pkg/errno"
)

// Update update a user info by the user identifier.
func (uc *UserController) Update(c *gin.Context) {
	log.L(c).Info("Update user function called")

	var r v1.UpdateUserRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)

		return
	}

	if _, err := govalidator.ValidateStruct(r); err != nil {
		core.WriteResponse(c, errno.ErrValidation.SetMessage(err.Error()), nil)

		return
	}

	if err := uc.b.Users().Update(c, c.Param("name"), &r); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, nil)
}
