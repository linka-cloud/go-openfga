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

	"go.linka.cloud/go-openfga/x"
)

type tx struct {
	c x.Tx
	*rw
}

func (t *tx) OnMissingIgnore() Tx {
	t2 := t.clone()
	t2.rw.onMissing = onMissingIgnore
	return t2
}

func (t *tx) OnMissingError() Tx {
	t2 := t.clone()
	t2.rw.onMissing = onMissingError
	return t2
}

func (t *tx) OnDuplicateIgnore() Tx {
	t2 := t.clone()
	t2.rw.onDuplicate = onDuplicateIgnore
	return t2
}

func (t *tx) OnDuplicateError() Tx {
	t2 := t.clone()
	t2.rw.onDuplicate = onDuplicateError
	return t2
}

func (t *tx) Commit(ctx context.Context) error {
	return t.c.Commit(ctx)
}

func (t *tx) Close() error {
	return t.c.Close()
}

func (t *tx) clone() *tx {
	return &tx{
		c:  t.c,
		rw: t.rw.clone(),
	}
}
