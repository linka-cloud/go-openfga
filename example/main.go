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

package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"time"

	"github.com/fullstorydev/grpchan/inprocgrpc"
	"github.com/openfga/openfga/pkg/server"
	"github.com/openfga/openfga/pkg/storage/memory"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"go.linka.cloud/go-openfga"
	example "go.linka.cloud/go-openfga/example/pb"
	"go.linka.cloud/go-openfga/interceptors"
)

//go:embed base.fga
var modelBase string

const userKey = "user"

func userContext(ctx context.Context, user string) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.Pairs(userKey, user))
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mem := memory.New()
	f, err := openfga.New(server.WithDatastore(mem))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	s, err := f.CreateStore(ctx, "default")
	if err != nil {
		log.Fatal(err)
	}

	model, err := s.WriteAuthorizationModel(ctx, modelBase, example.ResourceServiceModel)
	if err != nil {
		log.Fatal(err)
	}

	fga, err := interceptors.New(ctx, model, interceptors.WithUserFunc(func(ctx context.Context) (string, map[string]any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok || len(md.Get(userKey)) == 0 {
			return "", nil, status.Errorf(codes.Unauthenticated, "missing user from metadata")
		}
		return "user:" + md.Get(userKey)[0], nil, nil
	}))
	if err != nil {
		log.Fatal(err)
	}

	// register some users with system roles
	for _, v := range []string{
		example.ResourceServiceRoles.System.ResourceReader,
		example.ResourceServiceRoles.System.ResourceWriter,
		example.ResourceServiceRoles.System.ResourceAdmin,
		example.ResourceServiceRoles.System.ResourceWatcher,
	} {
		if err := model.Write(ctx, "system:default", v, fmt.Sprintf("user:%s", v)); err != nil {
			log.Fatal(err)
		}
	}

	// create the service
	svc := NewResourceService()

	// register the service permissions
	example.RegisterResourceServiceFGA(fga)

	// create the in-process grpc channel
	channel := (&inprocgrpc.Channel{}).
		WithServerUnaryInterceptor(fga.UnaryServerInterceptor()).
		WithServerStreamInterceptor(fga.StreamServerInterceptor())

	// register the service as usual
	example.RegisterResourceServiceServer(channel, svc)

	// create a client
	client := example.NewResourceServiceClient(channel)

	// validate checks
	if _, err := client.List(userContext(ctx, example.ResourceServiceRoles.System.ResourceReader), &example.ListRequest{}); err != nil {
		log.Fatal(err)
	}

	if _, err := client.Create(userContext(ctx, example.ResourceServiceRoles.System.ResourceReader), &example.CreateRequest{Resource: &example.Resource{ID: "0"}}); err == nil {
		log.Fatal("reader should not be able to create")
	}

	if _, err := client.Create(userContext(ctx, example.ResourceServiceRoles.System.ResourceWriter), &example.CreateRequest{Resource: &example.Resource{ID: "0"}}); err != nil {
		log.Fatal(err)
	}

	ss, err := client.Watch(userContext(ctx, example.ResourceServiceRoles.System.ResourceWriter), &example.WatchRequest{})
	if err != nil {
		log.Fatal(err)
	}
	// try to receive an event as the interceptor is not called when creating the stream
	if _, err := ss.Recv(); err == nil {
		log.Fatal("writer should not be able to watch")
	}

	ss, err = client.Watch(userContext(ctx, example.ResourceServiceRoles.System.ResourceAdmin), &example.WatchRequest{})
	if err != nil {
		log.Fatal(err)
	}

	// create a resource to trigger an event
	go func() {
		time.Sleep(100 * time.Millisecond)
		if _, err := client.Create(userContext(ctx, example.ResourceServiceRoles.System.ResourceWriter), &example.CreateRequest{Resource: &example.Resource{ID: "1"}}); err != nil {
			log.Fatal(err)
		}
	}()
	if _, err := ss.Recv(); err != nil {
		log.Fatal(err)
	}
}
