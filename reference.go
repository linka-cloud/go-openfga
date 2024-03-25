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
)

type Reference interface {
	Ref(id string) string
}

func NewReference(res string) Reference {
	return identifier(res)
}

type identifier string

func (i identifier) Ref(id string) string {
	return fmt.Sprintf("%s:%s", string(i), id)
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