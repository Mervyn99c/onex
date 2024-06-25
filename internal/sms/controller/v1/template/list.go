package template

import (
	"github.com/Rosas99/smsx/internal/pkg/core"
	v1 "github.com/Rosas99/smsx/pkg/api/sms/v1"
	"github.com/gin-gonic/gin"
)

func (b *TemplateController) List(c *gin.Context) {
	// query 示例
	//if err := c.ShouldBindQuery(&r); err != nil {
	//	core.WriteResponse(c, errors.WithCode(code.ErrBind, err.Error()), nil)
	//
	//	return
	//}

	var r v1.ListTemplateRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, err, nil)
	}

	template, err := b.svc.ListTemplate(c, &r)
	if err != nil {
		core.WriteResponse(c, err, nil)

	}
	core.WriteResponse(c, nil, template)

}
