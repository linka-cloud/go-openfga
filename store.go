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

package openfga

import (
	"context"
	"errors"
	"fmt"
	"time"

	openfgav1 "github.com/openfga/api/proto/openfga/v1"
)

type store[T any] struct {
	id        string
	name      string
	createdAt time.Time
	updatedAt time.Time
	c         *client[T]
}

func (s *store[T]) ID() string {
	return s.id
}

func (s *store[T]) Name() string {
	return s.name
}

func (s *store[T]) CreatedAt() time.Time {
	return s.createdAt
}

func (s *store[T]) UpdatedAt() time.Time {
	return s.updatedAt
}

func (s *store[T]) AuthorizationModel(ctx context.Context, id string) (Model[T], error) {
	res, err := s.c.c.ReadAuthorizationModel(ctx, &openfgav1.ReadAuthorizationModelRequest{StoreId: s.id, Id: id})
	if err != nil {
		return nil, err
	}
	return s.model(res.AuthorizationModel), nil
}

func (s *store[T]) ListAuthorizationModels(ctx context.Context) ([]Model[T], error) {
	res, err := s.c.c.ReadAuthorizationModels(ctx, &openfgav1.ReadAuthorizationModelsRequest{StoreId: s.id})
	if err != nil {
		return nil, fmt.Errorf("failed to list authorization models: %w", err)
	}
	var models []Model[T]
	for _, m := range res.AuthorizationModels {
		models = append(models, s.model(m))
	}
	return models, nil
}

func (s *store[T]) LastAuthorizationModel(ctx context.Context) (Model[T], error) {
	res, err := s.c.c.ReadAuthorizationModels(ctx, &openfgav1.ReadAuthorizationModelsRequest{StoreId: s.id})
	if err != nil {
		return nil, err
	}
	if len(res.AuthorizationModels) == 0 {
		return nil, errors.New("not found")
	}
	return s.model(res.AuthorizationModels[len(res.AuthorizationModels)-1]), nil
}

func (s *store[T]) WriteAuthorizationModel(ctx context.Context, dsl ...string) (Model[T], error) {
	m, err := combineModules(dsl...)
	if err != nil {
		return nil, fmt.Errorf("failed to parse model: %w", err)
	}
	res, err := s.c.c.WriteAuthorizationModel(ctx, &openfgav1.WriteAuthorizationModelRequest{
		StoreId:         s.id,
		TypeDefinitions: m.TypeDefinitions,
		SchemaVersion:   m.SchemaVersion,
		Conditions:      m.Conditions,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to write authorization model: %w", err)
	}
	return s.AuthorizationModel(ctx, res.AuthorizationModelId)
}

func (s *store[T]) model(m *openfgav1.AuthorizationModel) *model[T] {
	return &model[T]{s: s, rw: &rw{c: s.c.c, m: m, sid: s.id, mid: m.Id}, c: s.c.c}
}
