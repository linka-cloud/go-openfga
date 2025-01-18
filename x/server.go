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
	"errors"
	"io"

	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	"github.com/openfga/openfga/pkg/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.linka.cloud/go-openfga/storage"
	pbv1 "go.linka.cloud/go-openfga/x/pb/v1"
)

func Wrap[T any](s storage.Datastore[T], opts ...server.OpenFGAServiceV1Option) (pbv1.OpenFGAXServiceServer, error) {
	opts = append(opts, server.WithDatastore(s))
	srv, err := server.NewServerWithOpts(opts...)
	if err != nil {
		return nil, err
	}
	return &svc[T]{c: srv, s: s}, nil
}

type svc[T any] struct {
	pbv1.UnsafeOpenFGAXServiceServer
	c openfgav1.OpenFGAServiceServer
	s storage.Datastore[T]
}

func (s *svc[T]) Tx(g grpc.BidiStreamingServer[pbv1.TxRequest, pbv1.TxResponse]) error {
	ctx := g.Context()
	var tx storage.Tx[T]
	defer func() {
		if tx != nil {
			tx.Close()
		}
	}()
	for {
		req, err := g.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		if tx == nil {
			var o []storage.TxOption
			if req.ReadOnly {
				o = append(o, storage.WithReadOnly())
			}
			tx, err = s.s.Tx(g.Context(), o...)
			if err != nil {
				return err
			}
			ctx = storage.Context(ctx, tx)
		}
		var res *pbv1.TxResponse
		switch r := req.Request.(type) {
		case *pbv1.TxRequest_Read:
			v, err := s.c.Read(ctx, r.Read)
			if err != nil {
				return err
			}
			res = &pbv1.TxResponse{Response: &pbv1.TxResponse_Read{
				Read: v,
			}}
		case *pbv1.TxRequest_Write:
			v, err := s.c.Write(ctx, r.Write)
			if err != nil {
				return err
			}
			res = &pbv1.TxResponse{Response: &pbv1.TxResponse_Write{
				Write: v,
			}}
		case *pbv1.TxRequest_Check:
			v, err := s.c.Check(ctx, r.Check)
			if err != nil {
				return err
			}
			res = &pbv1.TxResponse{Response: &pbv1.TxResponse_Check{
				Check: v,
			}}
		case *pbv1.TxRequest_BatchCheck:
			v, err := s.c.BatchCheck(ctx, r.BatchCheck)
			if err != nil {
				return err
			}
			res = &pbv1.TxResponse{Response: &pbv1.TxResponse_BatchCheck{
				BatchCheck: v,
			}}
		case *pbv1.TxRequest_Expand:
			v, err := s.c.Expand(ctx, r.Expand)
			if err != nil {
				return err
			}
			res = &pbv1.TxResponse{Response: &pbv1.TxResponse_Expand{
				Expand: v,
			}}
		case *pbv1.TxRequest_ReadAuthorizationModels:
			v, err := s.c.ReadAuthorizationModels(ctx, r.ReadAuthorizationModels)
			if err != nil {
				return err
			}
			res = &pbv1.TxResponse{Response: &pbv1.TxResponse_ReadAuthorizationModels{
				ReadAuthorizationModels: v,
			}}
		case *pbv1.TxRequest_ReadAuthorizationModel:
			v, err := s.c.ReadAuthorizationModel(ctx, r.ReadAuthorizationModel)
			if err != nil {
				return err
			}
			res = &pbv1.TxResponse{Response: &pbv1.TxResponse_ReadAuthorizationModel{
				ReadAuthorizationModel: v,
			}}
		case *pbv1.TxRequest_WriteAuthorizationModel:
			v, err := s.c.WriteAuthorizationModel(ctx, r.WriteAuthorizationModel)
			if err != nil {
				return err
			}
			res = &pbv1.TxResponse{Response: &pbv1.TxResponse_WriteAuthorizationModel{
				WriteAuthorizationModel: v,
			}}
		case *pbv1.TxRequest_WriteAssertions:
			v, err := s.c.WriteAssertions(ctx, r.WriteAssertions)
			if err != nil {
				return err
			}
			res = &pbv1.TxResponse{Response: &pbv1.TxResponse_WriteAssertions{
				WriteAssertions: v,
			}}
		case *pbv1.TxRequest_ReadAssertions:
			v, err := s.c.ReadAssertions(ctx, r.ReadAssertions)
			if err != nil {
				return err
			}
			res = &pbv1.TxResponse{Response: &pbv1.TxResponse_ReadAssertions{
				ReadAssertions: v,
			}}
		case *pbv1.TxRequest_ReadChanges:
			v, err := s.c.ReadChanges(ctx, r.ReadChanges)
			if err != nil {
				return err
			}
			res = &pbv1.TxResponse{Response: &pbv1.TxResponse_ReadChanges{
				ReadChanges: v,
			}}
		case *pbv1.TxRequest_CreateStore:
			v, err := s.c.CreateStore(ctx, r.CreateStore)
			if err != nil {
				return err
			}
			res = &pbv1.TxResponse{Response: &pbv1.TxResponse_CreateStore{
				CreateStore: v,
			}}
		case *pbv1.TxRequest_UpdateStore:
			v, err := s.c.UpdateStore(ctx, r.UpdateStore)
			if err != nil {
				return err
			}
			res = &pbv1.TxResponse{Response: &pbv1.TxResponse_UpdateStore{
				UpdateStore: v,
			}}
		case *pbv1.TxRequest_DeleteStore:
			v, err := s.c.DeleteStore(ctx, r.DeleteStore)
			if err != nil {
				return err
			}
			res = &pbv1.TxResponse{Response: &pbv1.TxResponse_DeleteStore{
				DeleteStore: v,
			}}
		case *pbv1.TxRequest_GetStore:
			v, err := s.c.GetStore(ctx, r.GetStore)
			if err != nil {
				return err
			}
			res = &pbv1.TxResponse{Response: &pbv1.TxResponse_GetStore{
				GetStore: v,
			}}
		case *pbv1.TxRequest_ListStores:
			v, err := s.c.ListStores(ctx, r.ListStores)
			if err != nil {
				return err
			}
			res = &pbv1.TxResponse{Response: &pbv1.TxResponse_ListStores{
				ListStores: v,
			}}
		case *pbv1.TxRequest_ListObjects:
			v, err := s.c.ListObjects(ctx, r.ListObjects)
			if err != nil {
				return err
			}
			res = &pbv1.TxResponse{Response: &pbv1.TxResponse_ListObjects{
				ListObjects: v,
			}}
		case *pbv1.TxRequest_ListUsers:
			v, err := s.c.ListUsers(ctx, r.ListUsers)
			if err != nil {
				return err
			}
			res = &pbv1.TxResponse{Response: &pbv1.TxResponse_ListUsers{
				ListUsers: v,
			}}
		case *pbv1.TxRequest_Commit:
			if err := tx.Commit(ctx); err != nil {
				return err
			}
			return nil
		default:
			return status.Error(codes.InvalidArgument, "unexpected request type")
		}
		if err := g.Send(res); err != nil {
			return err
		}
	}
}
