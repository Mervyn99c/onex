// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/onex.
//

package sms

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/superproj/onex/internal/pkg/core"
	gin2 "github.com/superproj/onex/internal/pkg/middleware/gin"
	"github.com/superproj/onex/internal/sms/controller/v1/message"
	"github.com/superproj/onex/internal/sms/controller/v1/template"
	"github.com/superproj/onex/internal/sms/service"
	v1api "github.com/superproj/onex/pkg/api/sms/v1"
)

func installRouters(g *gin.Engine, svc *service.SmsServerService, accounts map[string]string) {
	// 注册 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		core.WriteResponse(c, v1api.ErrorOrderAlreadyExists("route not found"), nil)
	})

	// 注册 pprof 路由
	pprof.Register(g)

	// 创建 v1 路由分组，并添加认证中间件
	//v1 := g.Group("/v1", mw.BasicAuth(accounts))
	tl := template.New(svc)
	ms := message.New(svc)

	v1 := g.Group("/v1")
	{
		// 创建 blocks 路由分组
		templatev1 := v1.Group("/template")
		{

			templatev1.Use(gin2.Validator())

			templatev1.GET("", tl.Get)
			templatev1.GET("", tl.List)
			templatev1.GET("", tl.Create)
			templatev1.GET("", tl.Update)
			templatev1.GET("", tl.Delete)

		}

		msgv1 := v1.Group("/msg")
		{
			msgv1.GET("", ms.Send)

		}

	}

}
