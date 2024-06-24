// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/onex.
//

package message

import (
	"github.com/gin-gonic/gin"
	"github.com/superproj/onex/internal/pkg/core"
	v1 "github.com/superproj/onex/pkg/api/sms/v1"
)

func (b *MessageController) Send(c *gin.Context) {
	var r v1.CreateTemplateRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, err, nil)
		// todo 了解gin如何返回错误 如：
		/*
			c.JSON(http.StatusBadRequest, gin.H{
					"err": err.Error(),
				})
		*/
	}
	template, err := b.svc.Send(c, &r)
	if err != nil {
		core.WriteResponse(c, err, nil)

	}
	core.WriteResponse(c, nil, template)

	// todo log kpi

}
