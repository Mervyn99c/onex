// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package validate

import (
	"github.com/Rosas99/smsx/internal/sms/store"
	"net/http"

	"github.com/gin-gonic/gin"
)

// todo 使用这种方式更加简单
// 这样不需要使用自定义的了

// Validation make sure users have the right resource permission and operation.
func Validation(ds store.IStore) gin.HandlerFunc {
	return func(c *gin.Context) {

		switch c.FullPath() {
		// todo 根据url校验：
		// 参数非空 模板校验 手机号规则校验
		// 手机号白名单校验
		case "/v1/users":
			if c.Request.Method != http.MethodPost {
				//core.WriteResponse(c, errors.WithCode(code.ErrPermissionDenied, ""), nil)
				c.Abort()

				return
			}
		case "/v1/users/:name", "/v1/users/:name/change_password":
			username := c.GetString("username")
			if c.Request.Method == http.MethodDelete ||
				(c.Request.Method != http.MethodDelete && username != c.Param("name")) {
				//core.WriteResponse(c, errors.WithCode(code.ErrPermissionDenied, ""), nil)
				c.Abort()

				return
			}
		default:
		}

		c.Next()
	}
}
