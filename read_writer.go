// Copyright 2025 Linka Cloud  All rights reserved.
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

	"github.com/google/uuid"
	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	"github.com/openfga/openfga/pkg/tuple"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type rw struct {
	c   openfgav1.OpenFGAServiceServer
	m   *openfgav1.AuthorizationModel
	mid string
	sid string
}

func (m *rw) Check(ctx context.Context, object, relation, user string, contextKVs ...any) (bool, error) {
	return m.CheckTuple(ctx, tuple.NewTupleKey(object, relation, user), contextKVs...)
}

func (m *rw) CheckTuple(ctx context.Context, key *openfgav1.TupleKey, contextKVs ...any) (bool, error) {
	c, err := makeContext(contextKVs...)
	if err != nil {
		return false, err
	}
	res, err := m.c.Check(ctx, &openfgav1.CheckRequest{
		StoreId:              m.sid,
		AuthorizationModelId: m.mid,
		TupleKey: &openfgav1.CheckRequestTupleKey{
			User:     key.User,
			Relation: key.Relation,
			Object:   key.Object,
		},
		Context: c,
	})
	return res.GetAllowed(), err
}

func (m *rw) Read(ctx context.Context, object, relation, user string) ([]*openfgav1.Tuple, error) {
	var key *openfgav1.ReadRequestTupleKey
	if object != "" || user != "" || relation != "" {
		key = &openfgav1.ReadRequestTupleKey{User: user, Relation: relation, Object: object}
	}
	res, err := m.c.Read(ctx, &openfgav1.ReadRequest{
		StoreId:  m.sid,
		TupleKey: key,
	})
	return res.GetTuples(), err
}

func (m *rw) ReadWithPaging(ctx context.Context, object, relation, user string, pageSize int32, continuationToken string) ([]*openfgav1.Tuple, string, error) {
	var key *openfgav1.ReadRequestTupleKey
	if object != "" || user != "" || relation != "" {
		key = &openfgav1.ReadRequestTupleKey{User: user, Relation: relation, Object: object}
	}
	res, err := m.c.Read(ctx, &openfgav1.ReadRequest{
		StoreId:           m.sid,
		TupleKey:          key,
		PageSize:          wrapperspb.Int32(pageSize),
		ContinuationToken: continuationToken,
	})
	return res.GetTuples(), res.GetContinuationToken(), err
}

func (m *rw) Expand(ctx context.Context, object, relation string) (*openfgav1.UsersetTree, error) {
	res, err := m.c.Expand(ctx, &openfgav1.ExpandRequest{
		StoreId:              m.sid,
		TupleKey:             &openfgav1.ExpandRequestTupleKey{Relation: relation, Object: object},
		AuthorizationModelId: m.mid,
	})
	return res.GetTree(), err
}

func (m *rw) ListObjects(ctx context.Context, typ, relation, user string) ([]string, error) {
	res, err := m.c.ListObjects(ctx, &openfgav1.ListObjectsRequest{
		StoreId:              m.sid,
		AuthorizationModelId: m.mid,
		Type:                 typ,
		Relation:             relation,
		User:                 user,
	})
	return res.GetObjects(), err
}

func (m *rw) ListUsers(ctx context.Context, object, relation, userTyp string, contextKVs ...any) ([]string, error) {
	c, err := makeContext(contextKVs...)
	if err != nil {
		return nil, err
	}
	typ, id := tuple.SplitObject(object)
	res, err := m.c.ListUsers(ctx, &openfgav1.ListUsersRequest{
		StoreId:              m.sid,
		AuthorizationModelId: m.mid,
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

func (m *rw) ListRelations(ctx context.Context, object, user string, relations ...string) ([]string, error) {
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
	res, err := m.c.BatchCheck(ctx, &openfgav1.BatchCheckRequest{
		StoreId:              m.sid,
		AuthorizationModelId: m.mid,
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

func (m *rw) Write(ctx context.Context, object, relation, user string) error {
	return m.WriteTuples(ctx, tuple.NewTupleKey(object, relation, user))
}

func (m *rw) Delete(ctx context.Context, object, relation, user string) error {
	return m.DeleteTuples(ctx, tuple.NewTupleKey(object, relation, user))
}

func (m *rw) WriteTuples(ctx context.Context, keys ...*openfgav1.TupleKey) error {
	_, err := m.c.Write(ctx, &openfgav1.WriteRequest{
		StoreId:              m.sid,
		AuthorizationModelId: m.mid,
		Writes:               &openfgav1.WriteRequestWrites{TupleKeys: keys},
	})
	return err
}

func (m *rw) WriteWithCondition(ctx context.Context, object, relation, user string, condition string, kv ...any) error {
	c, err := makeContext(kv...)
	if err != nil {
		return fmt.Errorf("invalid condition: %w", err)
	}

	return m.WriteTuples(ctx, tuple.NewTupleKeyWithCondition(object, relation, user, condition, c))
}

func (m *rw) DeleteTuples(ctx context.Context, keys ...*openfgav1.TupleKey) error {
	var d []*openfgav1.TupleKeyWithoutCondition
	for _, key := range keys {
		d = append(d, &openfgav1.TupleKeyWithoutCondition{
			Object:   key.Object,
			Relation: key.Relation,
			User:     key.User,
		})

	}
	_, err := m.c.Write(ctx, &openfgav1.WriteRequest{
		StoreId:              m.sid,
		AuthorizationModelId: m.mid,
		Deletes:              &openfgav1.WriteRequestDeletes{TupleKeys: d},
	})
	return err
}
