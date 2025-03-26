# go-openfga & protoc-gen-go-openfga

[![Language: Go](https://img.shields.io/badge/lang-Go-6ad7e5.svg?style=flat-square&logo=go)](https://golang.org/)
[![Go Reference](https://pkg.go.dev/badge/go.linka.cloud/go-openfga.svg)](https://pkg.go.dev/go.linka.cloud/go-openfga)
[![Go Report Card](https://goreportcard.com/badge/go.linka.cloud/go-openfga)](https://goreportcard.com/report/go.linka.cloud/go-openfga)

**Project status: *alpha***

Not all planned features are completed.
The API, spec, status and other user facing objects are subject to change.
We do not support backward-compatibility for the alpha releases.


## `go-openfga`

### Overview

`go-openfga` is a simple implementation of an easy to use (and currently partial) client for the [OpenFGA](https://github.com/openfga/openfga) gRPC API.

It also provides a simple way to run openfga in process.

#### Import

```go
import openfga "go.linka.cloud/go-openfga"
```

## `protoc-gen-go-openfga`

### Overview

`protoc-gen-go-openfga` is a protoc plugin that generates openfga schema and go code to register access checks into the interceptor.


### Installation

```shell
go install go.linka.cloud/go-openfga/cmd/protoc-gen-go-openfga
```

### Usage

Use the plugin as any other protoc plugins.

### Generated code

For a given base [`openfga` module](example/base.fga):

```openfga
module base

type system
  relations
    define admin: [user]
    define writer: [user] or admin
    define reader: [user] or admin
    define watcher: [user] or admin

type user

```

For a given [`resource.proto`](example/pb/resource.proto):

```protobuf
syntax = "proto3";

package resource;

option go_package = "./resource";

import "openfga/openfga.proto";
import "patch/go.proto";

option (go.lint).all = true;

import "example/pb/types.proto";

service ResourceService {
  option (openfga.defaults) = { type: "system", id: "default" };
  option (openfga.module) = {
    name: "resource",
    extends: [ {
      type: "system",
      relations: [
        { define: "resource_admin", as: "[user, user with non_expired_grant] or admin" },
        { define: "resource_writer", as: "[user] or resource_admin" },
        { define: "resource_reader", as: "[user] or resource_admin or reader" },
        { define: "resource_watcher", as: "[user] or resource_admin or watcher" }
      ]
    } ],
    definitions: [ {
      type: "resource",
      relations: [
        { define: "system", as: "[system]" },
        { define: "admin", as: "[user] or resource_admin from system" },
        { define: "reader", as: "[user] or resource_reader from system" }
      ]
    }, {
      type: "sub",
      relations: [
        { define: "resource", as: "[resource]" },
        { define: "admin", as: "[user] or admin from resource" },
        { define: "reader", as: "[user] or reader from resource" }
      ]
    } ],
    conditions: [ "non_expired_grant(current_time: timestamp, grant_time: timestamp, grant_duration: duration) { current_time < grant_time + grant_duration }" ]
  };
  rpc Create (CreateRequest) returns (CreateResponse) {
    option (openfga.access) = { check: [ { as: "resource_writer" } ] };
  };
  rpc Read (ReadRequest) returns (ReadResponse) {
    option (openfga.access) = { check: [ { as: "resource_reader" }, { type: "resource", id: "{id}" as: "reader" } ] };
  }
  rpc Update (UpdateRequest) returns (UpdateResponse) {
    option (openfga.access) = { check: [ { as: "resource_writer" }, { type: "resource", id: "{resource.id}" as: "admin" } ] };
  }
  rpc AddSub (AddSubRequest) returns (AddSubResponse) {
    option (openfga.access) = { check: [
//      { as: "resource_writer" },
      { type: "resource", id: "{id}" as: "admin" }
    ] };
  }
  rpc ReadSub (ReadSubRequest) returns (ReadSubResponse) {
    option (openfga.access) = { check: [
//      { as: "resource_reader" },
      { type: "resource", id: "{resource_id}" as: "reader", ignore_not_found: true },
      { type: "sub", id: "{id}" as: "reader" }
    ] };
  }
  rpc Delete(DeleteRequest) returns (DeleteResponse) {
    option (openfga.access) = { check: [ { as: "resource_writer" }, { type: "resource", id: "{id}" as: "admin" } ] };
  }
  rpc List(ListRequest) returns (ListResponse) {
    option (openfga.access) = { check: [ { as: "resource_reader" } ] };
  }
  rpc Watch(WatchRequest) returns (stream Event) {
    option (openfga.access) = { check: [ { as: "resource_watcher" } ] };
  }
}


```

The following [`resource.fga`](example/pb/resource.fga) `openfga` module will be generated:

```openfga
# Code generated by protoc-gen-go-openfga. DO NOT EDIT.

module resource

extend type system
  relations
    define resource_admin: [user, user with non_expired_grant] or admin
    define resource_writer: [user] or resource_admin
    define resource_reader: [user] or resource_admin or reader
    define resource_watcher: [user] or resource_admin or watcher
    define can_resource_create: resource_writer
    define can_resource_read: resource_reader
    define can_resource_update: resource_writer
    define can_resource_delete: resource_writer
    define can_resource_list: resource_reader
    define can_resource_watch: resource_watcher

type resource
  relations
    define system: [system]
    define admin: [user] or resource_admin from system
    define reader: [user] or resource_reader from system
    define can_read: reader
    define can_update: admin
    define can_add_sub: admin
    define can_read_sub: reader
    define can_delete: admin
type sub
  relations
    define resource: [resource]
    define admin: [user] or admin from resource
    define reader: [user] or reader from resource
    define can_read: reader

condition non_expired_grant(current_time: timestamp, grant_time: timestamp, grant_duration: duration) { current_time < grant_time + grant_duration }

```

And following code will be generated:

```go
// Code generated by protoc-gen-go-openfga. DO NOT EDIT.
package resource

import (
	"context"
	_ "embed"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	fgainterceptors "go.linka.cloud/go-openfga/interceptors"
)

var (
	_ = codes.OK
	_ = status.New
	_ = fmt.Sprintf
	_ = context.Canceled
)

const (
	FGASystemType = "system"

	FGASystemCanResourceCreate = "can_resource_create"
	FGASystemCanResourceDelete = "can_resource_delete"
	FGASystemCanResourceList   = "can_resource_list"
	FGASystemCanResourceRead   = "can_resource_read"
	FGASystemCanResourceUpdate = "can_resource_update"
	FGASystemCanResourceWatch  = "can_resource_watch"
	FGASystemResourceAdmin     = "resource_admin"
	FGASystemResourceReader    = "resource_reader"
	FGASystemResourceWatcher   = "resource_watcher"
	FGASystemResourceWriter    = "resource_writer"

	FGAResourceType = "resource"

	FGAResourceAdmin      = "admin"
	FGAResourceCanAddSub  = "can_add_sub"
	FGAResourceCanDelete  = "can_delete"
	FGAResourceCanRead    = "can_read"
	FGAResourceCanReadSub = "can_read_sub"
	FGAResourceCanUpdate  = "can_update"
	FGAResourceReader     = "reader"
	FGAResourceSystem     = "system"

	FGASubType = "sub"

	FGASubAdmin    = "admin"
	FGASubCanRead  = "can_read"
	FGASubReader   = "reader"
	FGASubResource = "resource"
)

// FGASystemObject returns the object string for the system type, e.g. "system:id"
func FGASystemObject(id string) string {
	return FGASystemType + ":" + id
}

// FGAResourceObject returns the object string for the resource type, e.g. "resource:id"
func FGAResourceObject(id string) string {
	return FGAResourceType + ":" + id
}

// FGASubObject returns the object string for the sub type, e.g. "sub:id"
func FGASubObject(id string) string {
	return FGASubType + ":" + id
}

//go:embed resource.fga
var FGAModel string

// RegisterFGA registers the ResourceService service with the provided FGA interceptors.
func RegisterFGA(fga fgainterceptors.FGA) {
	fga.Register(ResourceService_Create_FullMethodName, func(ctx context.Context, req any, user string, kvs ...any) error {
		{
			object := "system" + ":" + fga.Normalize("default")
			msg := fmt.Sprintf("[%s]: not allowed to call %s", user, ResourceService_Create_FullMethodName)
			granted, err := fga.Check(ctx, object, FGASystemCanResourceCreate, user, kvs...)
			if err != nil {
				return status.Errorf(codes.Internal, "permission check failed: %v", err)
			}
			if !granted {
				return status.Error(codes.PermissionDenied, msg)
			}
		}
		return nil
	})
	fga.Register(ResourceService_Read_FullMethodName, func(ctx context.Context, req any, user string, kvs ...any) error {
		{
			object := "system" + ":" + fga.Normalize("default")
			msg := fmt.Sprintf("[%s]: not allowed to call %s", user, ResourceService_Read_FullMethodName)
			granted, err := fga.Check(ctx, object, FGASystemCanResourceRead, user, kvs...)
			if err != nil {
				return status.Errorf(codes.Internal, "permission check failed: %v", err)
			}
			if !granted {
				return status.Error(codes.PermissionDenied, msg)
			}
		}
		{
			r, ok := req.(*ReadRequest)
			if !ok {
				panic("unexpected request type: expected ReadRequest")
			}
			id := r.GetID()
			if id == "" {
				return status.Error(codes.InvalidArgument, "id is required")
			}
			object := "resource" + ":" + fga.Normalize(id)
			ok, err := fga.Has(ctx, object)
			if err != nil {
				return status.Errorf(codes.Internal, "permission check failed: %v", err)
			}
			if !ok {
				return status.Errorf(codes.NotFound, "resource %q not found", id)
			}
			msg := fmt.Sprintf("[%s]: not allowed to call %s on resource %q", user, ResourceService_Read_FullMethodName, id)
			granted, err := fga.Check(ctx, object, FGAResourceCanRead, user, kvs...)
			if err != nil {
				return status.Errorf(codes.Internal, "permission check failed: %v", err)
			}
			if !granted {
				return status.Error(codes.PermissionDenied, msg)
			}
		}
		return nil
	})
	fga.Register(ResourceService_Update_FullMethodName, func(ctx context.Context, req any, user string, kvs ...any) error {
		{
			object := "system" + ":" + fga.Normalize("default")
			msg := fmt.Sprintf("[%s]: not allowed to call %s", user, ResourceService_Update_FullMethodName)
			granted, err := fga.Check(ctx, object, FGASystemCanResourceUpdate, user, kvs...)
			if err != nil {
				return status.Errorf(codes.Internal, "permission check failed: %v", err)
			}
			if !granted {
				return status.Error(codes.PermissionDenied, msg)
			}
		}
		{
			r, ok := req.(*UpdateRequest)
			if !ok {
				panic("unexpected request type: expected UpdateRequest")
			}
			id := r.GetResource().GetID()
			if id == "" {
				return status.Error(codes.InvalidArgument, "resource.id is required")
			}
			object := "resource" + ":" + fga.Normalize(id)
			ok, err := fga.Has(ctx, object)
			if err != nil {
				return status.Errorf(codes.Internal, "permission check failed: %v", err)
			}
			if !ok {
				return status.Errorf(codes.NotFound, "resource %q not found", id)
			}
			msg := fmt.Sprintf("[%s]: not allowed to call %s on resource %q", user, ResourceService_Update_FullMethodName, id)
			granted, err := fga.Check(ctx, object, FGAResourceCanUpdate, user, kvs...)
			if err != nil {
				return status.Errorf(codes.Internal, "permission check failed: %v", err)
			}
			if !granted {
				return status.Error(codes.PermissionDenied, msg)
			}
		}
		return nil
	})
	fga.Register(ResourceService_AddSub_FullMethodName, func(ctx context.Context, req any, user string, kvs ...any) error {
		{
			r, ok := req.(*AddSubRequest)
			if !ok {
				panic("unexpected request type: expected AddSubRequest")
			}
			id := r.GetID()
			if id == "" {
				return status.Error(codes.InvalidArgument, "id is required")
			}
			object := "resource" + ":" + fga.Normalize(id)
			ok, err := fga.Has(ctx, object)
			if err != nil {
				return status.Errorf(codes.Internal, "permission check failed: %v", err)
			}
			if !ok {
				return status.Errorf(codes.NotFound, "resource %q not found", id)
			}
			msg := fmt.Sprintf("[%s]: not allowed to call %s on resource %q", user, ResourceService_AddSub_FullMethodName, id)
			granted, err := fga.Check(ctx, object, FGAResourceCanAddSub, user, kvs...)
			if err != nil {
				return status.Errorf(codes.Internal, "permission check failed: %v", err)
			}
			if !granted {
				return status.Error(codes.PermissionDenied, msg)
			}
		}
		return nil
	})
	fga.Register(ResourceService_ReadSub_FullMethodName, func(ctx context.Context, req any, user string, kvs ...any) error {
		{
			r, ok := req.(*ReadSubRequest)
			if !ok {
				panic("unexpected request type: expected ReadSubRequest")
			}
			id := r.GetResourceID()
			if id == "" {
				return status.Error(codes.InvalidArgument, "resource_id is required")
			}
			object := "resource" + ":" + fga.Normalize(id)
			msg := fmt.Sprintf("[%s]: not allowed to call %s on resource %q", user, ResourceService_ReadSub_FullMethodName, id)
			granted, err := fga.Check(ctx, object, FGAResourceCanReadSub, user, kvs...)
			if err != nil {
				return status.Errorf(codes.Internal, "permission check failed: %v", err)
			}
			if !granted {
				return status.Error(codes.PermissionDenied, msg)
			}
		}
		{
			r, ok := req.(*ReadSubRequest)
			if !ok {
				panic("unexpected request type: expected ReadSubRequest")
			}
			id := r.GetID()
			if id == "" {
				return status.Error(codes.InvalidArgument, "id is required")
			}
			object := "sub" + ":" + fga.Normalize(id)
			ok, err := fga.Has(ctx, object)
			if err != nil {
				return status.Errorf(codes.Internal, "permission check failed: %v", err)
			}
			if !ok {
				return status.Errorf(codes.NotFound, "sub %q not found", id)
			}
			msg := fmt.Sprintf("[%s]: not allowed to call %s on sub %q", user, ResourceService_ReadSub_FullMethodName, id)
			granted, err := fga.Check(ctx, object, FGASubCanRead, user, kvs...)
			if err != nil {
				return status.Errorf(codes.Internal, "permission check failed: %v", err)
			}
			if !granted {
				return status.Error(codes.PermissionDenied, msg)
			}
		}
		return nil
	})
	fga.Register(ResourceService_Delete_FullMethodName, func(ctx context.Context, req any, user string, kvs ...any) error {
		{
			object := "system" + ":" + fga.Normalize("default")
			msg := fmt.Sprintf("[%s]: not allowed to call %s", user, ResourceService_Delete_FullMethodName)
			granted, err := fga.Check(ctx, object, FGASystemCanResourceDelete, user, kvs...)
			if err != nil {
				return status.Errorf(codes.Internal, "permission check failed: %v", err)
			}
			if !granted {
				return status.Error(codes.PermissionDenied, msg)
			}
		}
		{
			r, ok := req.(*DeleteRequest)
			if !ok {
				panic("unexpected request type: expected DeleteRequest")
			}
			id := r.GetID()
			if id == "" {
				return status.Error(codes.InvalidArgument, "id is required")
			}
			object := "resource" + ":" + fga.Normalize(id)
			ok, err := fga.Has(ctx, object)
			if err != nil {
				return status.Errorf(codes.Internal, "permission check failed: %v", err)
			}
			if !ok {
				return status.Errorf(codes.NotFound, "resource %q not found", id)
			}
			msg := fmt.Sprintf("[%s]: not allowed to call %s on resource %q", user, ResourceService_Delete_FullMethodName, id)
			granted, err := fga.Check(ctx, object, FGAResourceCanDelete, user, kvs...)
			if err != nil {
				return status.Errorf(codes.Internal, "permission check failed: %v", err)
			}
			if !granted {
				return status.Error(codes.PermissionDenied, msg)
			}
		}
		return nil
	})
	fga.Register(ResourceService_List_FullMethodName, func(ctx context.Context, req any, user string, kvs ...any) error {
		{
			object := "system" + ":" + fga.Normalize("default")
			msg := fmt.Sprintf("[%s]: not allowed to call %s", user, ResourceService_List_FullMethodName)
			granted, err := fga.Check(ctx, object, FGASystemCanResourceList, user, kvs...)
			if err != nil {
				return status.Errorf(codes.Internal, "permission check failed: %v", err)
			}
			if !granted {
				return status.Error(codes.PermissionDenied, msg)
			}
		}
		return nil
	})
	fga.Register(ResourceService_Watch_FullMethodName, func(ctx context.Context, req any, user string, kvs ...any) error {
		{
			object := "system" + ":" + fga.Normalize("default")
			msg := fmt.Sprintf("[%s]: not allowed to call %s", user, ResourceService_Watch_FullMethodName)
			granted, err := fga.Check(ctx, object, FGASystemCanResourceWatch, user, kvs...)
			if err != nil {
				return status.Errorf(codes.Internal, "permission check failed: %v", err)
			}
			if !granted {
				return status.Error(codes.PermissionDenied, msg)
			}
		}
		return nil
	})
}

```

### Usage

See the [example](example) directory for complete example.

```go
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
	pb "go.linka.cloud/go-openfga/example/pb"
	"go.linka.cloud/go-openfga/interceptors"
)

//go:embed base.fga
var modelBase string

// userKey is the key used to store the user in the context metadata
const userKey = "user"

// defaultSystem is the default system object
var defaultSystem = pb.FGASystemObject("default")

// userContext returns a new context with the user set in the metadata
func userContext(ctx context.Context, user string) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.Pairs(userKey, user))
}

// contextUser returns the user from the context metadata
func contextUser(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md.Get(userKey)) == 0 {
		return "", status.Errorf(codes.Unauthenticated, "missing user from metadata")
	}
	return "user:" + md.Get(userKey)[0], nil
}

// mustContextUser returns the user from the context metadata or panics
func mustContextUser(ctx context.Context) string {
	user, err := contextUser(ctx)
	if err != nil {
		panic(err)
	}
	return user
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create the in-memory openfga server
	mem := memory.New()
	f, err := openfga.New(server.WithDatastore(mem))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// create the store
	s, err := f.CreateStore(ctx, "default")
	if err != nil {
		log.Fatal(err)
	}

	// write the model
	model, err := s.WriteAuthorizationModel(ctx, modelBase, pb.FGAModel)
	if err != nil {
		log.Fatal(err)
	}

	// create the interceptors
	fga, err := interceptors.New(ctx, model, interceptors.WithUserFunc(func(ctx context.Context) (string, map[string]any, error) {
		user, err := contextUser(ctx)
		if err != nil {
			return "", nil, err
		}
		return user, nil, nil
	}))
	if err != nil {
		log.Fatal(err)
	}

	// register some users with system roles
	for _, v := range []string{
		pb.FGASystemResourceReader,
		pb.FGASystemResourceWriter,
		pb.FGASystemResourceAdmin,
		pb.FGASystemResourceWatcher,
	} {
		if err := model.Write(ctx, defaultSystem, v, fmt.Sprintf("user:%s", v)); err != nil {
			log.Fatal(err)
		}
	}

	// create the service
	svc := NewResourceService()

	// register the service permissions
	pb.RegisterFGA(fga)

	// create the in-process grpc channel
	channel := (&inprocgrpc.Channel{}).
		WithServerUnaryInterceptor(fga.UnaryServerInterceptor()).
		WithServerStreamInterceptor(fga.StreamServerInterceptor())

	// register the service as usual
	pb.RegisterResourceServiceServer(channel, svc)

	// create a client
	client := pb.NewResourceServiceClient(channel)

	// validate checks
	if _, err := client.List(userContext(ctx, pb.FGASystemResourceReader), &pb.ListRequest{}); err != nil {
		log.Fatal(err)
	}

	if _, err := client.Create(userContext(ctx, pb.FGASystemResourceReader), &pb.CreateRequest{Resource: &pb.Resource{ID: "0"}}); err == nil {
		log.Fatal("reader should not be able to create")
	}

	if _, err := client.Create(userContext(ctx, pb.FGASystemResourceWriter), &pb.CreateRequest{Resource: &pb.Resource{ID: "0"}}); err != nil {
		log.Fatal(err)
	}

	wctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	ss, err := client.Watch(userContext(wctx, pb.FGASystemResourceWriter), &pb.WatchRequest{})
	if err != nil {
		log.Fatal(err)
	}
	// try to receive an event as the interceptor is not called when creating the stream
	if _, err := ss.Recv(); err == nil {
		log.Fatal("writer should not be able to watch")
	}

	wctx, cancel = context.WithTimeout(ctx, time.Second)
	defer cancel()
	ss, err = client.Watch(userContext(wctx, pb.FGASystemResourceAdmin), &pb.WatchRequest{})
	if err != nil {
		log.Fatal(err)
	}

	// create a resource to trigger an event
	go func() {
		time.Sleep(100 * time.Millisecond)
		if _, err := client.Create(userContext(ctx, pb.FGASystemResourceWriter), &pb.CreateRequest{Resource: &pb.Resource{ID: "1"}}); err != nil {
			log.Fatal(err)
		}
	}()
	if _, err := ss.Recv(); err != nil {
		log.Fatal(err)
	}
}

```
