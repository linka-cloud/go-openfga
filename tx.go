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
	"sync"

	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	"github.com/openfga/openfga/pkg/tuple"
)

type tx struct {
	m       *model
	writes  []*openfgav1.TupleKey
	deletes []*openfgav1.TupleKeyWithoutCondition

	mu   sync.Mutex
	done bool
}

func (t *tx) Write(object, relation, user string) error {
	return t.WriteTuples(tuple.NewTupleKey(object, relation, user))
}

func (t *tx) WriteWithCondition(object, relation, user string, condition string, kv ...any) error {
	c, err := makeContext(kv...)
	if err != nil {
		return fmt.Errorf("invalid condition: %w", err)
	}
	return t.WriteTuples(tuple.NewTupleKeyWithCondition(object, relation, user, condition, c))
}

func (t *tx) WriteTuples(data ...*openfgav1.TupleKey) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.done {
		return errors.New("transaction closed")
	}
	t.writes = append(t.writes, data...)
	return nil
}

func (t *tx) Delete(object, relation, user string) error {
	return t.DeleteTuples(tuple.NewTupleKey(object, relation, user))
}

func (t *tx) DeleteTuples(deletes ...*openfgav1.TupleKey) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.done {
		return errors.New("transaction closed")
	}
	var d []*openfgav1.TupleKeyWithoutCondition
	for _, key := range deletes {
		d = append(d, &openfgav1.TupleKeyWithoutCondition{
			Object:   key.Object,
			Relation: key.Relation,
			User:     key.User,
		})

	}
	t.deletes = append(t.deletes, d...)
	return nil
}

func (t *tx) Commit(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.done {
		return errors.New("transaction closed")
	}
	var w *openfgav1.WriteRequestWrites
	if len(t.writes) != 0 {
		w = &openfgav1.WriteRequestWrites{TupleKeys: t.writes}
	}
	var d *openfgav1.WriteRequestDeletes
	if len(t.deletes) != 0 {
		d = &openfgav1.WriteRequestDeletes{TupleKeys: t.deletes}
	}
	_, err := t.m.s.c.c.Write(ctx, &openfgav1.WriteRequest{
		StoreId:              t.m.s.id,
		Writes:               w,
		Deletes:              d,
		AuthorizationModelId: t.m.m.Id,
	})
	return err
}

func (t *tx) Close() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.done = true
	t.writes = nil
	t.deletes = nil
}
