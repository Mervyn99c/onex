// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Rosas99/smsx.
//

package biz

//go:generate mockgen -destination mock_biz.go -package biz github.com/Rosas99/smsx/internal/fakeserver/biz IBiz

import (
	"github.com/Rosas99/smsx/internal/pkg/idempotent"
	"github.com/Rosas99/smsx/internal/sms"
	"github.com/Rosas99/smsx/internal/sms/biz/message"
	"github.com/Rosas99/smsx/internal/sms/biz/template"
	"github.com/Rosas99/smsx/internal/sms/store"
	"github.com/redis/go-redis/v9"
)

// IBiz 定义了 Biz 层需要实现的方法.
type IBiz interface {
	Templates() template.TemplateBiz
	Messages() message.MessageBiz
}

// biz 是 IBiz 的一个具体实现.
type Biz struct {
	ds    store.IStore
	kafka *sms.Logger
	rds   *redis.Client
	idt   *idempotent.Idempotent
}

// 确保 biz 实现了 IBiz 接口.
var _ IBiz = (*Biz)(nil)

// NewBiz 创建一个 IBiz 类型的实例.
func NewBiz(ds store.IStore, kafka *sms.Logger, rds *redis.Client, idt *idempotent.Idempotent) *Biz {
	return &Biz{ds: ds, kafka: kafka, rds: rds, idt: idt}
}

// Orders 返回一个实现了 OrderBiz 接口的实例.
func (b *Biz) Templates() template.TemplateBiz {
	message.New(b.ds, b.kafka, b.rds, b.idt)
	return template.New(b.ds, b.kafka, b.rds)
}

func (b *Biz) Messages() message.MessageBiz {
	return message.New(b.ds, b.kafka, b.rds, b.idt)
}
