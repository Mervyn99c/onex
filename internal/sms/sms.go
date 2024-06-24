// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/onex.
//

package sms

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/superproj/onex/internal/pkg/client/usercenter"
	"github.com/superproj/onex/internal/pkg/idempotent"
	"github.com/superproj/onex/internal/pkg/middleware/header"
	"github.com/superproj/onex/internal/pkg/middleware/trace"
	"github.com/superproj/onex/internal/pkg/middleware/validate"
	"github.com/superproj/onex/internal/sms/biz"
	"github.com/superproj/onex/internal/sms/middleware/auth"
	"github.com/superproj/onex/internal/sms/mq"
	"github.com/superproj/onex/internal/sms/rule"
	"github.com/superproj/onex/internal/sms/service"
	"github.com/superproj/onex/internal/sms/store"
	"github.com/superproj/onex/internal/sms/store/mysql"
	"github.com/superproj/onex/pkg/db"
	"github.com/superproj/onex/pkg/log"
	genericoptions "github.com/superproj/onex/pkg/options"
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
	UserCenterOptions *usercenter.UserCenterOptions
	EtcdOptions       *genericoptions.EtcdOptions
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

	writer, err := mq.NewLogger(c.KafkaOptions1)
	if err != nil {
		return nil, err
	}
	//todo 注册rule
	factory := rule.NewRuleFactory()
	// 创建并注册 Rule 实例
	factory.RegisterRule("MESSAGE_COUNT_FOR_TEMPLATE_PER_DAY", &rule.MessageCountForTemplateRule{})
	// todo 其他规则

	//这里初始化所有writer 然后注入biz
	idt, err := idempotent.NewIdempotent(rds)
	if err != nil {
		return nil, err
	}

	biz := biz.NewBiz(ds, writer, rds, idt)

	srv := service.NewSmsServerService(biz)
	impl := usercenter.NewUserCenter(c.UserCenterOptions, c.EtcdOptions)

	// gin.Recovery() 中间件，用来捕获任何 panic，并恢复
	mws := []gin.HandlerFunc{gin.Recovery(), header.NoCache, header.Cors, header.Secure,
		trace.TraceID(), auth.BasicAuth(impl), validate.Validator()}

	// 设置 Gin 模式
	gin.SetMode(gin.ReleaseMode)

	// 创建 Gin 引擎
	g := gin.New()

	// 添加中间件
	g.Use(mws...)

	// 并初始化路由
	// 这里注册不同的路由可以分开，如是否使用人认证中间件，分别在use 认证中间件前后
	installRouters(g, srv, c.Accounts)
	// 考虑在这里install consumer

	httpsrv, err := NewHTTPServer(c.HTTPOptions, c.TLSOptions, g)
	if err != nil {
		return nil, err
	}

	// todo grpc mq server
	//kqConf := kq.KqConf{}
	//_ = copier.Copy(kqConf, c.Config.KqOptions)
	//background := context.Background()
	//queue1 := kq.MustNewQueue(kqConf, consumer.NewArticleLikeNumBiz(background, idt))
	//var queues []zeroservice.Service
	//queues = append(queues, queue1)
	////kq.MustNewQueue(kqConf, consumer.NewArticleLikeNumLogic(background, idt)) // others
	//
	//mqsrv := zeroservice.NewServiceGroup()
	//for _, item := range queues {
	//	mqsrv.Add(item)
	//}

	return &SmsServer{grpcsrv: nil, httpsrv: httpsrv}, nil
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
