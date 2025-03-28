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
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	parser "github.com/openfga/language/pkg/go/transformer"
	"github.com/openfga/openfga/pkg/tuple"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const SchemaVersion = "1.2"

type model struct {
	m  *openfgav1.AuthorizationModel
	s  *store
	mu sync.RWMutex
}

func (m *model) Check(ctx context.Context, object, relation, user string, contextKVs ...any) (bool, error) {
	return m.CheckTuple(ctx, tuple.NewTupleKey(object, relation, user), contextKVs...)
}

func (m *model) CheckTuple(ctx context.Context, key *openfgav1.TupleKey, contextKVs ...any) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, err := makeContext(contextKVs...)
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
	m.mu.RLock()
	defer m.mu.RUnlock()
	var key *openfgav1.ReadRequestTupleKey
	if object != "" || user != "" || relation != "" {
		key = &openfgav1.ReadRequestTupleKey{User: user, Relation: relation, Object: object}
	}
	res, err := m.s.c.c.Read(ctx, &openfgav1.ReadRequest{
		StoreId:  m.s.id,
		TupleKey: key,
	})
	return res.GetTuples(), err
}

func (m *model) ReadWithPaging(ctx context.Context, object, relation, user string, pageSize int32, continuationToken string) ([]*openfgav1.Tuple, string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var key *openfgav1.ReadRequestTupleKey
	if object != "" || user != "" || relation != "" {
		key = &openfgav1.ReadRequestTupleKey{User: user, Relation: relation, Object: object}
	}
	res, err := m.s.c.c.Read(ctx, &openfgav1.ReadRequest{
		StoreId:           m.s.id,
		TupleKey:          key,
		PageSize:          wrapperspb.Int32(pageSize),
		ContinuationToken: continuationToken,
	})
	return res.GetTuples(), res.GetContinuationToken(), err
}

func (m *model) Expand(ctx context.Context, object, relation string) (*openfgav1.UsersetTree, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res, err := m.s.c.c.Expand(ctx, &openfgav1.ExpandRequest{
		StoreId:              m.s.id,
		TupleKey:             &openfgav1.ExpandRequestTupleKey{Relation: relation, Object: object},
		AuthorizationModelId: m.m.Id,
	})
	return res.GetTree(), err
}

func (m *model) ListObjects(ctx context.Context, typ, relation, user string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res, err := m.s.c.c.ListObjects(ctx, &openfgav1.ListObjectsRequest{
		StoreId:              m.s.id,
		AuthorizationModelId: m.m.Id,
		Type:                 typ,
		Relation:             relation,
		User:                 user,
	})
	return res.GetObjects(), err
}

func (m *model) ListUsers(ctx context.Context, object, relation, userTyp string, contextKVs ...any) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, err := makeContext(contextKVs...)
	if err != nil {
		return nil, err
	}
	typ, id := tuple.SplitObject(object)
	res, err := m.s.c.c.ListUsers(ctx, &openfgav1.ListUsersRequest{
		StoreId:              m.s.id,
		AuthorizationModelId: m.m.Id,
		Relation:             relation,
		UserFilters:          []*openfgav1.UserTypeFilter{{Type: userTyp}},
		Object: &openfgav1.Object{
			Type: typ,
			Id:   id,
		},
		Context: c,
	})
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(res.GetUsers()))
	for _, u := range res.GetUsers() {
		out = append(out, tuple.UserProtoToString(u))
	}
	return out, err
}

func (m *model) ListRelations(ctx context.Context, object, user string, relations ...string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if len(relations) == 0 {
		objectType := strings.Split(object, ":")[0]
		for _, v := range m.m.TypeDefinitions {
			if v.Type != objectType {
				continue
			}
			for relation := range v.Relations {
				relations = append(relations, relation)
			}
			break
		}
	}
	if len(relations) == 0 {
		return nil, nil
	}
	var items []*openfgav1.BatchCheckItem
	var keys []string
	for _, v := range relations {
		k := uuid.NewString()
		keys = append(keys, k)
		items = append(items, &openfgav1.BatchCheckItem{
			CorrelationId: k,
			TupleKey: &openfgav1.CheckRequestTupleKey{
				User:     user,
				Relation: v,
				Object:   object,
			},
		})
	}
	res, err := m.s.c.c.BatchCheck(ctx, &openfgav1.BatchCheckRequest{
		StoreId:              m.s.id,
		AuthorizationModelId: m.m.Id,
		Checks:               items,
	})
	if err != nil {
		return nil, err
	}
	if len(res.Result) != len(relations) {
		return nil, errors.New("invalid response: relation length mismatch")
	}
	out := make([]string, 0, len(relations))
	for i, v := range keys {
		vv, ok := res.Result[v]
		if !ok {
			return nil, errors.New("invalid response: missing batch check result")
		}
		if vv.GetAllowed() {
			out = append(out, relations[i])
		}
	}
	return out, nil
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
	m.mu.RLock()
	defer m.mu.RUnlock()
	tx := m.Tx()
	defer tx.Close()
	if err := tx.WriteWithCondition(object, relation, user, condition, kv...); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (m *model) DeleteTuples(ctx context.Context, keys ...*openfgav1.TupleKey) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	tx := m.Tx()
	defer tx.Close()
	if err := tx.DeleteTuples(keys...); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (m *model) Reload(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	res, err := m.s.c.c.ReadAuthorizationModels(ctx, &openfgav1.ReadAuthorizationModelsRequest{StoreId: m.s.id})
	if err != nil {
		return err
	}
	if len(res.AuthorizationModels) == 0 {
		return errors.New("not found")
	}
	m.m = res.AuthorizationModels[len(res.AuthorizationModels)-1]
	return nil
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

func CombineModules(dsl ...string) (string, error) {
	mod, err := combineModules(dsl...)
	if err != nil {
		return "", err
	}
	// parser.TransformJSONProtoToDSL fails due to a bug in the parser on the nil `this`
	b, err := protojson.Marshal(mod)
	if err != nil {
		return "", err
	}
	mod.Reset()
	if err := protojson.Unmarshal(b, mod); err != nil {
		return "", err
	}
	out, err := parser.TransformJSONProtoToDSL(mod, parser.WithIncludeSourceInformation(true))
	if err != nil {
		return "", err
	}
	return out, nil
}

func combineModules(dsl ...string) (*openfgav1.AuthorizationModel, error) {
	if len(dsl) == 1 {
		return parser.TransformDSLToProto(dsl[0])
	}
	var mods []parser.ModuleFile
	for _, v := range dsl {
		mods = append(mods, parser.ModuleFile{Contents: v})
	}
	return parser.TransformModuleFilesToModel(mods, SchemaVersion)
}
