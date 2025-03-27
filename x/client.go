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
	"errors"
	"io"

	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	"google.golang.org/grpc"

	"go.linka.cloud/go-openfga/storage"
	pbv1 "go.linka.cloud/go-openfga/x/pb/v1"
)

var _ openfgav1.OpenFGAServiceClient = (*txc)(nil)

type Tx interface {
	openfgav1.OpenFGAServiceClient
	Commit(ctx context.Context) error
	Close() error
}

type Client interface {
	openfgav1.OpenFGAServiceClient
	Tx(ctx context.Context, opts ...storage.TxOption) (Tx, error)
}

func NewClient(cc grpc.ClientConnInterface) Client {
	return &xclient{c: pbv1.NewOpenFGAXServiceClient(cc), OpenFGAServiceClient: openfgav1.NewOpenFGAServiceClient(cc)}
}

type xclient struct {
	openfgav1.OpenFGAServiceClient
	c pbv1.OpenFGAXServiceClient
}

func (c *xclient) Tx(ctx context.Context, opts ...storage.TxOption) (Tx, error) {
	var o storage.TxOptions
	for _, v := range opts {
		v(&o)
	}
	t, err := c.c.Tx(ctx)
	if err != nil {
		return nil, err
	}
	return &txc{t: t, ro: o.ReadOnly}, nil
}

type txc struct {
	t  grpc.BidiStreamingClient[pbv1.TxRequest, pbv1.TxResponse]
	ro bool
}

func (t *txc) Read(ctx context.Context, in *openfgav1.ReadRequest, opts ...grpc.CallOption) (*openfgav1.ReadResponse, error) {
	if err := t.t.Send(&pbv1.TxRequest{ReadOnly: t.ro, Request: &pbv1.TxRequest_Read{Read: in}}); err != nil {
		return nil, err
	}
	res, err := t.t.Recv()
	if err != nil {
		return nil, err
	}
	v, ok := res.Response.(*pbv1.TxResponse_Read)
	if !ok {
		return nil, errors.New("unexpected response type")
	}
	return v.Read, nil
}

func (t *txc) Write(ctx context.Context, in *openfgav1.WriteRequest, opts ...grpc.CallOption) (*openfgav1.WriteResponse, error) {
	if err := t.t.Send(&pbv1.TxRequest{ReadOnly: t.ro, Request: &pbv1.TxRequest_Write{Write: in}}); err != nil {
		return nil, err
	}
	res, err := t.t.Recv()
	if err != nil {
		return nil, err
	}
	v, ok := res.Response.(*pbv1.TxResponse_Write)
	if !ok {
		return nil, errors.New("unexpected response type")
	}
	return v.Write, nil
}

func (t *txc) Check(ctx context.Context, in *openfgav1.CheckRequest, opts ...grpc.CallOption) (*openfgav1.CheckResponse, error) {
	if err := t.t.Send(&pbv1.TxRequest{ReadOnly: t.ro, Request: &pbv1.TxRequest_Check{Check: in}}); err != nil {
		return nil, err
	}
	res, err := t.t.Recv()
	if err != nil {
		return nil, err
	}
	v, ok := res.Response.(*pbv1.TxResponse_Check)
	if !ok {
		return nil, errors.New("unexpected response type")
	}
	return v.Check, nil
}

func (t *txc) BatchCheck(ctx context.Context, in *openfgav1.BatchCheckRequest, opts ...grpc.CallOption) (*openfgav1.BatchCheckResponse, error) {
	if err := t.t.Send(&pbv1.TxRequest{ReadOnly: t.ro, Request: &pbv1.TxRequest_BatchCheck{BatchCheck: in}}); err != nil {
		return nil, err
	}
	res, err := t.t.Recv()
	if err != nil {
		return nil, err
	}
	v, ok := res.Response.(*pbv1.TxResponse_BatchCheck)
	if !ok {
		return nil, errors.New("unexpected response type")
	}
	return v.BatchCheck, nil
}

func (t *txc) Expand(ctx context.Context, in *openfgav1.ExpandRequest, opts ...grpc.CallOption) (*openfgav1.ExpandResponse, error) {
	if err := t.t.Send(&pbv1.TxRequest{ReadOnly: t.ro, Request: &pbv1.TxRequest_Expand{Expand: in}}); err != nil {
		return nil, err
	}
	res, err := t.t.Recv()
	if err != nil {
		return nil, err
	}
	v, ok := res.Response.(*pbv1.TxResponse_Expand)
	if !ok {
		return nil, errors.New("unexpected response type")
	}
	return v.Expand, nil
}

