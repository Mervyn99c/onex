package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/superproj/onex/internal/pump/types"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/mongo"
	"strconv"
)

type ArticleLikeNumLogic struct {
	ctx context.Context
	db  *mongo.Database
	logx.Logger
}

func NewArticleLikeNumLogic(ctx context.Context, db *mongo.Database) *ArticleLikeNumLogic {
	return &ArticleLikeNumLogic{
		ctx:    ctx,
		db:     db,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ArticleLikeNumLogic) Consume(elem any) error {
	val := elem.(kafka.Message)

	var msg *types.CanalLikeMsg
	err := json.Unmarshal(val.Value, &msg)
	if err != nil {
		logx.Errorf("Consume val: %s error: %v", val, err)
		return err
	}

	return l.updateArticleLikeNum(l.ctx, msg)
}

func (l *ArticleLikeNumLogic) updateArticleLikeNum(ctx context.Context, msg *types.CanalLikeMsg) error {
	// 幂等检查
	// 如果是短信验证码，kpi记录时有所不同，其他不变
	// 供应商调用，失败则尝试下一个
	// 调用第三方接口使用什么方式
	// 错误处理 响应
	// 多线程消费
	db := l.db
	db.Name()

	if msg.BizID != types.ArticleBizID {
	}
	id, err := strconv.ParseInt(msg.ObjID, 10, 64)
	fmt.Println(id)
	if err != nil {
		logx.Errorf("strconv.ParseInt id: %s error: %v", msg.ID, err)
	}
	likeNum, err := strconv.ParseInt(msg.LikeNum, 10, 64)
	if err != nil {
		logx.Errorf("strconv.ParseInt likeNum: %s error: %v", msg.LikeNum, err)
	}
	fmt.Print(likeNum)
	//err = l.svcCtx.ArticleModel.UpdateLikeNum(ctx, id, likeNum)
	//if err != nil {
	//	logx.Errorf("UpdateLikeNum id: %d like: %d", id, likeNum)
	//}

	return nil
}
