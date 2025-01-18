package x_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	protodb2 "go.linka.cloud/protodb"

	"go.linka.cloud/go-openfga/storage/protodb"
	"go.linka.cloud/go-openfga/tests"
	"go.linka.cloud/go-openfga/x"
)

func TestServer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := protodb.New(ctx, protodb2.WithInMemory(true))
	require.NoError(t, err)
	defer db.Close()
	f, err := x.New(db)
	require.NoError(t, err)

	tx, err := f.Tx(ctx)
	require.NoError(t, err)
	defer tx.Close()

	s, err := tx.CreateStore(ctx, "default")
	require.NoError(t, err)
	require.NotEmpty(t, s.ID())
	assert.Equal(t, "default", s.Name())

	m, err := s.WriteAuthorizationModel(ctx, tests.DSL)
	require.NoError(t, err)
	require.NotEmpty(t, m.ID())
	require.NotNil(t, m.Store())
	dsl2, err := m.Show()
	require.NoError(t, err)
	assert.Equal(t, tests.DSL, dsl2)
	m2, err := s.LastAuthorizationModel(ctx)
	require.NoError(t, err)
	assert.Equal(t, m.ID(), m2.ID())
	require.NoError(t, m.Write(ctx, tests.Doc.Ref("doc1"), tests.DocRelations.Owner, tests.User.Ref("user1")))
	ok, err := m.Check(ctx, tests.Doc.Ref("doc1"), tests.DocRelations.CanWrite, tests.User.Ref("user1"))
	require.NoError(t, err)
	assert.True(t, ok)
	ok, err = m.Check(ctx, tests.Doc.Ref("doc1"), tests.DocRelations.CanWrite, tests.User.Ref("user2"))
	require.NoError(t, err)
	assert.False(t, ok)
	tree, err := m.Expand(ctx, tests.Doc.Ref("doc1"), tests.DocRelations.Owner)
	require.NoError(t, err)
	require.Len(t, tree.GetRoot().GetLeaf().GetUsers().GetUsers(), 1)
	assert.Equal(t, tests.User.Ref("user1"), tree.GetRoot().GetLeaf().GetUsers().GetUsers()[0])
	objs, err := m.ListObjects(ctx, "doc", tests.DocRelations.CanWrite, tests.User.Ref("user1"))
	require.NoError(t, err)
	require.Len(t, objs, 1)
	assert.Equal(t, tests.Doc.Ref("doc1"), objs[0])
	require.NoError(t, m.WriteWithCondition(ctx, tests.Doc.Ref("doc1"), tests.DocRelations.Viewer, tests.User.Ref("user2"), tests.ExpiringGrant.Name, tests.ExpiringGrant.Time, "2024-01-01T00:00:00Z", tests.ExpiringGrant.Duration, "10m"))
	ok, err = m.Check(ctx, tests.Doc.Ref("doc1"), tests.DocRelations.CanRead, tests.User.Ref("user2"), tests.ExpiringGrant.Current, "2024-01-01T00:00:00Z")
	require.NoError(t, err)
	assert.True(t, ok)
	ok, err = m.Check(ctx, tests.Doc.Ref("doc1"), tests.DocRelations.CanRead, tests.User.Ref("user2"), tests.ExpiringGrant.Current, "2024-01-01T00:11:00Z")
	require.NoError(t, err)
	assert.False(t, ok)
	us, err := m.ListUsers(ctx, tests.Doc.Ref("doc1"), tests.DocRelations.CanRead, tests.User.Type(), tests.ExpiringGrant.Current, "2024-01-01T00:00:00Z")
	require.NoError(t, err)
	assert.Len(t, us, 2)
	require.NoError(t, m.Delete(ctx, tests.Doc.Ref("doc1"), tests.DocRelations.Viewer, tests.User.Ref("user2")))
	us, err = m.ListUsers(ctx, tests.Doc.Ref("doc1"), tests.DocRelations.CanRead, tests.User.Type(), tests.ExpiringGrant.Current, "2024-01-01T00:00:00Z")
	require.NoError(t, err)
	assert.Len(t, us, 1)

	tx.Close()
	tx, err = f.Tx(ctx)
	require.NoError(t, err)
	defer tx.Close()
	s, err = tx.GetStore(ctx, s.ID())
	require.Error(t, err)
	require.Nil(t, s)
}
