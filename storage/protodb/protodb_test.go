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

package protodb

import (
	"context"
	"testing"

	"github.com/openfga/openfga/pkg/storage/test"
	"github.com/stretchr/testify/require"
	"go.linka.cloud/protodb"
)

func TestStorage(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := protodb.Open(ctx, protodb.WithInMemory(true))
	require.NoError(t, err)
	defer db.Close()
	ds, err := NewWithClient(ctx, db, false)
	require.NoError(t, err)
	test.RunAllTests(t, ds)
}
