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

	"go.linka.cloud/go-openfga/openfga"
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
		"module": m.module,
		"types": func(s pgs.Service) []Type {
			m.Push(s.Name().String())
			defer m.Pop()
			t := make(map[string]Type)
			mod := m.module(s)
			if mod != nil {
				for _, v := range append(mod.Extends, mod.Definitions...) {
					t[v.Type] = Type{Name: v.Type, Relations: make(map[string]struct{})}
					for _, vv := range v.Relations {
						t[v.Type].Relations[vv.Define] = struct{}{}
					}
				}
			}
			for _, v := range s.Methods() {
				m.Push(v.Name().String())
				a := m.access(v)
				if a == nil {
					m.Pop()
					continue
				}
				for _, vv := range a.Check {
					if vv.GetAs() != "" {
						continue
					}
					if _, ok := t[vv.GetType()]; !ok {
						t[vv.GetType()] = Type{Name: vv.GetType(), Relations: make(map[string]struct{})}
					}
					if _, ok := t[vv.GetType()].Relations[vv.GetCheck()]; !ok {
						t[vv.GetType()].Relations[vv.GetCheck()] = struct{}{}
					}
				}
				m.Pop()
			}
			var out []Type
			for _, v := range t {
				out = append(out, v)
			}
			return out
		},
		"access": m.access,
		"need_getter": func(s *openfga.Step) bool {
			id := s.GetID()
			return len(id) >= 2 && id[0] == '{' && id[len(id)-1] == '}'
		},
		"field": func(s *openfga.Step) string {
			id := s.GetID()
			if len(id) < 2 || id[0] != '{' || id[len(id)-1] != '}' {
				return ""
			}
			return id[1 : len(id)-1]
		},
		"getter": func(s *openfga.Step, me pgs.Method) string {
			t := me.Input()
			id := s.GetID()
			parts := strings.Split(id[1:len(id)-1], ".")
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
		"object": func(s *openfga.Step) string {
			return s.GetType() + ":" + s.GetID()
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

func (m *Module) module(s pgs.Service) *openfga.Module {
	m.Push(s.Name().String())
	defer m.Pop()
	mod := m.moduleExt(s)
	for _, v := range s.Methods() {
		m.Push(v.Name().String())
		a := m.accessExt(v)
		for _, vv := range a.Check {
			if vv == nil || vv.GetAs() == "" {
				continue
			}
			if mod == nil {
				m.Fail("\"as\" is not allowed without a module definition")
			}
			a = m.defaults(a, v)
			t := find(mod, vv.GetType())
			if t == nil {
				m.Failf("type %q not found in module", vv.GetAs())
			}
			t.Relations = append(t.Relations, &openfga.Relation{
				Define: define(mod, v, vv.GetType()),
				As:     vv.GetAs(),
			})
		}
		m.Pop()
	}
	return mod
}

func (m *Module) moduleExt(s pgs.Service) *openfga.Module {
	var mod openfga.Module
	ok, err := s.Extension(openfga.ExtModule, &mod)
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
}

func (m *Module) access(me pgs.Method) *openfga.Access {
	a := m.accessExt(me)
	if a == nil {
		return nil
	}
	mod := m.moduleExt(me.Service())
	a = m.defaults(a, me)
	for _, v := range a.Check {
		if v.GetAs() == "" {
			continue
		}
		if mod == nil {
			m.Fail("\"as\" is not allowed without a module definition")
		}
		v.Relation = &openfga.Step_Check{Check: define(mod, me, v.GetType())}
	}
	return a
}

func (m *Module) accessExt(me pgs.Method) *openfga.Access {
	var access openfga.Access
	ok, err := me.Extension(openfga.ExtAccess, &access)
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
}

func (m *Module) defaults(a *openfga.Access, me pgs.Method) *openfga.Access {
	var id, typ *string
	var ignore *bool
	if d := m.defaultAccess(me.Service()); d != nil {
		id, typ, ignore = d.ID, d.Type, d.IgnoreNotFound
	}
	for _, v := range a.Check {
		if v.ID == nil {
			v.ID = id
		}
		if v.Type == nil {
			v.Type = typ
		}
		if v.IgnoreNotFound == nil {
			v.IgnoreNotFound = ignore
		}
		if v.ID == nil || v.Type == nil {
			m.Failf("access %q is missing id or type", me.Name())
		}
	}
	return a
}

func (m *Module) defaultAccess(s pgs.Service) *openfga.DefaultAccess {
	var access openfga.DefaultAccess
	ok, err := s.Extension(openfga.ExtDefaults, &access)
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
	mod := m.module(s)
	if mod == nil {
		return nil
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
			var a openfga.Access
			if ok, _ := vv.Extension(openfga.ExtAccess, &a); !ok {
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

func define(mod *openfga.Module, m pgs.Method, typ string, prefix ...string) string {
	p := "can_"
	if len(prefix) > 0 {
		p = prefix[0] + "_"
	}
	method := m.Name().LowerSnakeCase().String()
	var name string
	if strings.Contains(method, "_"+typ) || strings.Contains(method, typ+"_") {
		name = p + method
	} else {
		name = fmt.Sprintf("%s%s_%s", p, mod.Name, method)
	}
	if strings.Contains(name, "_"+typ) {
		return strings.Replace(name, "_"+typ, "", 1)
	}
	return name
}

func find(mod *openfga.Module, name string) *openfga.Type {
	for _, v := range mod.Definitions {
		if v.Type == name {
			return v
		}
	}
	for _, v := range mod.Extends {
		if v.Type == name {
			return v
		}
	}
	return nil
}

type Type struct {
	Name      string
	Relations map[string]struct{}
}
