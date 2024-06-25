// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Rosas99/smsx.
//

//go:generate wire .
package pump

import (
	"github.com/Rosas99/smsx/internal/pkg/idempotent"
	flow2 "github.com/Rosas99/smsx/internal/pump/mq"
	"github.com/Rosas99/smsx/internal/pump/provider"
	"github.com/Rosas99/smsx/internal/pump/types"
	"github.com/Rosas99/smsx/internal/sms/store"
	"github.com/Rosas99/smsx/internal/sms/store/mysql"
	"github.com/Rosas99/smsx/pkg/db"
	"github.com/Rosas99/smsx/pkg/streams/flow"
	"github.com/jinzhu/copier"
	"time"

	"github.com/segmentio/kafka-go"
	"k8s.io/apimachinery/pkg/util/wait"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Rosas99/smsx/pkg/log"
	genericoptions "github.com/Rosas99/smsx/pkg/options"
	kafkaconnector "github.com/Rosas99/smsx/pkg/streams/connector/kafka"
)

// Config defines the config for the apiserver.
type Config struct {
	KafkaOptions *genericoptions.KafkaOptions
	MongoOptions *genericoptions.MongoOptions
	RedisOptions *genericoptions.RedisOptions
	MySQLOptions *genericoptions.MySQLOptions
}

// Server contains state for a Kubernetes cluster master/api server.
type Server struct {
	kafkaReader kafka.ReaderConfig
	colName     string
	db          *mongo.Database
	idt         *idempotent.Idempotent
	logger      *Logger
}

type completedConfig struct {
	*Config
}

// 思路1 函数使用依赖注入redi
// 思路2 函数修改为方法，但是map方法需要使用适配器兼容方法参数
// addUTC appends a UTC timestamp to the beginning of the message value.
var addUTC = func(msg kafka.Message) kafka.Message {
	timestamp := time.Now().Format(time.DateTime)

	// Concatenate the UTC timestamp with msg.Value
	msg.Value = []byte(timestamp + " " + string(msg.Value))
	// 这里拿到值后可以写入MySQL MongoDB，做日志记录 不一定是byte类型的
	return msg
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (cfg *Config) Complete() completedConfig {
	return completedConfig{cfg}
}

// New returns a new instance of Server from the given config.
// Certain config fields will be set to a default value if unset.
func (c completedConfig) New() (*Server, error) {
	client, err := c.MongoOptions.NewClient()
	if err != nil {
		return nil, err
	}

	var redisOptions db.RedisOptions
	value := &c.Config.RedisOptions
	_ = copier.Copy(&redisOptions, value)
	rds, err := db.NewRedis(&redisOptions)
	if err != nil {
		return nil, err
	}

	//这里初始化所有writer 然后注入biz
	idt, err := idempotent.NewIdempotent(rds)
	if err != nil {
		return nil, err
	}
	factory := provider.NewProviderFactory()
	factory.RegisterProvider(types.ProviderTypeWE, provider.NewWEProvider(rds))

	var dbOptions db.MySQLOptions
	_ = copier.Copy(&dbOptions, c.MySQLOptions)
	ins, err := db.NewMySQL(&dbOptions)
	if err != nil {
		return nil, err
	}
	var ds store.IStore
	ds = mysql.NewStore(ins)
	logger, err := NewLogger(ds.Templates())
	if err != nil {
		return nil, err
	}

	server := &Server{
		// 这里带有默认值的可以不配置
		kafkaReader: kafka.ReaderConfig{
			Brokers:           c.KafkaOptions.Brokers,
			Topic:             c.KafkaOptions.Topic,
			GroupID:           c.KafkaOptions.ReaderOptions.GroupID,
			QueueCapacity:     c.KafkaOptions.ReaderOptions.QueueCapacity,
			MinBytes:          c.KafkaOptions.ReaderOptions.MinBytes,
			MaxBytes:          c.KafkaOptions.ReaderOptions.MaxBytes,
			MaxWait:           c.KafkaOptions.ReaderOptions.MaxWait,
			ReadBatchTimeout:  c.KafkaOptions.ReaderOptions.ReadBatchTimeout,
			HeartbeatInterval: c.KafkaOptions.ReaderOptions.HeartbeatInterval,
			CommitInterval:    c.KafkaOptions.ReaderOptions.CommitInterval,
			RebalanceTimeout:  c.KafkaOptions.ReaderOptions.RebalanceTimeout,
			StartOffset:       c.KafkaOptions.ReaderOptions.StartOffset,
			MaxAttempts:       c.KafkaOptions.ReaderOptions.MaxAttempts,
		},
		colName: c.MongoOptions.Collection,
		db:      client.Database(c.MongoOptions.Database),
		// todo redis
		idt:    idt,
		logger: logger,
	}

	return server, nil
}

type PreparedServer struct {
	*Server
}

func (s *Server) PrepareRun() PreparedServer {
	return PreparedServer{s}
}

func (s PreparedServer) Run(stopCh <-chan struct{}) error {
	// 支持多个消费者

	ctx := wait.ContextForChannel(stopCh)
	// todo reader已经提供了默认值

	log.Infof("Successfully start pump server")

	// todo reader已经提供了默认值
	source2, err := kafkaconnector.NewKafkaSource(ctx, s.kafkaReader)
	if err != nil {
		return err
	}

	// todo 这里可以传入整个server，也可以只传入redis store等
	logic := flow2.NewHandlerMessageBiz(ctx, s.db, s.idt, s.logger)
	articleConsumer := flow.NewConsumer(logic, 1)
	// 这里通过map写入通道，通道是由sink初始化后开始消费
	source2.Via(articleConsumer)
	return err
}
