// Copyright 2022 Linka Cloud  All rights reserved.
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

package main

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"

	"go.linka.cloud/pubsub/typed"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"go.linka.cloud/go-openfga"
	pb "go.linka.cloud/go-openfga/example/pb"
)

type service struct {
	pb.UnsafeResourceServiceServer
	store  sync.Map
	pubsub pubsub.Publisher[*pb.Event]
}

func NewResourceService() pb.ResourceServiceServer {
	return &service{pubsub: pubsub.NewPublisher[*pb.Event](time.Second, 1)}
}

func (r *service) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	if err := openfga.MustFromContext(ctx).Write(ctx, pb.FGAResourceObject(req.GetResource().GetID()), pb.FGASystemType, defaultSystem); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to write permission: %v", err)
	}
	r.store.Store(req.GetResource().GetID(), req.GetResource())
	defer r.pubsub.Publish(&pb.Event{Type: pb.EventCreate, Resource: req.Resource})
	return &pb.CreateResponse{Resource: req.GetResource()}, nil
}

func (r *service) Read(_ context.Context, req *pb.ReadRequest) (*pb.ReadResponse, error) {
	v, ok := r.store.Load(req.GetID())
	if !ok {
		return nil, status.Errorf(codes.NotFound, "%s does not exists", req.GetID())
	}
	return &pb.ReadResponse{Resource: v.(*pb.Resource)}, nil
}

func (r *service) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	r.store.Store(req.GetResource().GetID(), req.GetResource())
	defer r.pubsub.Publish(&pb.Event{Type: pb.EventUpdate, Resource: req.Resource})
	return &pb.UpdateResponse{Resource: req.GetResource()}, nil
}

func (r *service) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	if err := openfga.MustFromContext(ctx).Delete(ctx, pb.FGAResourceObject(req.GetID()), pb.FGASystemType, defaultSystem); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete permission: %v", err)
	}
	v, ok := r.store.Load(req.GetID())
	if !ok {
		return nil, status.Errorf(codes.NotFound, "%s does not exists", req.GetID())
	}
	defer r.pubsub.Publish(&pb.Event{Type: pb.EventDelete, Resource: v.(*pb.Resource)})
	r.store.Delete(req.GetID())
	return &pb.DeleteResponse{}, nil
}

func (r *service) List(ctx context.Context, _ *pb.ListRequest) (*pb.ListResponse, error) {
	ids, err := openfga.MustFromContext(ctx).ListObjects(ctx, pb.FGAResourceType, pb.FGAResourceCanRead, mustContextUser(ctx))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list resources: %v", err)
	}
	res := &pb.ListResponse{}
	for _, v := range ids {
		if v, ok := r.store.Load(v); ok {
			res.Resources = append(res.Resources, v.(*pb.Resource))
		}
	}
	return res, nil
}

func (r *service) Watch(_ *pb.WatchRequest, server grpc.ServerStreamingServer[pb.Event]) error {
	ch := r.pubsub.Subscribe()
	defer r.pubsub.Evict(ch)
	for {
		select {
		case e := <-ch:
			ok, err := openfga.MustFromContext(server.Context()).Check(
				server.Context(), pb.FGAResourceObject(e.GetResource().GetID()), pb.FGAResourceCanRead, mustContextUser(server.Context()))
			if err != nil {
				return status.Errorf(codes.Internal, "permission check failed: %v", err)
			}
			if !ok {
				continue
			}
			if err := server.Send(proto.Clone(e).(*pb.Event)); err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}
				return err
			}
		case <-server.Context().Done():
			return server.Context().Err()
		}
	}
}
