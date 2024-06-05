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
	"fmt"
	"strings"
)

type Reference interface {
	Ref(id string) string
	IDs(refs ...string) ([]string, error)
	Type() string
}

func NewReference(res string) Reference {
	return identifier(res)
}

type identifier string

func (i identifier) Ref(id string) string {
	return fmt.Sprintf("%s:%s", string(i), id)
}

func (i identifier) IDs(refs ...string) (out []string, err error) {
	for _, ref := range refs {
		if !strings.HasPrefix(ref, string(i)+":") {
			return nil, fmt.Errorf("invalid reference: %s for %s", ref, i)
		}
		out = append(out, strings.TrimPrefix(ref, string(i)+":"))
	}
	return
}

func (i identifier) Type() string {
	return string(i)
}

func NewReferenceWithRelation(res, rel string) Reference {
	return identifierWithRelation{res: res, rel: rel}
}

type identifierWithRelation struct {
	res string
	rel string
}

func (i identifierWithRelation) Ref(id string) string {
	return fmt.Sprintf("%s:%s#%s", i.res, id, i.rel)
}

func (i identifierWithRelation) IDs(refs ...string) (out []string, err error) {
	for _, ref := range refs {
		if !strings.HasPrefix(ref, i.res+":") {
			return nil, fmt.Errorf("invalid reference: %s for %s#%s", ref, i.res, i.rel)
		}
		if !strings.HasSuffix(ref, "#"+i.rel) {
			return nil, fmt.Errorf("invalid reference: %s for %s#%s", ref, i.res, i.rel)
		}
		out = append(out, strings.Split(strings.TrimPrefix(ref, i.res+":"), "#")[0])
	}
	return
}

func (i identifierWithRelation) Type() string {
	return i.res
}
