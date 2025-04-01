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

	openfgav1 "github.com/openfga/api/proto/openfga/v1"

	"go.linka.cloud/go-openfga/storage"
	"go.linka.cloud/go-openfga/x"
	"go.linka.cloud/go-openfga/x/service"
)

type none struct{}

var _ x.OpenFGA = (*clientWrapper)(nil)

func wrap(c service.Client) x.OpenFGA {
	return &clientWrapper{c: c, wrapper: &wrapper{c: c}}
}

type clientWrapper struct {
	*wrapper
	c service.Client
}

func (c *clientWrapper) WithTx(_ any) x.Tx {
	panic("unsupported")
}

func (c *clientWrapper) Tx(ctx context.Context, opts ...storage.TxOption) (x.Tx, error) {
	tx, err := c.c.Tx(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &txWrapper{c: tx, wrapper: &wrapper{c: tx}}, nil
}

func (c *clientWrapper) Close() {}

type txWrapper struct {
	*wrapper
	c service.Tx
}

func (w *txWrapper) Commit(ctx context.Context) error {
	return w.c.Commit(ctx)
}

func (w *txWrapper) Close() error {
	return w.c.Close()
}

type wrapper struct {
	openfgav1.UnsafeOpenFGAServiceServer
	c openfgav1.OpenFGAServiceClient
}

func (c *wrapper) Read(ctx context.Context, req *openfgav1.ReadRequest) (*openfgav1.ReadResponse, error) {
	return c.c.Read(ctx, req)
}

func (c *wrapper) Write(ctx context.Context, req *openfgav1.WriteRequest) (*openfgav1.WriteResponse, error) {
	return c.c.Write(ctx, req)
}

func (c *wrapper) Check(ctx context.Context, req *openfgav1.CheckRequest) (*openfgav1.CheckResponse, error) {
	return c.c.Check(ctx, req)
}

func (c *wrapper) BatchCheck(ctx context.Context, req *openfgav1.BatchCheckRequest) (*openfgav1.BatchCheckResponse, error) {
	return c.c.BatchCheck(ctx, req)
}

func (c *wrapper) Expand(ctx context.Context, req *openfgav1.ExpandRequest) (*openfgav1.ExpandResponse, error) {
	return c.c.Expand(ctx, req)
}

func (c *wrapper) ReadAuthorizationModels(ctx context.Context, req *openfgav1.ReadAuthorizationModelsRequest) (*openfgav1.ReadAuthorizationModelsResponse, error) {
	return c.c.ReadAuthorizationModels(ctx, req)
}

func (c *wrapper) ReadAuthorizationModel(ctx context.Context, req *openfgav1.ReadAuthorizationModelRequest) (*openfgav1.ReadAuthorizationModelResponse, error) {
	return c.c.ReadAuthorizationModel(ctx, req)
}

func (c *wrapper) WriteAuthorizationModel(ctx context.Context, req *openfgav1.WriteAuthorizationModelRequest) (*openfgav1.WriteAuthorizationModelResponse, error) {
	return c.c.WriteAuthorizationModel(ctx, req)
}

func (c *wrapper) WriteAssertions(ctx context.Context, req *openfgav1.WriteAssertionsRequest) (*openfgav1.WriteAssertionsResponse, error) {
	return c.c.WriteAssertions(ctx, req)
}

func (c *wrapper) ReadAssertions(ctx context.Context, req *openfgav1.ReadAssertionsRequest) (*openfgav1.ReadAssertionsResponse, error) {
	return c.c.ReadAssertions(ctx, req)
}

func (c *wrapper) ReadChanges(ctx context.Context, req *openfgav1.ReadChangesRequest) (*openfgav1.ReadChangesResponse, error) {
	return c.c.ReadChanges(ctx, req)
}

func (c *wrapper) CreateStore(ctx context.Context, req *openfgav1.CreateStoreRequest) (*openfgav1.CreateStoreResponse, error) {
	return c.c.CreateStore(ctx, req)
}

func (c *wrapper) UpdateStore(ctx context.Context, req *openfgav1.UpdateStoreRequest) (*openfgav1.UpdateStoreResponse, error) {
	return c.c.UpdateStore(ctx, req)
}

func (c *wrapper) DeleteStore(ctx context.Context, req *openfgav1.DeleteStoreRequest) (*openfgav1.DeleteStoreResponse, error) {
	return c.c.DeleteStore(ctx, req)
}

func (c *wrapper) GetStore(ctx context.Context, req *openfgav1.GetStoreRequest) (*openfgav1.GetStoreResponse, error) {
	return c.c.GetStore(ctx, req)
}

func (c *wrapper) ListStores(ctx context.Context, req *openfgav1.ListStoresRequest) (*openfgav1.ListStoresResponse, error) {
	return c.c.ListStores(ctx, req)
}

func (c *wrapper) StreamedListObjects(req *openfgav1.StreamedListObjectsRequest, ss openfgav1.OpenFGAService_StreamedListObjectsServer) error {
	panic("unsupported")
}

func (c *wrapper) ListObjects(ctx context.Context, req *openfgav1.ListObjectsRequest) (*openfgav1.ListObjectsResponse, error) {
	return c.c.ListObjects(ctx, req)
}

func (c *wrapper) ListUsers(ctx context.Context, req *openfgav1.ListUsersRequest) (*openfgav1.ListUsersResponse, error) {
	return c.c.ListUsers(ctx, req)
}
