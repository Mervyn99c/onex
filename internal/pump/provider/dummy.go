package provider

import (
	"fmt"
	"github.com/Rosas99/smsx/internal/pump/types"
)

// DummyProvider 结构体
type DummyProvider struct{}

func NewDummyProvider() *DummyProvider {
	return &DummyProvider{}
}

func (p *DummyProvider) Send(request types.TemplateMsgRequest) (TemplateMsgResponse, error) {
	// 模拟发送短信的逻辑，不实际发送
	fmt.Printf("Simulating message send via DummyProvider to %s\n", request.PhoneNumber)
	return TemplateMsgResponse{}, nil
}
