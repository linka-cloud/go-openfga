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
	"fmt"
	"strconv"
	"strings"
	"time"

	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	"github.com/openfga/openfga/pkg/storage"
	"github.com/openfga/openfga/pkg/tuple"
	"go.linka.cloud/protodb"
	"go.linka.cloud/protodb/typed"
	"go.linka.cloud/protofilters/filters"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pbv1 "go.linka.cloud/go-openfga/storage/protodb/pb/v1"
)

var _ storage.OpenFGADatastore = (*pdb)(nil)

func New(ctx context.Context, option ...protodb.Option) (storage.OpenFGADatastore, error) {
	db, err := protodb.Open(ctx, option...)
	if err != nil {
		return nil, err
	}
	return newWithClient(ctx, db, false)

}

func NewWithClient(ctx context.Context, db protodb.Client) (storage.OpenFGADatastore, error) {
	return newWithClient(ctx, db, true)
}

func newWithClient(ctx context.Context, db protodb.Client, external bool) (storage.OpenFGADatastore, error) {
	for _, v := range []proto.Message{&openfgav1.Store{}, &pbv1.Tuple{}, &pbv1.Assertions{}, &pbv1.Change{}} {
		if err := db.Register(ctx, v.ProtoReflect().Descriptor().ParentFile()); err != nil {
			return nil, err
		}
	}
	return &pdb{db: db, external: external}, nil
}

type pdb struct {
	db       protodb.Client
	external bool
}

func (p *pdb) Read(ctx context.Context, store string, tupleKey *openfgav1.TupleKey, _ storage.ReadOptions) (storage.TupleIterator, error) {
	ts, _, err := typed.NewStore[pbv1.Tuple](p.db).Get(ctx, &pbv1.Tuple{}, protodb.WithFilter(makeFilter(store, tupleKey.User, tupleKey.Relation, tupleKey.Object)))
	if err != nil {
		return nil, err
	}
	return storage.NewStaticTupleIterator(toTuple(ts)), nil
}

func (p *pdb) ReadPage(ctx context.Context, store string, tupleKey *openfgav1.TupleKey, opts storage.ReadPageOptions) ([]*openfgav1.Tuple, []byte, error) {
	var typ string
	if tupleKey != nil {
		typ, _ = tuple.SplitObject(tupleKey.Object)
	}
	tk, err := parseToken(opts.Pagination.From, typ)
	if err != nil {
		return nil, nil, err
	}
	o := []protodb.GetOption{protodb.WithPaging(&protodb.Paging{Offset: tk.Offset, Limit: uint64(opts.Pagination.PageSize), Token: tk.Continuation})}
	if tupleKey != nil {
		o = append(o, protodb.WithFilter(makeFilter(store, tupleKey.User, tupleKey.Relation, tupleKey.Object)))
	} else {
		o = append(o, protodb.WithFilter(protodb.Where("key").StringHasPrefix(pbv1.TuplePrefix(store))))
	}
	ts, info, err := typed.NewStore[pbv1.Tuple](p.db).Get(ctx, &pbv1.Tuple{}, o...)
	if err != nil {
		return nil, nil, err
	}
	if !info.HasNext {
		return toTuple(ts), nil, nil
	}
	return toTuple(ts), newToken(info.Token, tk.Offset+uint64(len(ts)), typ).Bytes(), nil
}

func (p *pdb) ReadUserTuple(ctx context.Context, store string, tupleKey *openfgav1.TupleKey, _ storage.ReadUserTupleOptions) (*openfgav1.Tuple, error) {
	ts, _, err := typed.NewStore[pbv1.Tuple](p.db).Get(ctx, &pbv1.Tuple{Key: pbv1.TupleKey(store, tupleKey)})
	if err != nil {
		return nil, err
	}
	if len(ts) == 0 {
		return nil, storage.ErrNotFound
	}
	return &openfgav1.Tuple{Key: ts[0].TupleKey, Timestamp: ts[0].CreatedAt}, nil
}

func (p *pdb) ReadUsersetTuples(ctx context.Context, store string, filter storage.ReadUsersetTuplesFilter, _ storage.ReadUsersetTuplesOptions) (storage.TupleIterator, error) {
	prefix := fmt.Sprintf("%s/%s#%s@", store, filter.Object, filter.Relation)
	var f filters.Builder
	if len(filter.AllowedUserTypeRestrictions) != 0 {
		for _, v := range filter.AllowedUserTypeRestrictions {
			f2 := protodb.Where("key").StringHasPrefix(fmt.Sprintf("%s%s:", prefix, v.Type))
			if v.GetRelation() != "" {
				f2.And("key").StringHasSuffix(fmt.Sprintf("#%s", v.GetRelation()))
			}
			if f == nil {
				f = f2
			} else {
				f.Or(f2)
			}
		}
	} else {
		f = protodb.Where("key").StringHasPrefix(prefix)
	}
	ts, _, err := typed.NewStore[pbv1.Tuple](p.db).Get(ctx, &pbv1.Tuple{}, protodb.WithFilter(f))
	if err != nil {
		return nil, err
	}
	return storage.NewStaticTupleIterator(toTuple(ts)), nil
}

