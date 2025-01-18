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

package protodb

import (
	"context"

	"go.linka.cloud/protodb"

	"go.linka.cloud/go-openfga/storage"
)

func maybeTx(ctx context.Context, p *pdb, opts ...protodb.TxOption) (protodb.Tx, error) {
	if tx, ok := storage.From[protodb.Tx](ctx); ok {
		return tx, nil
	}
	return p.db.Tx(ctx, opts...)
}

func maybeCommit(ctx context.Context, tx protodb.Tx) error {
	if _, ok := storage.From[protodb.Tx](ctx); ok {
		return nil
	}
	return tx.Commit(ctx)
}

func maybeClose(ctx context.Context, tx protodb.Tx) func() {
	if _, ok := storage.From[protodb.Tx](ctx); ok {
		return func() {}
	}
	return tx.Close
}

func withTx(ctx context.Context, p *pdb, fn func(ctx context.Context, txn *tx) error, opts ...protodb.TxOption) error {
	txn, err := maybeTx(ctx, p, opts...)
	if err != nil {
		return err
	}
	defer maybeClose(ctx, txn)()
	if err = fn(ctx, &tx{tx: txn}); err != nil {
		return err
	}
	if err := maybeCommit(ctx, txn); err != nil {
		return err
	}
	return nil
}

func withTx2[R any](ctx context.Context, p *pdb, fn func(ctx context.Context, txn *tx) (R, error), opts ...protodb.TxOption) (R, error) {
	var o1 R
	txn, err := maybeTx(ctx, p, opts...)
	if err != nil {
		return o1, err
	}
	defer maybeClose(ctx, txn)()
	o1, err = fn(ctx, &tx{tx: txn})
	if err != nil {
		return o1, err
	}
	if err := maybeCommit(ctx, txn); err != nil {
		return o1, err
	}
	return o1, nil
}

func withTx3[R any, T any](ctx context.Context, p *pdb, fn func(ctx context.Context, txn *tx) (R, T, error), opts ...protodb.TxOption) (R, T, error) {
	var (
		o1 R
		o2 T
	)
	txn, err := maybeTx(ctx, p, opts...)
	if err != nil {
		return o1, o2, err
	}
	defer maybeClose(ctx, txn)()
	o1, o2, err = fn(ctx, &tx{tx: txn})
	if err != nil {
		return o1, o2, err
	}
	if err := maybeCommit(ctx, txn); err != nil {
		return o1, o2, err
	}
	return o1, o2, nil
}
