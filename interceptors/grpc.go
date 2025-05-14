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

type CheckFunc func(ctx context.Context, req any, user string, kvs ...any) error

type UserFunc func(ctx context.Context) (string, map[string]any, error)

type FGA interface {
	Interceptors
	Register(fqn string, obj CheckFunc)
	// Normalize can be used to normalize object ID
	Normalize(id string) string
	// Check checks if the user has the relation to the object
	Check(ctx context.Context, object, relation, user string, contextKVs ...any) (bool, error)
	// Has checks if the object exists, it can be used to be able to return "not found" errors
	Has(ctx context.Context, object string) (bool, error)
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
		reg:   make(map[string]CheckFunc),
		model: model,
	}
	for _, v := range opts {
		v(fga)
	}
	if fga.user == nil {
		return nil, errors.New("grpc openfga: missing user function")
	}
	if fga.normalize == nil {
		fga.normalize = func(s string) string { return s }
	}
	return fga, nil
}

type fga struct {
	reg       map[string]CheckFunc
	user      UserFunc
	model     openfga.Model
	normalize func(string) string
}

func (f *fga) Normalize(id string) string {
	return f.normalize(id)
}

func (f *fga) Check(ctx context.Context, object, relation, user string, contextKVs ...any) (bool, error) {
	return f.model.Check(ctx, object, relation, user, contextKVs...)
}

func (f *fga) Has(ctx context.Context, object string) (bool, error) {
	out, _, err := f.model.ReadWithPaging(ctx, object, "", "", 1, "")
	return len(out) != 0, err
}

func (f *fga) Register(fqn string, obj CheckFunc) {
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
		return handler(srv, &wrapper{ctx: ctx, ServerStream: ss, method: info.FullMethod, f: f})
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
	return fn(ctx, req, u, kvs...)
}

type wrapper struct {
	grpc.ServerStream
	ctx    context.Context
	method string
	f      *fga
}

func (w *wrapper) Context() context.Context {
	return w.ctx
}

func (w *wrapper) RecvMsg(m any) error {
	if err := w.ServerStream.RecvMsg(m); err != nil {
		return err
	}
	if err := w.f.check(w.ctx, w.method, m); err != nil {
		return err
	}
	return nil
}
