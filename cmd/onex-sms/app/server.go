// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Rosas99/smsx.
//

package app

import (
	genericapiserver "k8s.io/apiserver/pkg/server"

	"github.com/Rosas99/smsx/cmd/onex-sms/app/options"
	"github.com/Rosas99/smsx/internal/sms"
	"github.com/Rosas99/smsx/pkg/app"
)

const commandDesc = `The sms server is a standard, specification-compliant demo 
example of the onex service.

Find more onex-sms information at:
    https://github.com/Rosas99/smsx/blob/master/docs/guide/en-US/cmd/onex-sms.md`

// NewApp creates an App object with default parameters.
func NewApp(name string) *app.App {
	opts := options.NewOptions()
	application := app.NewApp(name, "Launch a onex sms server",
		app.WithDescription(commandDesc),
		app.WithOptions(opts),
		app.WithDefaultValidArgs(),
		app.WithRunFunc(run(opts)),
	)

	return application
}

func run(opts *options.Options) app.RunFunc {
	return func() error {
		cfg, err := opts.Config()
		if err != nil {
			return err
		}

		return Run(cfg, genericapiserver.SetupSignalHandler())
	}
}

// Run runs the specified APIServer. This should never exit.
func Run(c *sms.Config, stopCh <-chan struct{}) error {
	server, err := c.Complete().New()
	if err != nil {
		return err
	}

	return server.Run(stopCh)
}
