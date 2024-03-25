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

package openfga

import (
	"context"

	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	"github.com/openfga/openfga/pkg/tuple"
)

type model struct {
	id string
	s  *store
}

func (m *model) Check(ctx context.Context, object, relation, user string) (bool, error) {
	return m.CheckTuple(ctx, tuple.NewTupleKey(object, relation, user))
}

func (m *model) CheckTuple(ctx context.Context, key *openfgav1.TupleKey) (bool, error) {
	res, err := m.s.c.c.Check(ctx, &openfgav1.CheckRequest{
		StoreId:              m.s.id,
		AuthorizationModelId: m.id,
		TupleKey: &openfgav1.CheckRequestTupleKey{
			User:     key.User,
			Relation: key.Relation,
			Object:   key.Object,
		},
	})
	if err != nil {
		return false, err
	}
	return res.Allowed, nil
}

func (m *model) Write(ctx context.Context, object, relation, user string) error {
	return m.WriteTuples(ctx, tuple.NewTupleKey(object, relation, user))
}

func (m *model) Delete(ctx context.Context, object, relation, user string) error {
	return m.DeleteTuples(ctx, tuple.NewTupleKey(object, relation, user))
}

func (m *model) WriteTuples(ctx context.Context, keys ...*openfgav1.TupleKey) error {
	tx := m.Tx()
	defer tx.Close()
	if err := tx.WriteTuples(keys...); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (m *model) DeleteTuples(ctx context.Context, keys ...*openfgav1.TupleKey) error {
	tx := m.Tx()
	defer tx.Close()
	if err := tx.DeleteTuples(keys...); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (m *model) Tx() Tx {
	return &tx{m: m}
}

func (m *model) ID() string {
	return m.id
}

func (m *model) Store() Store {
	return m.s
}
