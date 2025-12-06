package openfga_test

import (
	"context"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.linka.cloud/grpc-toolkit/logger"
	protodb2 "go.linka.cloud/protodb"

	"go.linka.cloud/go-openfga"
	"go.linka.cloud/go-openfga/storage/protodb"
	"go.linka.cloud/go-openfga/tests"
)

func TestServer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := protodb.New(ctx, true, protodb2.WithInMemory(true))
	require.NoError(t, err)
	defer db.Close()
	f, err := openfga.New(db)
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
	rs, err := m.ListRelations(ctx, tests.Doc.Ref("doc1"), tests.User.Ref("user1"))
	require.NoError(t, err)
	slices.Sort(rs)
	assert.Equal(t, []string{"can_change_owner", "can_read", "can_share", "can_write", "owner"}, rs)
	assert.Len(t, us, 1)

	tx, err := m.Tx(ctx)
	require.NoError(t, err)
	require.NoError(t, tx.Write(ctx, tests.Doc.Ref("doc1"), tests.DocRelations.Owner, tests.User.Ref("user3")))
	ok, err = tx.Check(ctx, tests.Doc.Ref("doc1"), tests.DocRelations.Owner, tests.User.Ref("user3"))
	require.NoError(t, err)
	assert.True(t, ok)
	require.NoError(t, tx.Close())

	ok, err = m.Check(ctx, tests.Doc.Ref("doc1"), tests.DocRelations.Owner, tests.User.Ref("user3"))
	require.NoError(t, err)
	assert.False(t, ok)

	require.NoError(t, m.OnMissingIgnore().Delete(ctx, tests.Doc.Ref("noop"), tests.DocRelations.Viewer, tests.User.Ref("noop")))
	require.Error(t, m.OnMissingError().Delete(ctx, tests.Doc.Ref("noop"), tests.DocRelations.Viewer, tests.User.Ref("noop")))
}

func TestWithTx(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := protodb2.Open(ctx, protodb2.WithInMemory(true), protodb2.WithLogger(logger.C(ctx).WithOffset(1)))
	require.NoError(t, err)
	defer db.Close()
	fdb, err := protodb.NewWithClient(ctx, db, true)
	require.NoError(t, err)
	defer db.Close()
	f, err := openfga.New(fdb)
	require.NoError(t, err)
	defer f.Close()

	s, err := f.CreateStore(ctx, "default")
	require.NoError(t, err)
	m, err := s.WriteAuthorizationModel(ctx, tests.DSL)
	require.NoError(t, err)
	txn, err := db.Tx(ctx)
	require.NoError(t, err)
	defer txn.Close()
	tx := m.WithTx(txn)
	require.NoError(t, tx.Write(ctx, tests.Doc.Ref("doc1"), tests.DocRelations.Owner, tests.User.Ref("user1")))
	require.NoError(t, tx.OnMissingIgnore().Delete(ctx, tests.Doc.Ref("noop"), tests.DocRelations.Viewer, tests.User.Ref("noop")))
	require.Error(t, tx.OnMissingError().Delete(ctx, tests.Doc.Ref("noop"), tests.DocRelations.Viewer, tests.User.Ref("noop")))
	require.NoError(t, txn.Commit(ctx))
	require.Error(t, tx.Commit(ctx))
	ts, err := m.Read(ctx, "", "", "")
	require.NoError(t, err)
	require.Len(t, ts, 1)

}
