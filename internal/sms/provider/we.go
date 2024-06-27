package provider

import (
	"fmt"
	"github.com/Rosas99/smsx/internal/pump/types"
	"github.com/redis/go-redis/v9"
)

// WEProvider 结构体
type WEProvider struct {
	rds *redis.Client
}

// todo 依赖注入
func NewWEProvider(rds *redis.Client) *WEProvider {
	return &WEProvider{
		rds: rds,
	}
}

// Send 实现发送短信的方法
func (p *WEProvider) Send(request types.TemplateMsgRequest) (TemplateMsgResponse, error) {
	// 这里应该是调用微信的API发送短信的逻辑
	fmt.Printf("Sending message via WEProvider to %s\n", request.PhoneNumber)
	// 返回示例响应
	return TemplateMsgResponse{MessageID: "123456"}, nil
}
