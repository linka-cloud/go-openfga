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
	"fmt"

	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	parser "github.com/openfga/language/pkg/go/transformer"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"

	"go.linka.cloud/go-openfga/storage"
	"go.linka.cloud/go-openfga/x"
)

const SchemaVersion = "1.2"

type model struct {
	s *store
	c x.OpenFGA
	*rw
}

// func (m *model) Reload(ctx context.Context) error {
// 	res, err := m.s.c.c.ReadAuthorizationModels(ctx, &openfgav1.ReadAuthorizationModelsRequest{StoreId: m.s.id})
// 	if err != nil {
// 		return err
// 	}
// 	if len(res.AuthorizationModels) == 0 {
// 		return errors.New("not found")
// 	}
// 	m.m = res.AuthorizationModels[len(res.AuthorizationModels)-1]
// 	return nil
// }

func (m *model) Tx(ctx context.Context, opts ...storage.TxOption) (Tx, error) {
	c, err := m.c.Tx(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &tx{c: c, rw: &rw{c: c, m: m.m, mid: m.mid, sid: m.sid}}, nil
}

func (m *model) WithTx(t any) Tx {
	c := m.c.WithTx(t)
	return &tx{c: c, rw: &rw{c: c, m: m.m, mid: m.mid, sid: m.sid}}
}

func (m *model) ID() string {
	return m.m.Id
}

func (m *model) Store() Store {
	return m.s
}

func (m *model) Show() (string, error) {
	return parser.TransformJSONProtoToDSL(m.m, parser.WithIncludeSourceInformation(true))
}

func makeContext(kv ...any) (*structpb.Struct, error) {
	m := make(map[string]any)
	for i := 0; i < len(kv); i += 2 {
		k, ok := kv[i].(string)
		if !ok {
			return nil, fmt.Errorf("invalid key %v: must be string", kv[i])
		}
		m[k] = kv[i+1]
	}
	return structpb.NewStruct(m)
}

func CombineModules(dsl ...string) (string, error) {
	mod, err := combineModules(dsl...)
	if err != nil {
		return "", err
	}
	// parser.TransformJSONProtoToDSL fails due to a bug in the parser on the nil `this`
	b, err := protojson.Marshal(mod)
	if err != nil {
		return "", err
	}
	mod.Reset()
	if err := protojson.Unmarshal(b, mod); err != nil {
		return "", err
	}
	out, err := parser.TransformJSONProtoToDSL(mod, parser.WithIncludeSourceInformation(true))
	if err != nil {
		return "", err
	}
	return out, nil
}

func combineModules(dsl ...string) (*openfgav1.AuthorizationModel, error) {
	if len(dsl) == 1 {
		return parser.TransformDSLToProto(dsl[0])
	}
	var mods []parser.ModuleFile
	for _, v := range dsl {
		mods = append(mods, parser.ModuleFile{Contents: v})
	}
	return parser.TransformModuleFilesToModel(mods, SchemaVersion)
}
