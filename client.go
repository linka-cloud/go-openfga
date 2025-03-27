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
	"google.golang.org/grpc"

	"go.linka.cloud/go-openfga/x"
)

func NewClient(c grpc.ClientConnInterface) Client {
	return &client{c: x.NewClient(c)}
}

type client struct {
	c x.Client
}

func (c *client) CreateStore(ctx context.Context, name string) (Store, error) {
	res, err := c.c.CreateStore(ctx, &openfgav1.CreateStoreRequest{Name: name})
	if err != nil {
		return nil, err
	}
	return &store{c: c, id: res.Id, name: res.Name, createdAt: res.CreatedAt.AsTime(), updatedAt: res.UpdatedAt.AsTime()}, nil
}

func (c *client) GetStore(ctx context.Context, id string) (Store, error) {
	res, err := c.c.GetStore(ctx, &openfgav1.GetStoreRequest{StoreId: id})
	if err != nil {
		return nil, err
	}
	return &store{c: c, id: res.Id, name: res.Name, createdAt: res.CreatedAt.AsTime(), updatedAt: res.UpdatedAt.AsTime()}, nil
}

func (c *client) ListStores(ctx context.Context) ([]Store, error) {
	res, err := c.c.ListStores(ctx, &openfgav1.ListStoresRequest{})
	if err != nil {
		return nil, err
	}
	var stores []Store
	for _, s := range res.Stores {
		stores = append(stores, &store{c: c, id: s.Id, name: s.Name, createdAt: s.CreatedAt.AsTime(), updatedAt: s.UpdatedAt.AsTime()})
	}
	return stores, nil
}

func (c *client) DeleteStore(ctx context.Context, id string) error {
	_, err := c.c.DeleteStore(ctx, &openfgav1.DeleteStoreRequest{StoreId: id})
	return err
}
