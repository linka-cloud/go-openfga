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
	_ "embed"
	"fmt"
	"strings"
	"text/template"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"

	openfgav1alpha1 "go.linka.cloud/go-openfga/openfga/v1alpha1"
)

var _ pgs.Module = (*Module)(nil)

//go:embed register.go.tpl
var registerTemplate string

//go:embed module.fga.tpl
var moduleTemplate string

var fgaTemplate = template.Must(template.New("fga").Parse(moduleTemplate))

type Module struct {
	*pgs.ModuleBase
	ctx pgsgo.Context
	tpl *template.Template
}

func FGA() *Module {
	return &Module{
		ModuleBase: &pgs.ModuleBase{},
	}
}

func (m *Module) InitContext(c pgs.BuildContext) {
	m.ModuleBase.InitContext(c)
	m.ctx = pgsgo.InitContext(c.Parameters())
	m.tpl = template.Must(template.New("register").Funcs(map[string]any{
		"package": m.ctx.PackageName,
		"name":    m.ctx.Name,
		"file": func(s pgs.Service, ext string) string {
			return s.File().InputPath().BaseName() + ext
		},
		"comment": func(s string) string {
			var out string
			parts := strings.Split(s, "\n")
			for i, v := range parts {
				if i == len(parts)-1 && v == "" {
					return out
				}
				out += "//" + v + "\n"
			}
			return out
		},
		"module": func(s pgs.Service) *openfgav1alpha1.Module {
			var mod openfgav1alpha1.Module
			ok, err := s.Extension(openfgav1alpha1.ExtModule, &mod)
			if err != nil {
				m.Fail(err)
			}
			if !ok {
				return nil
			}
			if err := mod.ValidateAll(); err != nil {
				m.Fail(err)
			}
			return &mod
		},
		"access": func(me pgs.Method) *openfgav1alpha1.Access {
			var access openfgav1alpha1.Access
			ok, err := me.Extension(openfgav1alpha1.ExtAccess, &access)
			if err != nil {
				m.Fail(err)
			}
			if !ok {
				return nil
			}
			if err := access.ValidateAll(); err != nil {
				m.Fail(err)
			}
			return &access
		},
		"need_getter": func(a *openfgav1alpha1.Access) bool {
			return len(a.ID) >= 2 && a.ID[0] == '{' && a.ID[len(a.ID)-1] == '}'
		},
		"field": func(a *openfgav1alpha1.Access) string {
			return a.ID[1 : len(a.ID)-1]
		},
		"getter": func(a *openfgav1alpha1.Access, me pgs.Method) string {
			t := me.Input()
			parts := strings.Split(a.ID[1:len(a.ID)-1], ".")
			var getter string
			i := 0
			for i < len(parts) {
				if t == nil {
					m.Failf("field %q is not a message", parts[i])
				}
				var f pgs.Field
				for _, v := range t.Fields() {
					if v.Name().String() == parts[i] {
						f = v
						break
					}
				}
				if f == nil {
					m.Failf("field %q not found in %s", parts[i], t.Name())
				}
				g := fmt.Sprintf("Get%s", m.ctx.Name(f))
				if getter == "" {
					getter = g
				} else {
					getter += "()." + g
				}
				t = f.Type().Embed()
				i++
			}
			return getter
		},
		"object": func(a *openfgav1alpha1.Access) string {
			return a.Type + ":" + a.ID
		},
		"upperCamelCase": func(s string) string {
			return pgs.Name(s).UpperCamelCase().String()
		},
	}).Parse(registerTemplate))
}

func (m *Module) Execute(targets map[string]pgs.File, _ map[string]pgs.Package) []pgs.Artifact {
	for _, f := range targets {
		m.generate(f)
	}
	return m.Artifacts()
}

func (m *Module) Name() string {
	return "go-fga"
}

func (m *Module) generate(f pgs.File) {
	m.Push(f.Name().String())
	defer m.Pop()

	for _, v := range f.Services() {
		if err := m.generateModule(f, v); err != nil {
			m.Fail(err)
		}
	}
	if err := m.generateRegister(f); err != nil {
		m.Fail(err)
	}
}

func (m *Module) generateModule(f pgs.File, s pgs.Service) error {
	var mod openfgav1alpha1.Module
	ok, err := s.Extension(openfgav1alpha1.ExtModule, &mod)
	if err != nil {
		return fmt.Errorf("unable to read module extension: %w", err)
	}
	if !ok {
		return nil
	}
	if err := mod.ValidateAll(); err != nil {
		m.Failf("fga module is invalid: %v", err)
	}
	name := f.InputPath().SetExt(".fga")
	m.AddGeneratorTemplateFile(name.String(), fgaTemplate, mod)
	return nil
}

func (m *Module) generateRegister(f pgs.File) error {
	var found bool
	for _, v := range f.Services() {
		m.Push(v.Name().String())
		for _, vv := range v.Methods() {
			m.Push(vv.Name().String())
			var a openfgav1alpha1.Access
			if ok, _ := vv.Extension(openfgav1alpha1.ExtAccess, &a); !ok {
				m.Pop()
				continue
			}
			found = true
			if err := a.ValidateAll(); err != nil {
				return fmt.Errorf("fga access is invalid: %w", err)
			}
			m.Pop()
		}
		m.Pop()
	}
	if !found {
		return nil
	}
	name := m.ctx.OutputPath(f).SetExt(".fga.go")
	m.AddGeneratorTemplateFile(name.String(), m.tpl, f)
	return nil
}
