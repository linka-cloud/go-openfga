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

package x

import (
	"context"

	"github.com/fullstorydev/grpchan/inprocgrpc"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/openfga/openfga/pkg/middleware/requestid"
	"github.com/openfga/openfga/pkg/middleware/storeid"
	"github.com/openfga/openfga/pkg/middleware/validator"
	"github.com/openfga/openfga/pkg/server"
	"google.golang.org/grpc"

	"go.linka.cloud/go-openfga"
	"go.linka.cloud/go-openfga/storage"
	pbv1 "go.linka.cloud/go-openfga/x/pb/v1"
)

type FGA interface {
	Tx(ctx context.Context, opts ...storage.TxOption) (Tx, error)
}

type Tx interface {
	openfga.Client
	Commit(ctx context.Context) error
	Close()
}

func New[T any](s storage.Datastore[T], opts ...server.OpenFGAServiceV1Option) (FGA, error) {
	svc, err := Wrap(s, opts...)
	if err != nil {
		return nil, err
	}
	ch := &inprocgrpc.Channel{}
	ch.WithServerUnaryInterceptor(
		grpcmiddleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(), // needed for logging
			requestid.NewUnaryInterceptor(),       // add request_id to ctxtags
			storeid.NewUnaryInterceptor(),         // if available, add store_id to ctxtags
			// logging.NewLoggingInterceptor(s.Logger), // needed to log invalid requests
			validator.UnaryServerInterceptor(),
		),
	)
	ch.WithServerStreamInterceptor(
		grpcmiddleware.ChainStreamServer(
			[]grpc.StreamServerInterceptor{
				requestid.NewStreamingInterceptor(),
				validator.StreamServerInterceptor(),
				grpc_ctxtags.StreamServerInterceptor(),
			}...,
		),
	)
	pbv1.RegisterOpenFGAXServiceServer(ch, svc)
	return &fga[T]{c: newClient(ch)}, nil
}

type fga[T any] struct {
	c *client
}

func (f *fga[T]) Tx(ctx context.Context, opts ...storage.TxOption) (Tx, error) {
	t, err := f.c.Tx(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &tx{tx: t, Client: openfga.FomClient(t)}, nil
}

type tx struct {
	tx *txc
	openfga.Client
}

func (t *tx) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *tx) Close() {
	t.tx.Close()
}
