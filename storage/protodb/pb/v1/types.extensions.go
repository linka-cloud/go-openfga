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

package pbv1

import (
	"fmt"

	"github.com/oklog/ulid/v2"
	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TuplePrefix(storeID string) string {
	return fmt.Sprintf("%s/", storeID)
}

func TupleKey(storeID string, x *openfgav1.TupleKey) string {
	return fmt.Sprintf("%s/%s#%s@%s", storeID, x.Object, x.Relation, x.User)
}

func (x *Tuple) SetKey() *Tuple {
	x.Key = fmt.Sprintf("%s/%s#%s@%s", x.StoreId, x.TupleKey.Object, x.TupleKey.Relation, x.TupleKey.User)
	return x
}

func AssertionsKey(storeID, modelID string) string {
	return fmt.Sprintf("%s/%s", storeID, modelID)
}

func (x *Assertions) SetKey() *Assertions {
	x.Key = fmt.Sprintf("%s/%s", x.StoreId, x.ModelId)
	return x
}

func ModelPrefix(storeID string) string {
	return fmt.Sprintf("%s/", storeID)
}

func (x *Model) SetKey() *Model {
	x.Key = fmt.Sprintf("%s/%s", x.StoreId, x.Model.Id)
	return x
}

func NewWriteChange(storeID string, t *Tuple) *Change {
	return &Change{
		Key:     fmt.Sprintf("%s/%s", storeID, ulid.MustNew(ulid.Timestamp(t.CreatedAt.AsTime()), ulid.DefaultEntropy()).String()),
		StoreId: storeID,
		Change: &openfgav1.TupleChange{
			TupleKey:  t.TupleKey,
			Operation: openfgav1.TupleOperation_TUPLE_OPERATION_WRITE,
			Timestamp: t.CreatedAt,
		},
	}
}

func NewDeleteChange(storeID string, t *Tuple) *Change {
	ts := timestamppb.Now()
	return &Change{
		Key:     fmt.Sprintf("%s/%s", storeID, ulid.MustNew(ulid.Timestamp(ts.AsTime()), ulid.DefaultEntropy()).String()),
		StoreId: storeID,
		Change: &openfgav1.TupleChange{
			TupleKey:  t.TupleKey,
			Operation: openfgav1.TupleOperation_TUPLE_OPERATION_DELETE,
			Timestamp: timestamppb.Now(),
		},
	}
}
