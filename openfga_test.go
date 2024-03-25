package openfga

import (
	"context"
	"testing"

	"github.com/openfga/openfga/pkg/server"
	"github.com/openfga/openfga/pkg/storage/memory"
	"github.com/stretchr/testify/require"
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
    define viewer: [user, user:*, group#member] or owner or viewer from parent

type doc
  relations
    define can_change_owner: owner
    define owner: [user, group]
    define parent: [folder]
    define can_read: viewer or owner or viewer from parent
    define can_share: owner or owner from parent
    define viewer: [user, user:*, group#member]
    define can_write: owner or owner from parent
`

var (
	Goup         = NewReference("group")
	User         = NewReference("user")
	GroupMembers = NewReferenceWithRelation("group", "member")
	Doc          = NewReference("doc")
	Folder       = NewReference("folder")
)

func TestServer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mem := memory.New()
	f, err := New(server.WithDatastore(mem))
	require.NoError(t, err)
	defer f.Close()

	_, err = f.GetStore(ctx, "default")
	require.Error(t, err)

	s, err := f.CreateStore(ctx, "default")
	require.NoError(t, err)
	require.NotEmpty(t, s.ID())
	require.Equal(t, "default", s.Name())
	_, err = s.AuthorizationModel(ctx, "")
	require.Error(t, err)
	m, err := s.WriteAuthorizationModel(ctx, dsl)
	require.NoError(t, err)
	require.NotEmpty(t, m.ID())
	require.NotNil(t, m.Store())
	m2, err := s.LastAuthorizationModel(ctx)
	require.NoError(t, err)
	require.Equal(t, m.ID(), m2.ID())
	require.NoError(t, m.Write(ctx, Doc.Ref("doc1"), "owner", User.Ref("user1")))
	ok, err := m.Check(ctx, Doc.Ref("doc1"), "can_write", User.Ref("user1"))
	require.NoError(t, err)
	require.True(t, ok)
	ok, err = m.Check(ctx, Doc.Ref("doc1"), "can_write", User.Ref("user2"))
	require.NoError(t, err)
	require.False(t, ok)
}
