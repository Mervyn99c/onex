// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/onex.
//

package template

import (
	"github.com/gin-gonic/gin"
	"github.com/superproj/onex/internal/pkg/core"
	v1 "github.com/superproj/onex/pkg/api/sms/v1"
)

func (b *TemplateController) Delete(c *gin.Context) {
	var r v1.DeleteTemplateRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, err, nil)
	}
	template, err := b.svc.DeleteTemplate(c, &r)
	if err != nil {
		core.WriteResponse(c, err, nil)

	}
	core.WriteResponse(c, nil, template)

}
