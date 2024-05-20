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
	"fmt"

	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	parser "github.com/openfga/language/pkg/go/transformer"
	"github.com/openfga/openfga/pkg/tuple"
	"google.golang.org/protobuf/types/known/structpb"
)

type model struct {
	m *openfgav1.AuthorizationModel
	s *store
}

func (m *model) Check(ctx context.Context, object, relation, user string) (bool, error) {
	return m.CheckTuple(ctx, tuple.NewTupleKey(object, relation, user))
}

func (m *model) CheckWithContext(ctx context.Context, object, relation, user string, kv ...any) (bool, error) {
	return m.CheckTupleWithContext(ctx, tuple.NewTupleKey(object, relation, user), kv...)
}

func (m *model) CheckTuple(ctx context.Context, key *openfgav1.TupleKey) (bool, error) {
	res, err := m.s.c.c.Check(ctx, &openfgav1.CheckRequest{
		StoreId:              m.s.id,
		AuthorizationModelId: m.m.Id,
		TupleKey: &openfgav1.CheckRequestTupleKey{
			User:     key.User,
			Relation: key.Relation,
			Object:   key.Object,
		},
	})
	return res.GetAllowed(), err
}

func (m *model) CheckTupleWithContext(ctx context.Context, key *openfgav1.TupleKey, kv ...any) (bool, error) {
	c, err := makeContext(kv...)
	if err != nil {
		return false, err
	}
	res, err := m.s.c.c.Check(ctx, &openfgav1.CheckRequest{
		StoreId:              m.s.id,
		AuthorizationModelId: m.m.Id,
		TupleKey: &openfgav1.CheckRequestTupleKey{
			User:     key.User,
			Relation: key.Relation,
			Object:   key.Object,
		},
		Context: c,
	})
	return res.GetAllowed(), err
}

func (m *model) Read(ctx context.Context, object, relation, user string) ([]*openfgav1.Tuple, error) {
	res, err := m.s.c.c.Read(ctx, &openfgav1.ReadRequest{
		StoreId:  m.s.id,
		TupleKey: &openfgav1.ReadRequestTupleKey{User: user, Relation: relation, Object: object},
	})
	return res.GetTuples(), err
}

func (m *model) Expand(ctx context.Context, object, relation string) (*openfgav1.UsersetTree, error) {
	res, err := m.s.c.c.Expand(ctx, &openfgav1.ExpandRequest{
		StoreId:              m.s.id,
		TupleKey:             &openfgav1.ExpandRequestTupleKey{Relation: relation, Object: object},
		AuthorizationModelId: m.m.Id,
	})
	return res.GetTree(), err
}

func (m *model) List(ctx context.Context, typ, relation, user string) ([]string, error) {
	res, err := m.s.c.c.ListObjects(ctx, &openfgav1.ListObjectsRequest{
		StoreId:              m.s.id,
		AuthorizationModelId: m.m.Id,
		Type:                 typ,
		Relation:             relation,
		User:                 user,
	})
	return res.GetObjects(), err
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

func (m *model) WriteWithCondition(ctx context.Context, object, relation, user string, condition string, kv ...any) error {
	tx := m.Tx()
	defer tx.Close()
	if err := tx.WriteWithCondition(object, relation, user, condition, kv...); err != nil {
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
	return m.m.Id
}

func (m *model) Store() Store {
	return m.s
}

func (m *model) Show() (string, error) {
	return parser.TransformJSONProtoToDSL(m.m, parser.WithIncludeSourceInformation(true))
}

func makeContext(kv ...any) (*structpb.Struct, error) {
	m := make(map[string]any)
	for i := 0; i < len(kv); i += 2 {
		k, ok := kv[i].(string)
		if !ok {
			return nil, fmt.Errorf("invalid key %v: must be string", kv[i])
		}
		m[k] = kv[i+1]
	}
	return structpb.NewStruct(m)
}
