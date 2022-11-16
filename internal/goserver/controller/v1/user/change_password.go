// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/component-base/pkg/core"

	v1 "github.com/marmotedu/goserver/pkg/api/goserver/v1"
	"github.com/marmotedu/goserver/pkg/errno"

	"github.com/marmotedu/iam/pkg/log"
)

// ChangePassword change the user's password by the user identifier.
func (uc *UserController) ChangePassword(c *gin.Context) {
	log.L(c).Info("Change password function called.")

	var r v1.ChangePasswordRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)

		return
	}

	if _, err := govalidator.ValidateStruct(r); err != nil {
		core.WriteResponse(c, errno.ErrValidation.SetMessage(err.Error()), nil)

		return
	}

	if err := uc.b.Users().ChangePassword(c, c.Param("name"), &r); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, nil)
}