func (p *pdb) ReadStartingWithUser(ctx context.Context, store string, filter storage.ReadStartingWithUserFilter, _ storage.ReadStartingWithUserOptions) (storage.TupleIterator, error) {
	// TODO(adphi): handle filter.ObjectIDs
	var of []filters.Builder
	if filter.ObjectIDs != nil {
		for _, v := range filter.ObjectIDs.Values() {
			of = append(of, protodb.Where("key").StringHasPrefix(fmt.Sprintf("%s/%s:%s#", store, filter.ObjectType, v)))
		}
	} else {
		of = append(of, protodb.Where("key").StringHasPrefix(fmt.Sprintf("%s/%s:", store, filter.ObjectType)))
	}
	var f filters.Builder
	for _, v := range filter.UserFilter {
		for _, o := range of {
			pf := o.Clone()
			s := fmt.Sprintf("#%s@%s", filter.Relation, v.Object)
			if v.Relation != "" {
				s = fmt.Sprintf("#%s@%s#%s", filter.Relation, v.Object, v.Relation)
			}
			if f == nil {
				f = pf.And("key").StringHasSuffix(s)
			} else {
				f.Or(pf.And("key").StringHasSuffix(s))
			}
		}
	}
	ts, _, err := typed.NewStore[pbv1.Tuple](p.db).Get(ctx, &pbv1.Tuple{}, protodb.WithFilter(f))
	if err != nil {
		return nil, err
	}
	return storage.NewStaticTupleIterator(toTuple(ts)), nil
}

func (p *pdb) Write(ctx context.Context, store string, d storage.Deletes, w storage.Writes) error {
	return typed.WithTypedTx[pbv1.Tuple](ctx, p.db, func(ctx context.Context, tx typed.Tx[pbv1.Tuple, *pbv1.Tuple]) error {
		for _, v := range w {
			t := (&pbv1.Tuple{StoreId: store, TupleKey: v, CreatedAt: timestamppb.Now()}).SetKey()
			got, _, err := tx.Get(ctx, t)
			if err != nil {
				return err
			}
			if len(got) != 0 {
				return storage.InvalidWriteInputError(v, openfgav1.TupleOperation_TUPLE_OPERATION_WRITE)
			}
			// don't know why the tests want that...
			if t.TupleKey.Condition != nil && v.Condition.Context == nil {
				t.TupleKey.Condition.Context = &structpb.Struct{}
			}
			if _, err := tx.Raw().Set(ctx, pbv1.NewWriteChange(store, t)); err != nil {
				return err
			}
			if _, err := tx.Set(ctx, t); err != nil {
				return err
			}
		}
		for _, v := range d {
			t := &pbv1.Tuple{Key: pbv1.TupleKey(store, tuple.TupleKeyWithoutConditionToTupleKey(v))}
			got, _, err := tx.Get(ctx, t)
			if err != nil {
				return err
			}
			if len(got) == 0 {
				return storage.InvalidWriteInputError(v, openfgav1.TupleOperation_TUPLE_OPERATION_DELETE)
			}
			if err := tx.Delete(ctx, t); err != nil {
				return err
			}
			got[0].TupleKey.Condition = nil
			if _, err := tx.Raw().Set(ctx, pbv1.NewDeleteChange(store, got[0])); err != nil {
				return err
			}
		}
		return nil
	})
}

func (p *pdb) MaxTuplesPerWrite() int {
	// TODO(adphi): implement
	return 100
}

func (p *pdb) ReadAuthorizationModel(ctx context.Context, store string, id string) (*openfgav1.AuthorizationModel, error) {
	return typed.WithTypedTx2(ctx, p.db, func(ctx context.Context, tx typed.Tx[pbv1.Model, *pbv1.Model]) (*openfgav1.AuthorizationModel, error) {
		ms, _, err := tx.Get(ctx, &pbv1.Model{Key: fmt.Sprintf("%s/%s", store, id)})
		if err != nil {
			return nil, err
		}
		if len(ms) == 0 {
			return nil, storage.ErrNotFound
		}
		return ms[0].Model, nil
	})
}

