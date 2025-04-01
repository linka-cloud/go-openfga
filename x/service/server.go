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

package service

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

type Server interface {
	openfgav1.OpenFGAServiceServer
	pbv1.OpenFGAXServiceServer
	Close()
}

func Wrap(s storage.Datastore, opts ...server.OpenFGAServiceV1Option) (Server, error) {
	opts = append(opts, server.WithDatastore(s))
	srv, err := server.NewServerWithOpts(opts...)
	if err != nil {
		return nil, err
	}
	return &svc{OpenFGAServiceServer: srv, s: s}, nil
}

type svc struct {
	pbv1.UnsafeOpenFGAXServiceServer
	openfgav1.OpenFGAServiceServer
	s storage.Datastore
}

func (s *svc) Tx(g grpc.BidiStreamingServer[pbv1.TxRequest, pbv1.TxResponse]) error {
	ctx := g.Context()
	var tx storage.Tx
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
			v, err := s.OpenFGAServiceServer.Read(ctx, r.Read)
			if err != nil {
				return err
			}
			res = &pbv1.TxResponse{Response: &pbv1.TxResponse_Read{
				Read: v,
			}}
		case *pbv1.TxRequest_Write:
			v, err := s.OpenFGAServiceServer.Write(ctx, r.Write)
			if err != nil {
				return err
			}
			res = &pbv1.TxResponse{Response: &pbv1.TxResponse_Write{
				Write: v,
			}}
		case *pbv1.TxRequest_Check:
			v, err := s.OpenFGAServiceServer.Check(ctx, r.Check)
			if err != nil {
				return err
			}
			res = &pbv1.TxResponse{Response: &pbv1.TxResponse_Check{
				Check: v,
			}}
		case *pbv1.TxRequest_BatchCheck:
			v, err := s.OpenFGAServiceServer.BatchCheck(ctx, r.BatchCheck)
			if err != nil {
				return err
			}
			res = &pbv1.TxResponse{Response: &pbv1.TxResponse_BatchCheck{
				BatchCheck: v,
			}}
		case *pbv1.TxRequest_Expand:
			v, err := s.OpenFGAServiceServer.Expand(ctx, r.Expand)
			if err != nil {
				return err
			}
			res = &pbv1.TxResponse{Response: &pbv1.TxResponse_Expand{
				Expand: v,
			}}
		case *pbv1.TxRequest_ReadAuthorizationModels:
			v, err := s.OpenFGAServiceServer.ReadAuthorizationModels(ctx, r.ReadAuthorizationModels)
			if err != nil {
				return err
			}
			res = &pbv1.TxResponse{Response: &pbv1.TxResponse_ReadAuthorizationModels{
				ReadAuthorizationModels: v,
			}}
		case *pbv1.TxRequest_ReadAuthorizationModel:
			v, err := s.OpenFGAServiceServer.ReadAuthorizationModel(ctx, r.ReadAuthorizationModel)
			if err != nil {
				return err
			}
			res = &pbv1.TxResponse{Response: &pbv1.TxResponse_ReadAuthorizationModel{
				ReadAuthorizationModel: v,
			}}
		case *pbv1.TxRequest_WriteAuthorizationModel:
			v, err := s.OpenFGAServiceServer.WriteAuthorizationModel(ctx, r.WriteAuthorizationModel)
			if err != nil {
				return err
			}
			res = &pbv1.TxResponse{Response: &pbv1.TxResponse_WriteAuthorizationModel{
				WriteAuthorizationModel: v,
			}}
		case *pbv1.TxRequest_WriteAssertions:
			v, err := s.OpenFGAServiceServer.WriteAssertions(ctx, r.WriteAssertions)
			if err != nil {
				return err
			}
			res = &pbv1.TxResponse{Response: &pbv1.TxResponse_WriteAssertions{
				WriteAssertions: v,
			}}
		case *pbv1.TxRequest_ReadAssertions:
			v, err := s.OpenFGAServiceServer.ReadAssertions(ctx, r.ReadAssertions)
			if err != nil {
				return err
			}
			res = &pbv1.TxResponse{Response: &pbv1.TxResponse_ReadAssertions{
				ReadAssertions: v,
			}}
		case *pbv1.TxRequest_ReadChanges:
			v, err := s.OpenFGAServiceServer.ReadChanges(ctx, r.ReadChanges)
			if err != nil {
				return err
			}
			res = &pbv1.TxResponse{Response: &pbv1.TxResponse_ReadChanges{
				ReadChanges: v,
			}}
		case *pbv1.TxRequest_CreateStore:
			v, err := s.OpenFGAServiceServer.CreateStore(ctx, r.CreateStore)
			if err != nil {
				return err
			}
			res = &pbv1.TxResponse{Response: &pbv1.TxResponse_CreateStore{
				CreateStore: v,
			}}
		case *pbv1.TxRequest_UpdateStore:
			v, err := s.OpenFGAServiceServer.UpdateStore(ctx, r.UpdateStore)
			if err != nil {
				return err
			}
			res = &pbv1.TxResponse{Response: &pbv1.TxResponse_UpdateStore{
				UpdateStore: v,
			}}
		case *pbv1.TxRequest_DeleteStore:
			v, err := s.OpenFGAServiceServer.DeleteStore(ctx, r.DeleteStore)
			if err != nil {
				return err
			}
			res = &pbv1.TxResponse{Response: &pbv1.TxResponse_DeleteStore{
				DeleteStore: v,
			}}
		case *pbv1.TxRequest_GetStore:
			v, err := s.OpenFGAServiceServer.GetStore(ctx, r.GetStore)
			if err != nil {
				return err
			}
			res = &pbv1.TxResponse{Response: &pbv1.TxResponse_GetStore{
				GetStore: v,
			}}
		case *pbv1.TxRequest_ListStores:
			v, err := s.OpenFGAServiceServer.ListStores(ctx, r.ListStores)
			if err != nil {
				return err
			}
			res = &pbv1.TxResponse{Response: &pbv1.TxResponse_ListStores{
				ListStores: v,
			}}
		case *pbv1.TxRequest_ListObjects:
			v, err := s.OpenFGAServiceServer.ListObjects(ctx, r.ListObjects)
			if err != nil {
				return err
			}
			res = &pbv1.TxResponse{Response: &pbv1.TxResponse_ListObjects{
				ListObjects: v,
			}}
		case *pbv1.TxRequest_ListUsers:
			v, err := s.OpenFGAServiceServer.ListUsers(ctx, r.ListUsers)
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

func (s *svc) Close() {
	s.s.Close()
}
