// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/onex.
//

package store

import (
	"sync"

	"github.com/google/wire"

	ucstore "github.com/superproj/onex/internal/usercenter/store"
)

// ProviderSet is store providers.
var ProviderSet = wire.NewSet(NewStore, wire.Bind(new(Interface), new(*datastore)))

var (
	once sync.Once
	S    *datastore
)

// Interface defines the storage interface.
type Interface interface {
	UserCenter() ucstore.IStore
}

type datastore struct {
	uc ucstore.IStore
}

var _ Interface = (*datastore)(nil)

func (ds *datastore) UserCenter() ucstore.IStore {
	return ds.uc
}

func NewStore(uc ucstore.IStore) *datastore {
	once.Do(func() {
		S = &datastore{uc: uc}
	})

	return S
}
