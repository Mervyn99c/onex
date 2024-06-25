// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Rosas99/smsx.
//

package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Rosas99/smsx/internal/pkg/zid"
)

func (t *TemplateM) BeforeCreate(tx *gorm.DB) (err error) {
	if t.TemplateCode == "" {
		// Generate a new UUID for SecretKey.
		t.TemplateCode = uuid.New().String()
	}

	// 可以考虑zid的生成方式
	return nil
}

// AfterCreate todo 待确认 生成带盐值的随机id 非主键
func (t *TemplateM) AfterCreate(tx *gorm.DB) (err error) {
	t.TemplateCode = zid.Template.New(uint64(t.ID)) // Generate and set a new order ID.

	return tx.Save(t).Error // Save the updated order record to the database.
}
