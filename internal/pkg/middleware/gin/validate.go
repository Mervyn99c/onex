// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/onex.
//

package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/superproj/onex/internal/pkg/core"
	known "github.com/superproj/onex/internal/pkg/known/toyblc"
	"github.com/superproj/onex/pkg/api/zerrors"
)

type validator interface {
	Validate() error
}

// IValidator defines methods to implement a custom validator.
//type IValidator interface {
//	Validate(ctx context.Context, rq any) error
//}

func Validator() gin.HandlerFunc {
	return func(c *gin.Context) {
		var rq any
		if err := c.ShouldBindJSON(&rq); err != nil {
		}

		if v, ok := rq.(validator); ok {
			// Kratos validation method
			if err := v.Validate(); err != nil {
				core.WriteResponse(c, nil, zerrors.ErrorInvalidParameter(err.Error()).WithCause(err))
				c.Abort()
				return
			}
		}

		c.Set(known.UsernameKey, "userID")
		c.Next()
	}

}

//func CustomValidator(vd IValidator) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		var rq any
//		if err := c.ShouldBindJSON(&rq); err != nil {
//		}
//
//		if v, ok := rq.(validator); ok {
//			// Kratos validation method
//			if err := v.Validate(); err != nil {
//				if se := new(errors.Error); errors.As(err, &se) {
//					core.WriteResponse(c, nil, se)
//					c.Abort()
//					return
//
//				}
//				core.WriteResponse(c, nil, zerrors.ErrorInvalidParameter(err.Error()).WithCause(err))
//				c.Abort()
//				return
//			}
//		}
//
//		//Custom validation, specific to the API interface
//		if err := vd.Validate(c, rq); err != nil {
//			if se := new(errors.Error); errors.As(err, &se) {
//				core.WriteResponse(c, nil, se)
//				c.Abort()
//				return
//			}
//
//			core.WriteResponse(c, nil, zerrors.ErrorInvalidParameter(err.Error()).WithCause(err))
//			c.Abort()
//			return
//		}
//
//		c.Set(known.UsernameKey, "userID")
//		c.Next()
//	}
//}
