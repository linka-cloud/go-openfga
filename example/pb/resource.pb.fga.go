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
