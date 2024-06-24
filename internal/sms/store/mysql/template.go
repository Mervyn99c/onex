// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/onex.
//

package mysql

import (
	"context"
	"errors"
	"github.com/superproj/onex/internal/pkg/meta"
	"github.com/superproj/onex/internal/sms/model"
	"gorm.io/gorm"
)

type templates struct {
	db *gorm.DB
}

func newTemplates(db *gorm.DB) *templates {
	return &templates{db: db}
}

func (t *templates) Create(ctx context.Context, template *model.TemplateM) error {
	return t.db.Create(&template).Error
}

func (t *templates) Get(ctx context.Context, templateCode string) (*model.TemplateM, error) {
	var template model.TemplateM
	if err := t.db.Where("template_code = ?", templateCode).First(&template).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

func (t *templates) Update(ctx context.Context, template *model.TemplateM) error {
	return t.db.Save(&template).Error
}
func (t *templates) List(ctx context.Context, templateCode string, opts ...meta.ListOption) (count int64, ret []*model.TemplateM, err error) {

	options := meta.NewListOptions(opts...)
	if templateCode != "" {
		options.Filters["template_code"] = templateCode
	}
	// todo 对比 ucenter
	ans := t.db.
		Where(options.Filters).
		Offset(options.Offset).
		Limit(options.Limit).
		Order("id desc").
		Find(ret).
		Limit(-1).
		Count(&count)

	return count, ret, ans.Error
}
func (t *templates) Delete(ctx context.Context, id int64) error {
	err := t.db.Where("id = ?", id).Delete(&model.TemplateM{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return nil
}
