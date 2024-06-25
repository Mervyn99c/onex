// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Rosas99/smsx.
//

package sms

import (
	"github.com/Rosas99/smsx/internal/pkg/client/usercenter"
	"github.com/Rosas99/smsx/internal/pkg/idempotent"
	"github.com/Rosas99/smsx/internal/pkg/middleware/header"
	"github.com/Rosas99/smsx/internal/pkg/middleware/trace"
	"github.com/Rosas99/smsx/internal/sms/biz"
	"github.com/Rosas99/smsx/internal/sms/middleware/auth"
	"github.com/Rosas99/smsx/internal/sms/middleware/validate"
	"github.com/Rosas99/smsx/internal/sms/service"
	"github.com/Rosas99/smsx/internal/sms/store"
	"github.com/Rosas99/smsx/internal/sms/store/mysql"
	"github.com/Rosas99/smsx/pkg/db"
	"github.com/Rosas99/smsx/pkg/log"
	genericoptions "github.com/Rosas99/smsx/pkg/options"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

// Config represents the configuration of the service.
type Config struct {
	FakeStore     bool
	GRPCOptions   *genericoptions.GRPCOptions
	HTTPOptions   *genericoptions.HTTPOptions
	TLSOptions    *genericoptions.TLSOptions
	MySQLOptions  *genericoptions.MySQLOptions
	RedisOptions  *genericoptions.RedisOptions
	KafkaOptions1 *genericoptions.KafkaOptions
	KqOptions     *genericoptions.KqOptions
	Address       string
	Accounts      map[string]string

	// todo
	EtcdOptions *genericoptions.EtcdOptions
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (cfg *Config) Complete() completedConfig {
	return completedConfig{cfg}
}

type completedConfig struct {
	*Config
}

// SmsServer represents the fake server.
type SmsServer struct {
	httpsrv Server
	grpcsrv Server
	config  completedConfig
}

// New returns a new instance of Server from the given config.
func (c completedConfig) New() (*SmsServer, error) {
	var ds store.IStore

	var dbOptions db.MySQLOptions
	_ = copier.Copy(&dbOptions, c.MySQLOptions)

	ins, err := db.NewMySQL(&dbOptions)
	if err != nil {
		return nil, err
	}
	// todo 这里需要指定model
	//ins.AutoMigrate(&model.OrderM{})
	ds = mysql.NewStore(ins)

	var redisOptions db.RedisOptions
	value := &c.Config.RedisOptions
	_ = copier.Copy(&redisOptions, value)
	rds, err := db.NewRedis(&redisOptions)
	if err != nil {
		return nil, err
	}

	writer, err := NewLogger(c.KafkaOptions1, ds.Templates())
	if err != nil {
		return nil, err
	}

	//这里初始化所有writer 然后注入biz
	idt, err := idempotent.NewIdempotent(rds)
	if err != nil {
		return nil, err
	}

	biz := biz.NewBiz(ds, writer, rds, idt)

	srv := service.NewSmsServerService(biz)

	// 设置 Gin 模式
	gin.SetMode(gin.ReleaseMode)

	// 创建 Gin 引擎
	g := gin.New()

	// 并初始化路由
	// 这里注册不同的路由可以分开，如是否使用人认证中间件，分别在use 认证中间件前后
	installRouters(g, srv, c.Accounts)
	// 考虑在这里install consumer

	httpsrv, err := NewHTTPServer(c.HTTPOptions, c.TLSOptions, g)
	if err != nil {
		return nil, err
	}

	grpcsrv, err := NewGRPCServer(c.GRPCOptions, c.TLSOptions, srv)
	if err != nil {
		return nil, err
	}
	impl := usercenter.NewFakeServer()

	// gin.Recovery() 中间件，用来捕获任何 panic，并恢复
	mws := []gin.HandlerFunc{gin.Recovery(), header.NoCache, header.Cors, header.Secure,
		// todo 这里传入rds ds
		// 注意验证链路的顺序
		trace.TraceID(), auth.BasicAuth(impl), validate.Validation(ds)}
	// 添加中间件
	g.Use(mws...)
	// Need start grpc server first. http server depends on grpc sever.
	go grpcsrv.RunOrDie()

	return &SmsServer{grpcsrv: grpcsrv, httpsrv: httpsrv}, nil
}

func (s *SmsServer) Run(stopCh <-chan struct{}) error {

	log.Infof("Successfully start pump server")

	go s.httpsrv.RunOrDie()

	<-stopCh

	// The most gracefully way is to shutdown the dependent service first,
	// and then shutdown the depended service.
	s.httpsrv.GracefulStop()
	s.grpcsrv.GracefulStop()

	return nil
}
