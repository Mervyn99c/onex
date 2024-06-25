package mq

import (
	"context"
	"encoding/json"
	"github.com/Rosas99/smsx/internal/pkg/idempotent"
	"github.com/Rosas99/smsx/internal/pump"
	factory "github.com/Rosas99/smsx/internal/pump/provider"
	"github.com/Rosas99/smsx/internal/pump/types"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/mongo"
)

type HandleMessageBiz struct {
	ctx    context.Context
	db     *mongo.Database
	idt    *idempotent.Idempotent
	logger *pump.Logger
}

func NewHandlerMessageBiz(ctx context.Context, db *mongo.Database, idt *idempotent.Idempotent, logger *pump.Logger) *HandleMessageBiz {
	return &HandleMessageBiz{
		ctx:    ctx,
		db:     db,
		idt:    idt,
		logger: logger,
	}
}

func (l *HandleMessageBiz) Consume(elem any) error {
	val := elem.(kafka.Message)

	var msg *types.CanalLikeMsg
	err := json.Unmarshal(val.Value, &msg)
	if err != nil {
		logx.Errorf("Consume val: %s error: %v", val, err)
		return err
	}

	return l.handleSmsRequest(l.ctx, msg)
}

func (l *HandleMessageBiz) handleSmsRequest(ctx context.Context, msg *types.CanalLikeMsg) error {

	// 消息id
	ok := l.idt.Check(ctx, msg.BizID)
	if !ok {
		// 消费失败
	}

	db := l.db
	db.Name()
	tc := types.TemplateMsgRequest{}
	// 日志打印，手机号脱敏
	providers := tc.Providers
	for _, provider := range providers {
		providerFactory := factory.NewProviderFactory()
		templateProvider, err := providerFactory.GetSMSTemplateProvider(types.ProviderType(provider))
		if err != nil {
			break
		}
		_, err = templateProvider.Send(tc)
		if err != nil {
			continue
		}
		break
	}

	// todo 记录到history
	l.logger.LogHistory("")
	return nil
}