func (t *txc) ReadAuthorizationModels(ctx context.Context, in *openfgav1.ReadAuthorizationModelsRequest, opts ...grpc.CallOption) (*openfgav1.ReadAuthorizationModelsResponse, error) {
	if err := t.t.Send(&pbv1.TxRequest{ReadOnly: t.ro, Request: &pbv1.TxRequest_ReadAuthorizationModels{ReadAuthorizationModels: in}}); err != nil {
		return nil, err
	}
	res, err := t.t.Recv()
	if err != nil {
		return nil, err
	}
	v, ok := res.Response.(*pbv1.TxResponse_ReadAuthorizationModels)
	if !ok {
		return nil, errors.New("unexpected response type")
	}
	return v.ReadAuthorizationModels, nil
}

func (t *txc) ReadAuthorizationModel(ctx context.Context, in *openfgav1.ReadAuthorizationModelRequest, opts ...grpc.CallOption) (*openfgav1.ReadAuthorizationModelResponse, error) {
	if err := t.t.Send(&pbv1.TxRequest{ReadOnly: t.ro, Request: &pbv1.TxRequest_ReadAuthorizationModel{ReadAuthorizationModel: in}}); err != nil {
		return nil, err
	}
	res, err := t.t.Recv()
	if err != nil {
		return nil, err
	}
	v, ok := res.Response.(*pbv1.TxResponse_ReadAuthorizationModel)
	if !ok {
		return nil, errors.New("unexpected response type")
	}
	return v.ReadAuthorizationModel, nil
}

func (t *txc) WriteAuthorizationModel(ctx context.Context, in *openfgav1.WriteAuthorizationModelRequest, opts ...grpc.CallOption) (*openfgav1.WriteAuthorizationModelResponse, error) {
	if err := t.t.Send(&pbv1.TxRequest{ReadOnly: t.ro, Request: &pbv1.TxRequest_WriteAuthorizationModel{WriteAuthorizationModel: in}}); err != nil {
		return nil, err
	}
	res, err := t.t.Recv()
	if err != nil {
		return nil, err
	}
	v, ok := res.Response.(*pbv1.TxResponse_WriteAuthorizationModel)
	if !ok {
		return nil, errors.New("unexpected response type")
	}
	return v.WriteAuthorizationModel, nil
}

func (t *txc) WriteAssertions(ctx context.Context, in *openfgav1.WriteAssertionsRequest, opts ...grpc.CallOption) (*openfgav1.WriteAssertionsResponse, error) {
	if err := t.t.Send(&pbv1.TxRequest{ReadOnly: t.ro, Request: &pbv1.TxRequest_WriteAssertions{WriteAssertions: in}}); err != nil {
		return nil, err
	}
	res, err := t.t.Recv()
	if err != nil {
		return nil, err
	}
	v, ok := res.Response.(*pbv1.TxResponse_WriteAssertions)
	if !ok {
		return nil, errors.New("unexpected response type")
	}
	return v.WriteAssertions, nil
}

func (t *txc) ReadAssertions(ctx context.Context, in *openfgav1.ReadAssertionsRequest, opts ...grpc.CallOption) (*openfgav1.ReadAssertionsResponse, error) {
	if err := t.t.Send(&pbv1.TxRequest{ReadOnly: t.ro, Request: &pbv1.TxRequest_ReadAssertions{ReadAssertions: in}}); err != nil {
		return nil, err
	}
	res, err := t.t.Recv()
	if err != nil {
		return nil, err
	}
	v, ok := res.Response.(*pbv1.TxResponse_ReadAssertions)
	if !ok {
		return nil, errors.New("unexpected response type")
	}
	return v.ReadAssertions, nil
}

func (t *txc) ReadChanges(ctx context.Context, in *openfgav1.ReadChangesRequest, opts ...grpc.CallOption) (*openfgav1.ReadChangesResponse, error) {
	if err := t.t.Send(&pbv1.TxRequest{ReadOnly: t.ro, Request: &pbv1.TxRequest_ReadChanges{ReadChanges: in}}); err != nil {
		return nil, err
	}
	res, err := t.t.Recv()
	if err != nil {
		return nil, err
	}
	v, ok := res.Response.(*pbv1.TxResponse_ReadChanges)
	if !ok {
		return nil, errors.New("unexpected response type")
	}
	return v.ReadChanges, nil
}

func (t *txc) CreateStore(ctx context.Context, in *openfgav1.CreateStoreRequest, opts ...grpc.CallOption) (*openfgav1.CreateStoreResponse, error) {
	if err := t.t.Send(&pbv1.TxRequest{ReadOnly: t.ro, Request: &pbv1.TxRequest_CreateStore{CreateStore: in}}); err != nil {
		return nil, err
	}
	res, err := t.t.Recv()
	if err != nil {
		return nil, err
	}
	v, ok := res.Response.(*pbv1.TxResponse_CreateStore)
	if !ok {
		return nil, errors.New("unexpected response type")
	}
	return v.CreateStore, nil
}

