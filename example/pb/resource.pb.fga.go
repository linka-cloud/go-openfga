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
package example

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

//go:embed resource.fga
var ResourceServiceModel string

var ResourceServiceRoles = struct {
	System struct {
		ResourceAdmin     string
		ResourceWriter    string
		ResourceReader    string
		ResourceWatcher   string
		CanCreateResource string
		CanListResources  string
		CanWatchResources string
	}
	Resource struct {
		System    string
		Admin     string
		Reader    string
		CanRead   string
		CanUpdate string
		CanDelete string
	}
}{
	System: struct {
		ResourceAdmin     string
		ResourceWriter    string
		ResourceReader    string
		ResourceWatcher   string
		CanCreateResource string
		CanListResources  string
		CanWatchResources string
	}{
		ResourceAdmin:     "resource_admin",
		ResourceWriter:    "resource_writer",
		ResourceReader:    "resource_reader",
		ResourceWatcher:   "resource_watcher",
		CanCreateResource: "can_create_resource",
		CanListResources:  "can_list_resources",
		CanWatchResources: "can_watch_resources",
	},
	Resource: struct {
		System    string
		Admin     string
		Reader    string
		CanRead   string
		CanUpdate string
		CanDelete string
	}{
		System:    "system",
		Admin:     "admin",
		Reader:    "reader",
		CanRead:   "can_read",
		CanUpdate: "can_update",
		CanDelete: "can_delete",
	},
}

func RegisterResourceServiceFGA(fga fgainterceptors.FGA) {
	fga.Register(ResourceService_Create_FullMethodName, func(ctx context.Context, req any) (object string, relation string, err error) {
		return "system:default", "can_create_resource", nil
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
		return "resource:" + id, "can_read", nil
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
		return "resource:" + id, "can_update", nil
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
		return "resource:" + id, "can_delete", nil
	})
	fga.Register(ResourceService_List_FullMethodName, func(ctx context.Context, req any) (object string, relation string, err error) {
		return "system:default", "can_list_resources", nil
	})
	fga.Register(ResourceService_Watch_FullMethodName, func(ctx context.Context, req any) (object string, relation string, err error) {
		return "system:default", "can_watch_resources", nil
	})
}
