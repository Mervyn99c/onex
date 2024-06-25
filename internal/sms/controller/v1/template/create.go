// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Rosas99/smsx.
//

package template

import (
	"github.com/Rosas99/smsx/internal/pkg/core"
	v1 "github.com/Rosas99/smsx/pkg/api/sms/v1"
	"github.com/gin-gonic/gin"
)

func (b *TemplateController) Create(c *gin.Context) {
	var r v1.CreateTemplateRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, err, nil)
	}

	order, _ := b.svc.CreateTemplate(c, &r)
	// todo 增加validate
	core.WriteResponse(c, nil, order)

}
