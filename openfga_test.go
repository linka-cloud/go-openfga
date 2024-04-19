package openfga

import (
	"context"
	"testing"

	"github.com/openfga/openfga/pkg/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	protodb2 "go.linka.cloud/protodb"

	"go.linka.cloud/go-openfga/storage/protodb"
)

const dsl = `model
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
	Goup         = NewReference("group")
	User         = NewReference("user")
	GroupMembers = NewReferenceWithRelation("group", "member")
	Doc          = NewReference("doc")
	Folder       = NewReference("folder")

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

func TestServer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := protodb.New(ctx, protodb2.WithInMemory(true))
	require.NoError(t, err)
	f, err := New(server.WithDatastore(db))
	require.NoError(t, err)
	defer f.Close()

	_, err = f.GetStore(ctx, "default")
	require.Error(t, err)

	s, err := f.CreateStore(ctx, "default")
	require.NoError(t, err)
	require.NotEmpty(t, s.ID())
	assert.Equal(t, "default", s.Name())
	_, err = s.AuthorizationModel(ctx, "")
	require.Error(t, err)
	m, err := s.WriteAuthorizationModel(ctx, dsl)
	require.NoError(t, err)
	require.NotEmpty(t, m.ID())
	require.NotNil(t, m.Store())
	dsl2, err := m.Show()
	require.NoError(t, err)
	assert.Equal(t, dsl, dsl2)
	m2, err := s.LastAuthorizationModel(ctx)
	require.NoError(t, err)
	assert.Equal(t, m.ID(), m2.ID())
	require.NoError(t, m.Write(ctx, Doc.Ref("doc1"), DocRelations.Owner, User.Ref("user1")))
	ok, err := m.Check(ctx, Doc.Ref("doc1"), DocRelations.CanWrite, User.Ref("user1"))
	require.NoError(t, err)
	assert.True(t, ok)
	ok, err = m.Check(ctx, Doc.Ref("doc1"), DocRelations.CanWrite, User.Ref("user2"))
	require.NoError(t, err)
	assert.False(t, ok)
	tree, err := m.Expand(ctx, Doc.Ref("doc1"), DocRelations.Owner)
	require.NoError(t, err)
	require.Len(t, tree.GetRoot().GetLeaf().GetUsers().GetUsers(), 1)
	assert.Equal(t, User.Ref("user1"), tree.GetRoot().GetLeaf().GetUsers().GetUsers()[0])
	objs, err := m.ListObjects(ctx, "doc", DocRelations.CanWrite, User.Ref("user1"))
	require.NoError(t, err)
	require.Len(t, objs, 1)
	assert.Equal(t, Doc.Ref("doc1"), objs[0])
	require.NoError(t, m.WriteWithCondition(ctx, Doc.Ref("doc1"), DocRelations.Viewer, User.Ref("user2"), ExpiringGrant.Name, ExpiringGrant.Time, "2024-01-01T00:00:00Z", ExpiringGrant.Duration, "10m"))
	ok, err = m.Check(ctx, Doc.Ref("doc1"), DocRelations.CanRead, User.Ref("user2"), ExpiringGrant.Current, "2024-01-01T00:00:00Z")
	require.NoError(t, err)
	assert.True(t, ok)
	ok, err = m.Check(ctx, Doc.Ref("doc1"), DocRelations.CanRead, User.Ref("user2"), ExpiringGrant.Current, "2024-01-01T00:11:00Z")
	require.NoError(t, err)
	assert.False(t, ok)
	us, err := m.ListUsers(ctx, Doc.Ref("doc1"), DocRelations.CanRead, User.Type(), ExpiringGrant.Current, "2024-01-01T00:00:00Z")
	require.NoError(t, err)
	assert.Len(t, us, 2)
	require.NoError(t, m.Delete(ctx, Doc.Ref("doc1"), DocRelations.Viewer, User.Ref("user2")))
	us, err = m.ListUsers(ctx, Doc.Ref("doc1"), DocRelations.CanRead, User.Type(), ExpiringGrant.Current, "2024-01-01T00:00:00Z")
	require.NoError(t, err)
	assert.Len(t, us, 1)
}