func (p *pdb) ReadAuthorizationModels(ctx context.Context, store string, opts storage.ReadAuthorizationModelsOptions) (out []*openfgav1.AuthorizationModel, continuation []byte, err error) {
	tk, err := parseToken(opts.Pagination.From, modelType)
	if err != nil {
		return nil, nil, err
	}
	ms, info, err := typed.NewStore[pbv1.Model](p.db).Get(ctx, &pbv1.Model{}, protodb.WithFilter(protodb.Where("key").StringHasPrefix(pbv1.ModelPrefix(store))), protodb.WithPaging(&protodb.Paging{Token: tk.Continuation, Offset: tk.Offset, Limit: uint64(opts.Pagination.PageSize)}), protodb.WithReverse())
	if err != nil {
		return nil, nil, err
	}
	for _, v := range ms {
		out = append(out, v.Model)
	}
	if !info.HasNext {
		return out, nil, nil
	}
	return out, newToken(info.Token, tk.Offset+uint64(len(ms)), modelType).Bytes(), nil

}

func (p *pdb) FindLatestAuthorizationModel(ctx context.Context, store string) (*openfgav1.AuthorizationModel, error) {
	ms, _, err := typed.NewStore[pbv1.Model](p.db).Get(ctx, &pbv1.Model{}, protodb.WithFilter(protodb.Where("key").StringHasPrefix(pbv1.ModelPrefix(store))), protodb.WithReverse())
	if err != nil {
		return nil, err
	}
	if len(ms) == 0 {
		return nil, storage.ErrNotFound
	}
	return ms[0].Model, nil
}

func (p *pdb) MaxTypesPerAuthorizationModel() int {
	// TODO implement me
	return 100
}

func (p *pdb) WriteAuthorizationModel(ctx context.Context, store string, model *openfgav1.AuthorizationModel) error {
	_, err := p.db.Set(ctx, (&pbv1.Model{StoreId: store, Model: model}).SetKey())
	return err
}

func (p *pdb) CreateStore(ctx context.Context, store *openfgav1.Store) (*openfgav1.Store, error) {
	return typed.WithTypedTx2[openfgav1.Store, *openfgav1.Store](ctx, p.db, func(ctx context.Context, tx typed.Tx[openfgav1.Store, *openfgav1.Store]) (*openfgav1.Store, error) {
		ss, _, err := tx.Get(ctx, &openfgav1.Store{Id: store.Id})
		if err != nil {
			return nil, err
		}
		if len(ss) != 0 {
			return nil, storage.ErrCollision
		}
		return typed.NewStore[openfgav1.Store](p.db).Set(ctx, store)
	})
}

func (p *pdb) DeleteStore(ctx context.Context, id string) error {
	return protodb.WithTx(ctx, p.db, func(ctx context.Context, tx protodb.Tx) error {
		ss, _, err := tx.Get(ctx, &openfgav1.Store{Id: id})
		if err != nil {
			return err
		}
		if len(ss) == 0 {
			return storage.ErrNotFound
		}
		ts, _, err := tx.Get(ctx, &pbv1.Tuple{}, protodb.WithFilter(protodb.Where("key").StringHasPrefix(pbv1.TuplePrefix(id))))
		if err != nil {
			return err
		}
		for _, v := range ts {
			if err := tx.Delete(ctx, v); err != nil {
				return err
			}
		}
		cs, _, err := tx.Get(ctx, &pbv1.Change{}, protodb.WithFilter(protodb.Where("key").StringHasPrefix(pbv1.TuplePrefix(id))))
		if err != nil {
			return err
		}
		for _, v := range cs {
			if err := tx.Delete(ctx, v); err != nil {
				return err
			}
		}
		as, _, err := tx.Get(ctx, &pbv1.Assertions{}, protodb.WithFilter(protodb.Where("key").StringHasPrefix(pbv1.TuplePrefix(id))))
		if err != nil {
			return err
		}
		for _, v := range as {
			if err := tx.Delete(ctx, v); err != nil {
				return err
			}
		}
		return tx.Delete(ctx, &openfgav1.Store{Id: id})
	})
}

func (p *pdb) GetStore(ctx context.Context, id string) (*openfgav1.Store, error) {
	ss, _, err := typed.NewStore[openfgav1.Store](p.db).Get(ctx, &openfgav1.Store{Id: id})
	if err != nil {
		return nil, err
	}
	if len(ss) == 0 {
		return nil, storage.ErrNotFound
	}
	return ss[0], nil
}

func (p *pdb) ListStores(ctx context.Context, opts storage.ListStoresOptions) ([]*openfgav1.Store, []byte, error) {
	tk, err := parseToken(opts.Pagination.From, storeType)
	if err != nil {
		return nil, nil, err
	}
	ss, info, err := typed.NewStore[openfgav1.Store](p.db).Get(ctx, &openfgav1.Store{}, protodb.WithPaging(&protodb.Paging{
		Offset: tk.Offset,
		Limit:  uint64(opts.Pagination.PageSize),
		Token:  tk.Continuation,
	}))
	if err != nil {
		return nil, nil, err
	}
	if !info.HasNext {
		return ss, nil, nil
	}
	return ss, newToken(info.Token, tk.Offset+uint64(len(ss)), storeType).Bytes(), nil
}

