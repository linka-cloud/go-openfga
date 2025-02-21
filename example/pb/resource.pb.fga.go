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

// Code generated by protoc-gen-go-openfga. DO NOT EDIT.
package resource

import (
	"context"
	_ "embed"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	fgainterceptors "go.linka.cloud/go-openfga/interceptors"
)

var (
	_ = codes.OK
	_ = status.New
)

const (
	FGASystemType = "system"

	FGASystemCanCreateResource = "can_create_resource"
	FGASystemCanListResources  = "can_list_resources"
	FGASystemCanWatchResources = "can_watch_resources"
	FGASystemResourceAdmin     = "resource_admin"
	FGASystemResourceReader    = "resource_reader"
	FGASystemResourceWatcher   = "resource_watcher"
	FGASystemResourceWriter    = "resource_writer"

	FGAResourceType = "resource"

	FGAResourceAdmin     = "admin"
	FGAResourceCanDelete = "can_delete"
	FGAResourceCanRead   = "can_read"
	FGAResourceCanUpdate = "can_update"
	FGAResourceReader    = "reader"
	FGAResourceSystem    = "system"
)

// FGASystemObject returns the object string for the system type, e.g. "system:id"
func FGASystemObject(id string) string {
	return FGASystemType + ":" + id
}

// FGAResourceObject returns the object string for the resource type, e.g. "resource:id"
func FGAResourceObject(id string) string {
	return FGAResourceType + ":" + id
}

//go:embed resource.fga
var FGAModel string

// RegisterFGA registers the ResourceService service with the provided FGA interceptors.
func RegisterFGA(fga fgainterceptors.FGA) {
	fga.Register(ResourceService_Create_FullMethodName, func(ctx context.Context, req any) (object string, relation string, err error) {
		return FGASystemType + ":" + "default", FGASystemCanCreateResource, nil
	})
	fga.Register(ResourceService_Read_FullMethodName, func(ctx context.Context, req any) (object string, relation string, err error) {
		r, ok := req.(*ReadRequest)
		if !ok {
			panic("unexpected request type: expected ReadRequest")
		}
		id := r.GetID()
		if id == "" {
			return "", "", status.Error(codes.InvalidArgument, "id is required")
		}
		return FGAResourceObject(id), FGAResourceCanRead, nil
	})
	fga.Register(ResourceService_Update_FullMethodName, func(ctx context.Context, req any) (object string, relation string, err error) {
		r, ok := req.(*UpdateRequest)
		if !ok {
			panic("unexpected request type: expected UpdateRequest")
		}
		id := r.GetResource().GetID()
		if id == "" {
			return "", "", status.Error(codes.InvalidArgument, "resource.id is required")
		}
		return FGAResourceObject(id), FGAResourceCanUpdate, nil
	})
	fga.Register(ResourceService_Delete_FullMethodName, func(ctx context.Context, req any) (object string, relation string, err error) {
		r, ok := req.(*DeleteRequest)
		if !ok {
			panic("unexpected request type: expected DeleteRequest")
		}
		id := r.GetID()
		if id == "" {
			return "", "", status.Error(codes.InvalidArgument, "id is required")
		}
		return FGAResourceObject(id), FGAResourceCanDelete, nil
	})
	fga.Register(ResourceService_List_FullMethodName, func(ctx context.Context, req any) (object string, relation string, err error) {
		return FGASystemType + ":" + "default", FGASystemCanListResources, nil
	})
	fga.Register(ResourceService_Watch_FullMethodName, func(ctx context.Context, req any) (object string, relation string, err error) {
		return FGASystemType + ":" + "default", FGASystemCanWatchResources, nil
	})
}
