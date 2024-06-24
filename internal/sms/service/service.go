// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/onex.
//

package service

import (
	"github.com/superproj/onex/internal/sms/biz"
)

type SmsServerService struct {
	biz biz.IBiz
}

func NewSmsServerService(biz biz.IBiz) *SmsServerService {
	return &SmsServerService{biz: biz}
}
