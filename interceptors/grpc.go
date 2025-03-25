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

package interceptors

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.linka.cloud/go-openfga"
)

type RegisterFunc func(fga FGA)

type ObjectFunc func(ctx context.Context, req any) (objectType, objectID, relation string, err error)

type UserFunc func(ctx context.Context) (string, map[string]any, error)

type FGA interface {
	Interceptors
	Register(fqn string, obj ObjectFunc)
}

type Interceptors interface {
	UnaryServerInterceptor() grpc.UnaryServerInterceptor
	StreamServerInterceptor() grpc.StreamServerInterceptor
	UnaryClientInterceptor() grpc.UnaryClientInterceptor
	StreamClientInterceptor() grpc.StreamClientInterceptor
}

func New(_ context.Context, model openfga.Model, opts ...Option) (FGA, error) {
	if model == nil {
		return nil, errors.New("model is nil")
	}
	fga := &fga{
		reg:   make(map[string]ObjectFunc),
		model: model,
	}
	for _, v := range opts {
		v(fga)
	}
	if fga.user == nil {
		return nil, errors.New("grpc openfga: missing user function")
	}
	return fga, nil
}

type fga struct {
	reg       map[string]ObjectFunc
	user      UserFunc
	model     openfga.Model
	normalize func(string) string
}

func (f *fga) Register(fqn string, obj ObjectFunc) {
	f.reg[fqn] = obj
}

func (f *fga) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx = openfga.Context(ctx, f.model)
		if err := f.check(ctx, info.FullMethod, req); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

func (f *fga) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := openfga.Context(ss.Context(), f.model)
		if err := f.check(ss.Context(), info.FullMethod, nil); err != nil {
			return err
		}
		return handler(srv, &wrapper{ctx: ctx, ServerStream: ss})
	}
}

func (f *fga) UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if err := f.check(ctx, method, req); err != nil {
			return err
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (f *fga) StreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		if err := f.check(ctx, method, nil); err != nil {
			return nil, err
		}
		return streamer(ctx, desc, cc, method, opts...)
	}
}

func (f *fga) check(ctx context.Context, fullMethod string, req any) error {
	u, kv, err := f.user(ctx)
	if err != nil {
		return err
	}
	var kvs []any
	for k, v := range kv {
		kvs = append(kvs, k, v)
	}
	fn, ok := f.reg[fullMethod]
	if !ok {
		return status.Errorf(codes.Internal, "permission for '%s' not found", fullMethod)
	}
	t, id, r, err := fn(ctx, req)
	if err != nil {
		return err
	}
	if f.normalize != nil {
		id = f.normalize(id)
	}
	granted, err := f.model.Check(ctx, t+":"+id, r, u, kvs...)
	if err != nil {
		return status.Errorf(codes.Internal, "permission check failed: %v", err)
	}
	if !granted {
		return status.Errorf(codes.PermissionDenied, "[%s]: not allowed to call %s", u, fullMethod)
	}
	return nil
}

type wrapper struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrapper) Context() context.Context {
	return w.ctx
}
