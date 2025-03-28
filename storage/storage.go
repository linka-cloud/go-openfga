// Copyright 2024 Linka Cloud  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package storage

import (
	"context"

	"github.com/openfga/openfga/pkg/storage"
)

type TxOptions struct {
	ReadOnly bool
}

type TxOption func(*TxOptions)

func WithReadOnly() TxOption {
	return func(o *TxOptions) {
		o.ReadOnly = true
	}
}

type Tx[T any] interface {
	storage.TupleBackend
	storage.AuthorizationModelBackend
	storage.StoresBackend
	storage.AssertionsBackend
	Commit(ctx context.Context) error
	Close()
	Unwrap() T
}

type TxProvider[T any] interface {
	Tx(ctx context.Context, opts ...TxOption) (Tx[T], error)
	WithTx(tx T) Tx[T]
}

type Datastore[T any] interface {
	storage.OpenFGADatastore
	TxProvider[T]
}

type txKey struct{}

func Context[T any](ctx context.Context, tx Tx[T]) context.Context {
	return context.WithValue(ctx, txKey{}, tx.Unwrap())
}

func From[T any](ctx context.Context) (T, bool) {
	tx, ok := ctx.Value(txKey{}).(T)
	return tx, ok
}