func (t *txc) UpdateStore(ctx context.Context, in *openfgav1.UpdateStoreRequest, opts ...grpc.CallOption) (*openfgav1.UpdateStoreResponse, error) {
	if err := t.t.Send(&pbv1.TxRequest{ReadOnly: t.ro, Request: &pbv1.TxRequest_UpdateStore{UpdateStore: in}}); err != nil {
		return nil, err
	}
	res, err := t.t.Recv()
	if err != nil {
		return nil, err
	}
	v, ok := res.Response.(*pbv1.TxResponse_UpdateStore)
	if !ok {
		return nil, errors.New("unexpected response type")
	}
	return v.UpdateStore, nil
}

func (t *txc) DeleteStore(ctx context.Context, in *openfgav1.DeleteStoreRequest, opts ...grpc.CallOption) (*openfgav1.DeleteStoreResponse, error) {
	if err := t.t.Send(&pbv1.TxRequest{ReadOnly: t.ro, Request: &pbv1.TxRequest_DeleteStore{DeleteStore: in}}); err != nil {
		return nil, err
	}
	res, err := t.t.Recv()
	if err != nil {
		return nil, err
	}
	v, ok := res.Response.(*pbv1.TxResponse_DeleteStore)
	if !ok {
		return nil, errors.New("unexpected response type")
	}
	return v.DeleteStore, nil
}

func (t *txc) GetStore(ctx context.Context, in *openfgav1.GetStoreRequest, opts ...grpc.CallOption) (*openfgav1.GetStoreResponse, error) {
	if err := t.t.Send(&pbv1.TxRequest{ReadOnly: t.ro, Request: &pbv1.TxRequest_GetStore{GetStore: in}}); err != nil {
		return nil, err
	}
	res, err := t.t.Recv()
	if err != nil {
		return nil, err
	}
	v, ok := res.Response.(*pbv1.TxResponse_GetStore)
	if !ok {
		return nil, errors.New("unexpected response type")
	}
	return v.GetStore, nil
}

func (t *txc) ListStores(ctx context.Context, in *openfgav1.ListStoresRequest, opts ...grpc.CallOption) (*openfgav1.ListStoresResponse, error) {
	if err := t.t.Send(&pbv1.TxRequest{ReadOnly: t.ro, Request: &pbv1.TxRequest_ListStores{ListStores: in}}); err != nil {
		return nil, err
	}
	res, err := t.t.Recv()
	if err != nil {
		return nil, err
	}
	v, ok := res.Response.(*pbv1.TxResponse_ListStores)
	if !ok {
		return nil, errors.New("unexpected response type")
	}
	return v.ListStores, nil
}

func (t *txc) StreamedListObjects(ctx context.Context, in *openfgav1.StreamedListObjectsRequest, opts ...grpc.CallOption) (openfgav1.OpenFGAService_StreamedListObjectsClient, error) {
	return nil, errors.New("unsupported")
}

func (t *txc) ListObjects(ctx context.Context, in *openfgav1.ListObjectsRequest, opts ...grpc.CallOption) (*openfgav1.ListObjectsResponse, error) {
	if err := t.t.Send(&pbv1.TxRequest{ReadOnly: t.ro, Request: &pbv1.TxRequest_ListObjects{ListObjects: in}}); err != nil {
		return nil, err
	}
	res, err := t.t.Recv()
	if err != nil {
		return nil, err
	}
	v, ok := res.Response.(*pbv1.TxResponse_ListObjects)
	if !ok {
		return nil, errors.New("unexpected response type")
	}
	return v.ListObjects, nil
}

func (t *txc) ListUsers(ctx context.Context, in *openfgav1.ListUsersRequest, opts ...grpc.CallOption) (*openfgav1.ListUsersResponse, error) {
	if err := t.t.Send(&pbv1.TxRequest{ReadOnly: t.ro, Request: &pbv1.TxRequest_ListUsers{ListUsers: in}}); err != nil {
		return nil, err
	}
	res, err := t.t.Recv()
	if err != nil {
		return nil, err
	}
	v, ok := res.Response.(*pbv1.TxResponse_ListUsers)
	if !ok {
		return nil, errors.New("unexpected response type")
	}
	return v.ListUsers, nil
}

func (t *txc) Commit(ctx context.Context) error {
	if err := t.t.Send(&pbv1.TxRequest{ReadOnly: t.ro, Request: &pbv1.TxRequest_Commit{Commit: &pbv1.CommitRequest{}}}); err != nil {
		return err
	}
	_, err := t.t.Recv()
	if errors.Is(err, io.EOF) {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

func (t *txc) Close() error {
	return t.t.CloseSend()
}
