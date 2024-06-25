package service

import (
	"context"
	v1 "github.com/Rosas99/smsx/pkg/api/sms/v1"
	"github.com/Rosas99/smsx/pkg/log"
)

// todo 这里模板和配置关联 限流配置不用单独一个CURD

func (s *SmsServerService) Send(ctx context.Context, rq *v1.CreateTemplateRequest) (*v1.CreateTemplateResponse, error) {
	log.C(ctx).Infow("CreateOrder function called")
	return s.biz.Messages().Send(ctx, rq)
}
