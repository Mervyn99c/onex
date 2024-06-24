package consumer

import (
	"context"
	"github.com/superproj/onex/internal/pump"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/queue"
	"github.com/zeromicro/go-zero/core/service"
)

func Init(ctx context.Context, config kq.KqConf, s *pump.PreparedServer) queue.MessageQueue {
	return kq.MustNewQueue(config, NewArticleLikeNumLogic(ctx, s))
}

func Consumers(ctx context.Context, config kq.KqConf, s *pump.PreparedServer) []service.Service {
	return []service.Service{
		kq.MustNewQueue(config, NewArticleLikeNumLogic(ctx, s)),
		//kq.MustNewQueue(config, NewArticleLikeNumLogic(ctx, s)), // 其他消费类
	}
}
