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

package tests

import (
	"go.linka.cloud/go-openfga"
)

const DSL = `model
  schema 1.1

type user

type group
  relations
    define member: [user, group#member]

type folder
  relations
    define can_create_file: owner
    define owner: [user]
    define parent: [folder]
    define viewer: [user with non_expired_grant, group#member] or owner or viewer from parent

type doc
  relations
    define can_change_owner: owner
    define can_read: viewer or owner or viewer from parent
    define can_share: owner or owner from parent
    define can_write: owner or owner from parent
    define owner: [user, group]
    define parent: [folder]
    define viewer: [user with non_expired_grant, group#member]

condition non_expired_grant(current_time: timestamp, grant_duration: duration, grant_time: timestamp) {
  current_time < grant_time + grant_duration
}
`

var (
	Group        = openfga.NewReference("group")
	User         = openfga.NewReference("user")
	GroupMembers = openfga.NewReferenceWithRelation("group", "member")
	Doc          = openfga.NewReference("doc")
	Folder       = openfga.NewReference("folder")

	GroupRelations = struct {
		Member string
	}{
		Member: "member",
	}

	FolderRelations = struct {
		CanCreateFile string
		Owner         string
		Parent        string
		Viewer        string
	}{
		CanCreateFile: "can_create_file",
		Owner:         "owner",
		Parent:        "parent",
		Viewer:        "viewer",
	}

	DocRelations = struct {
		CanChangeOwner string
		Owner          string
		Parent         string
		CanRead        string
		CanShare       string
		Viewer         string
		CanWrite       string
	}{
		CanChangeOwner: "can_change_owner",
		Owner:          "owner",
		Parent:         "parent",
		CanRead:        "can_read",
		CanShare:       "can_share",
		Viewer:         "viewer",
		CanWrite:       "can_write",
	}

	ExpiringGrant = struct {
		Name     string
		Current  string
		Time     string
		Duration string
	}{
		Name:     "non_expired_grant",
		Current:  "current_time",
		Time:     "grant_time",
		Duration: "grant_duration",
	}
)
