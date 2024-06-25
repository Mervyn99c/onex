// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Rosas99/smsx.
//

package template

import (
	"github.com/Rosas99/smsx/internal/pkg/core"
	"github.com/gin-gonic/gin"
)

func (b *TemplateController) Update(c *gin.Context) {
	_, _ = b.svc.GetTemplate(c, nil)

	// todo 临时测试用
	core.WriteResponse(c, nil, "order")

}
