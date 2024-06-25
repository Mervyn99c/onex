package template

import "github.com/Rosas99/smsx/internal/sms/service"

type TemplateController struct {
	svc *service.SmsServerService
}

func New(svc *service.SmsServerService) *TemplateController {
	return &TemplateController{svc: svc}
}
