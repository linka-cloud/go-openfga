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

	"go.linka.cloud/pubsub"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"go.linka.cloud/go-openfga/example/pb"
)

type resourceService struct {
	example.UnsafeResourceServiceServer
	store  sync.Map
	pubsub pubsub.Publisher
}

func NewResourceService() example.ResourceServiceServer {
	return &resourceService{pubsub: pubsub.NewPublisher(time.Second, 1)}
}

func (r *resourceService) Create(ctx context.Context, request *example.CreateRequest) (*example.CreateResponse, error) {
	r.store.Store(request.GetResource().GetID(), request.GetResource())
	defer r.pubsub.Publish(&example.Event{Type: example.EventCreate, Resource: request.Resource})
	return &example.CreateResponse{Resource: request.GetResource()}, nil
}

func (r *resourceService) Read(ctx context.Context, request *example.ReadRequest) (*example.ReadResponse, error) {
	v, ok := r.store.Load(request.GetID())
	if !ok {
		return nil, status.Errorf(codes.NotFound, "%s does not exists", request.GetID())
	}
	return &example.ReadResponse{Resource: v.(*example.Resource)}, nil
}

func (r *resourceService) Update(ctx context.Context, request *example.UpdateRequest) (*example.UpdateResponse, error) {
	r.store.Store(request.GetResource().GetID(), request.GetResource())
	defer r.pubsub.Publish(&example.Event{Type: example.EventUpdate, Resource: request.Resource})
	return &example.UpdateResponse{Resource: request.GetResource()}, nil
}

func (r *resourceService) Delete(ctx context.Context, request *example.DeleteRequest) (*example.DeleteResponse, error) {
	v, ok := r.store.Load(request.GetID())
	if !ok {
		return nil, status.Errorf(codes.NotFound, "%s does not exists", request.GetID())
	}
	defer r.pubsub.Publish(&example.Event{Type: example.EventDelete, Resource: v.(*example.Resource)})
	r.store.Delete(request.GetID())
	return &example.DeleteResponse{}, nil
}

func (r *resourceService) List(ctx context.Context, request *example.ListRequest) (*example.ListResponse, error) {
	res := &example.ListResponse{}
	r.store.Range(func(key, value interface{}) bool {
		res.Resources = append(res.Resources, value.(*example.Resource))
		return true
	})
	return res, nil
}

func (r *resourceService) Watch(_ *example.WatchRequest, server grpc.ServerStreamingServer[example.Event]) error {
	ch := r.pubsub.Subscribe()
	defer r.pubsub.Evict(ch)
	for {
		select {
		case e := <-ch:
			if err := server.Send(proto.Clone(e.(*example.Event)).(*example.Event)); err != nil {
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
