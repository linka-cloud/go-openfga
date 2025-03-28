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

package x

import (
	"context"

	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	"github.com/openfga/openfga/pkg/server"

	"go.linka.cloud/go-openfga/storage"
)

type OpenFGA[T any] interface {
	openfgav1.OpenFGAServiceServer
	Tx(ctx context.Context, opts ...storage.TxOption) (Tx, error)
	WithTx(tx T) Tx
	Close()
}

type Tx interface {
	openfgav1.OpenFGAServiceServer
	Commit(ctx context.Context) error
	Close() error
}

func New[T any](s storage.Datastore[T], opts ...server.OpenFGAServiceV1Option) (OpenFGA[T], error) {
	opts = append(opts, server.WithDatastore(s))
	srv, err := server.NewServerWithOpts(opts...)
	if err != nil {
		return nil, err
	}
	return &svc[T]{Server: srv, s: s}, nil
}

type svc[T any] struct {
	*server.Server
	s storage.Datastore[T]
}

func (s *svc[T]) Tx(ctx context.Context, opts ...storage.TxOption) (Tx, error) {
	t, err := s.s.Tx(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &tx[T]{c: s, tx: t}, nil
}

func (s *svc[T]) WithTx(t T) Tx {
	return &tx[T]{c: s, tx: s.s.WithTx(t)}
}

type tx[T any] struct {
	openfgav1.UnsafeOpenFGAServiceServer
	c  openfgav1.OpenFGAServiceServer
	tx storage.Tx[T]
}

func (t *tx[T]) Read(ctx context.Context, req *openfgav1.ReadRequest) (*openfgav1.ReadResponse, error) {
	return t.c.Read(storage.Context(ctx, t.tx), req)
}

func (t *tx[T]) Write(ctx context.Context, req *openfgav1.WriteRequest) (*openfgav1.WriteResponse, error) {
	return t.c.Write(storage.Context(ctx, t.tx), req)
}

func (t *tx[T]) Check(ctx context.Context, req *openfgav1.CheckRequest) (*openfgav1.CheckResponse, error) {
	return t.c.Check(storage.Context(ctx, t.tx), req)
}

func (t *tx[T]) BatchCheck(ctx context.Context, req *openfgav1.BatchCheckRequest) (*openfgav1.BatchCheckResponse, error) {
	return t.c.BatchCheck(storage.Context(ctx, t.tx), req)
}

func (t *tx[T]) Expand(ctx context.Context, req *openfgav1.ExpandRequest) (*openfgav1.ExpandResponse, error) {
	return t.c.Expand(storage.Context(ctx, t.tx), req)
}

func (t *tx[T]) ReadAuthorizationModels(ctx context.Context, req *openfgav1.ReadAuthorizationModelsRequest) (*openfgav1.ReadAuthorizationModelsResponse, error) {
	return t.c.ReadAuthorizationModels(storage.Context(ctx, t.tx), req)
}

func (t *tx[T]) ReadAuthorizationModel(ctx context.Context, req *openfgav1.ReadAuthorizationModelRequest) (*openfgav1.ReadAuthorizationModelResponse, error) {
	return t.c.ReadAuthorizationModel(storage.Context(ctx, t.tx), req)
}

func (t *tx[T]) WriteAuthorizationModel(ctx context.Context, req *openfgav1.WriteAuthorizationModelRequest) (*openfgav1.WriteAuthorizationModelResponse, error) {
	return t.c.WriteAuthorizationModel(storage.Context(ctx, t.tx), req)
}

func (t *tx[T]) WriteAssertions(ctx context.Context, req *openfgav1.WriteAssertionsRequest) (*openfgav1.WriteAssertionsResponse, error) {
	return t.c.WriteAssertions(storage.Context(ctx, t.tx), req)
}

func (t *tx[T]) ReadAssertions(ctx context.Context, req *openfgav1.ReadAssertionsRequest) (*openfgav1.ReadAssertionsResponse, error) {
	return t.c.ReadAssertions(storage.Context(ctx, t.tx), req)
}

func (t *tx[T]) ReadChanges(ctx context.Context, req *openfgav1.ReadChangesRequest) (*openfgav1.ReadChangesResponse, error) {
	return t.c.ReadChanges(storage.Context(ctx, t.tx), req)
}

func (t *tx[T]) CreateStore(ctx context.Context, req *openfgav1.CreateStoreRequest) (*openfgav1.CreateStoreResponse, error) {
	return t.c.CreateStore(storage.Context(ctx, t.tx), req)
}

func (t *tx[T]) UpdateStore(ctx context.Context, req *openfgav1.UpdateStoreRequest) (*openfgav1.UpdateStoreResponse, error) {
	return t.c.UpdateStore(storage.Context(ctx, t.tx), req)
}

func (t *tx[T]) DeleteStore(ctx context.Context, req *openfgav1.DeleteStoreRequest) (*openfgav1.DeleteStoreResponse, error) {
	return t.c.DeleteStore(storage.Context(ctx, t.tx), req)
}

func (t *tx[T]) GetStore(ctx context.Context, req *openfgav1.GetStoreRequest) (*openfgav1.GetStoreResponse, error) {
	return t.c.GetStore(storage.Context(ctx, t.tx), req)
}

func (t *tx[T]) ListStores(ctx context.Context, req *openfgav1.ListStoresRequest) (*openfgav1.ListStoresResponse, error) {
	return t.c.ListStores(storage.Context(ctx, t.tx), req)
}

func (t *tx[T]) StreamedListObjects(req *openfgav1.StreamedListObjectsRequest, server openfgav1.OpenFGAService_StreamedListObjectsServer) error {
	panic("unsupported")
}

func (t *tx[T]) ListObjects(ctx context.Context, req *openfgav1.ListObjectsRequest) (*openfgav1.ListObjectsResponse, error) {
	return t.c.ListObjects(storage.Context(ctx, t.tx), req)
}

func (t *tx[T]) ListUsers(ctx context.Context, req *openfgav1.ListUsersRequest) (*openfgav1.ListUsersResponse, error) {
	return t.c.ListUsers(storage.Context(ctx, t.tx), req)
}

func (t *tx[T]) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *tx[T]) Close() error {
	t.tx.Close()
	return nil
}