func (p *pdb) WriteAssertions(ctx context.Context, store, modelID string, assertions []*openfgav1.Assertion) error {
	_, err := p.db.Set(ctx, (&pbv1.Assertions{StoreId: store, ModelId: modelID, Assertions: assertions}).SetKey())
	return err
}

func (p *pdb) ReadAssertions(ctx context.Context, store, modelID string) ([]*openfgav1.Assertion, error) {
	as, _, err := typed.NewStore[pbv1.Assertions](p.db).Get(ctx, &pbv1.Assertions{Key: pbv1.AssertionsKey(store, modelID)})
	if err != nil {
		return nil, err
	}
	if len(as) == 0 {
		return nil, nil
	}
	return as[0].Assertions, nil
}

func (p *pdb) ReadChanges(ctx context.Context, store, objectType string, opts storage.ReadChangesOptions, horizonOffset time.Duration) ([]*openfgav1.TupleChange, []byte, error) {
	tk, err := parseToken(opts.Pagination.From, objectType)
	if err != nil {
		return nil, nil, err
	}
	f := protodb.Where("key").StringHasPrefix(store + "/")
	// we do not set the horizontal offset in the filter because it may change
	if objectType != "" {
		f = f.And("change.tuple_key.object").StringHasPrefix(objectType + ":")
	}
	cs, info, err := typed.NewStore[pbv1.Change](p.db).Get(ctx, &pbv1.Change{}, protodb.WithFilter(f), protodb.WithPaging(&protodb.Paging{
		Offset: tk.Offset,
		Limit:  uint64(opts.Pagination.PageSize),
		Token:  tk.Continuation,
	}))
	if err != nil {
		return nil, nil, err
	}
	var out []*openfgav1.TupleChange
	for _, v := range cs {
		if v.Change.Timestamp.AsTime().After(time.Now().Add(-horizonOffset)) {
			continue
		}
		out = append(out, v.Change)
	}
	if len(out) == 0 {
		return nil, nil, storage.ErrNotFound
	}
	return out, newToken(info.Token, tk.Offset+uint64(len(out)), objectType).Bytes(), nil
}

func (p *pdb) IsReady(ctx context.Context) (storage.ReadinessStatus, error) {
	return storage.ReadinessStatus{IsReady: true}, nil
}

func (p *pdb) Close() {
	if !p.external {
		p.db.Close()
	}
}

const (
	storeType = "store"
	modelType = "model"
)

func newToken(continuation string, offset uint64, objectType string) *token {
	return &token{Offset: offset, Continuation: continuation, ObjectType: objectType}
}

type token struct {
	Offset       uint64
	Continuation string
	ObjectType   string
}

func (t *token) Bytes() []byte {
	return []byte(fmt.Sprintf("%d:%s:%s", t.Offset, t.Continuation, t.ObjectType))
}

func (t *token) validate(objectType string) error {
	if t.ObjectType == "" && t.Offset == 0 {
		return nil
	}
	if t.ObjectType != objectType {
		return storage.ErrMismatchObjectType
	}
	return nil
}

func parseToken(t string, objectType string) (*token, error) {
	if len(t) == 0 {
		return &token{}, (&token{}).validate(objectType)
	}
	parts := strings.SplitN(t, ":", 3)
	if len(parts) != 3 {
		return nil, storage.ErrInvalidContinuationToken
	}
	offset, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return nil, storage.ErrInvalidContinuationToken
	}
	tk := &token{Offset: offset, Continuation: parts[1], ObjectType: parts[2]}
	if err := tk.validate(objectType); err != nil {
		return nil, err
	}
	return tk, nil
}

func makeFilter(storeID, user, relation, object string) filters.FieldFilterer {
	var f filters.Builder
	if object != "" {
		if o, id := tuple.SplitObject(object); id != "" {
			prefix := storeID + "/" + o + ":" + id + "#"
			if relation != "" {
				prefix += relation + "@"
			}
			f = protodb.Where("key").StringHasPrefix(prefix)
		} else {
			f = protodb.Where("key").StringHasPrefix(fmt.Sprintf("%s/%s:", storeID, o))
		}
	} else {
		f = protodb.Where("key").StringHasPrefix(fmt.Sprintf("%s/", storeID))
	}
	if user != "" {
		suffix := "@" + user
		if relation != "" {
			suffix = "#" + relation + suffix
		}
		f.And("key").StringHasSuffix(suffix)
	}
	return f
}

func toTuple(ts []*pbv1.Tuple) []*openfgav1.Tuple {
	var out []*openfgav1.Tuple
	for _, v := range ts {
		out = append(out, &openfgav1.Tuple{Key: v.TupleKey, Timestamp: v.CreatedAt})
	}
	return out
}
