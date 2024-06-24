package template

import (
	"github.com/gin-gonic/gin"
	"github.com/superproj/onex/internal/pkg/core"
	v1 "github.com/superproj/onex/pkg/api/sms/v1"
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
