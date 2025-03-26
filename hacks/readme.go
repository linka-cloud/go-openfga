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
	"os"
	"text/template"
)

//go:embed README.gotpl.md
var readme string

var funcs = template.FuncMap{
	"file": func(name string) string {
		b, err := os.ReadFile(name)
		if err != nil {
			panic(err)
		}
		return string(b)
	},
}

func main() {
	if err := template.Must(template.New("readme").Funcs(funcs).Parse(readme)).Execute(os.Stdout, nil); err != nil {
		panic(err)
	}
}
