// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/onex.
//

// Package options contains flags and options for initializing an apiserver
package options

import (
	"github.com/superproj/onex/internal/sms"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	cliflag "k8s.io/component-base/cli/flag"

	"github.com/superproj/onex/internal/pkg/feature"
	"github.com/superproj/onex/pkg/app"
	"github.com/superproj/onex/pkg/log"
	genericoptions "github.com/superproj/onex/pkg/options"
)

const (
	// UserAgent is the userAgent name when starting onex-sms server.
	UserAgent = "onex-sms"
)

var _ app.CliOptions = (*Options)(nil)

// Options contains state for master/api server.
type Options struct {
	//GRPCOptions  *genericoptions.GRPCOptions    `json:"grpc" mapstructure:"grpc"`
	HTTPOptions *genericoptions.HTTPOptions `json:"http" mapstructure:"http"`
	//TLSOptions   *genericoptions.TLSOptions     `json:"tls" mapstructure:"tls"`
	MySQLOptions *genericoptions.MySQLOptions   `json:"mysql" mapstructure:"mysql"`
	Metrics      *genericoptions.MetricsOptions `json:"metrics" mapstructure:"metrics"`
	//Redis options for configuring Redis related options.
	RedisOptions *genericoptions.RedisOptions `json:"redis" mapstructure:"redis"`
	// Kafka options for configuring Kafka related options.
	KafkaOptions1 *genericoptions.KafkaOptions `json:"kafka1" mapstructure:"kafka"`
	KqOptions1    *genericoptions.KqOptions    `json:"kqConf" mapstructure:"kqConf"`
	// Kafka options for configuring Kafka related options.
	//KafkaOptions2 *genericoptions.KafkaOptions `json:"kafka2" mapstructure:"kafka2"`
	Log *log.Options `json:"log" mapstructure:"log"`

	// Path to kubeconfig file with authorization and master location information.
	Kubeconfig   string          `json:"kubeconfig" mapstructure:"kubeconfig"`
	FeatureGates map[string]bool `json:"feature-gates" mapstructure:"-"`
}

// NewOptions returns initialized Options.
func NewOptions() *Options {
	o := &Options{
		//GRPCOptions:   genericoptions.NewGRPCOptions(),
		HTTPOptions: genericoptions.NewHTTPOptions(),
		//TLSOptions:    genericoptions.NewTLSOptions(),
		MySQLOptions:  genericoptions.NewMySQLOptions(),
		Metrics:       genericoptions.NewMetricsOptions(),
		RedisOptions:  genericoptions.NewRedisOptions(),
		KafkaOptions1: genericoptions.NewKafkaOptions(),
		KqOptions1:    genericoptions.NewKqOptions(),
		//KafkaOptions2: genericoptions.NewKafkaOptions(),
		Log: log.NewOptions(),
	}

	return o
}

// Flags returns flags for a specific server by section name.
func (o *Options) Flags() (fss cliflag.NamedFlagSets) {
	//o.GRPCOptions.AddFlags(fss.FlagSet("grpc"))
	o.HTTPOptions.AddFlags(fss.FlagSet("http"))
	//o.TLSOptions.AddFlags(fss.FlagSet("tls"))
	o.MySQLOptions.AddFlags(fss.FlagSet("mysql"))
	o.Metrics.AddFlags(fss.FlagSet("metrics"))
	o.Log.AddFlags(fss.FlagSet("log"))
	o.RedisOptions.AddFlags(fss.FlagSet("redis"))
	//o.KafkaOptions1.AddFlags(fss.FlagSet("kafka"))
	//o.KafkaOptions2.AddFlags(fss.FlagSet("kafka"))
	// Note: the weird ""+ in below lines seems to be the only way to get gofmt to
	// arrange these text blocks sensibly. Grrr.
	fs := fss.FlagSet("misc")
	feature.DefaultMutableFeatureGate.AddFlag(fs)

	return fss
}

// Complete completes all the required options.
func (o *Options) Complete() error {

	_ = feature.DefaultMutableFeatureGate.SetFromMap(o.FeatureGates)
	return nil
}

// Validate validates all the required options.
func (o *Options) Validate() error {
	errs := []error{}

	//errs = append(errs, o.GRPCOptions.Validate()...)
	errs = append(errs, o.HTTPOptions.Validate()...)
	//errs = append(errs, o.TLSOptions.Validate()...)
	errs = append(errs, o.MySQLOptions.Validate()...)
	errs = append(errs, o.Metrics.Validate()...)
	errs = append(errs, o.Log.Validate()...)
	errs = append(errs, o.RedisOptions.Validate()...)
	//errs = append(errs, o.KafkaOptions1.Validate()...)
	//errs = append(errs, o.KqOptions2.Validate()...)
	return utilerrors.NewAggregate(errs)
}

// ApplyTo fills up onex-fakeserver config with options.
func (o *Options) ApplyTo(c *sms.Config) error {
	//c.GRPCOptions = o.GRPCOptions
	c.HTTPOptions = o.HTTPOptions
	//c.TLSOptions = o.TLSOptions
	c.MySQLOptions = o.MySQLOptions
	c.RedisOptions = o.RedisOptions
	c.KafkaOptions1 = o.KafkaOptions1
	//c.KqOptions = o.KqOptions
	return nil
}

// Config return a onex-fakeserver config object.
func (o *Options) Config() (*sms.Config, error) {
	c := &sms.Config{}

	if err := o.ApplyTo(c); err != nil {
		return nil, err
	}

	return c, nil
}
