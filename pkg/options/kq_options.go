// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/onex.
//

package options

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/zeromicro/go-zero/core/service"
)

var _ IOptions = (*KqOptions)(nil)

// KafkaOptions defines options for kafka cluster.
// Common options for kafka-go reader and writer.
type KqOptions struct {
	service.ServiceConf
	Brokers     []string `mapstructure:"brokers"`
	Group       string   `mapstructure:"group"`
	Topic       string   `mapstructure:"topic"`
	Offset      string   `mapstructure:"offest" ,json:",options=first|last,default=last"` // 这里不做可填项的限制？
	Conns       int      `mapstructure:"conns" ,json:",default=1"`
	Consumers   int      `mapstructure:"consumers" ,json:",default=8"`
	Processors  int      `mapstructure:"processors" ,json:",default=8"`
	MinBytes    int      `mapstructure:"minBytes" ,json:",default=10240"`
	MaxBytes    int      `mapstructure:"maxBytes" ,json:",default=10485760"`
	Username    string   `mapstructure:"username" ,json:",optional"`
	Password    string   `mapstructure:"password" ,json:",optional"`
	ForceCommit bool     `mapstructure:"forceCommit" ,json:",default=true"`
}

func (k KqOptions) Validate() []error {
	errs := []error{}

	// todo确认校验方法
	if len(k.Brokers) == 0 {
		errs = append(errs, fmt.Errorf("kafka broker can not be empty"))
	}
	return errs

}

func (k KqOptions) AddFlags(fs *pflag.FlagSet, prefixes ...string) {

	//TODO implement me 暂时用不到
	panic("implement me")
}

// NewKafkaOptions create a `zero` value instance.
func NewKqOptions() *KqOptions {
	// todo 默认值在这里设置
	// 这里reader 和writer都有了
	return &KqOptions{
		Conns:       1,
		Consumers:   8,
		Processors:  8,
		MinBytes:    10240,
		MaxBytes:    10485760,
		ForceCommit: true,
	}
}
