package template

import "github.com/superproj/onex/internal/sms/service"

type TemplateController struct {
	svc *service.SmsServerService
}

func New(svc *service.SmsServerService) *TemplateController {
	return &TemplateController{svc: svc}
}
