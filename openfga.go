package openfga

import (
	"context"
	"time"

	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	"github.com/openfga/openfga/pkg/server"

	"go.linka.cloud/go-openfga/storage"
	"go.linka.cloud/go-openfga/x"
	"go.linka.cloud/go-openfga/x/service"
)

type FGA interface {
	CreateStore(ctx context.Context, name string) (Store, error)
	DeleteStore(ctx context.Context, name string) error
	GetStore(ctx context.Context, name string) (Store, error)
	ListStores(ctx context.Context) ([]Store, error)
	Close()
}

type Store interface {
	AuthorizationModel(ctx context.Context, id string) (Model, error)
	LastAuthorizationModel(ctx context.Context) (Model, error)
	ListAuthorizationModels(ctx context.Context) ([]Model, error)
	WriteAuthorizationModel(ctx context.Context, dsl ...string) (Model, error)

	ID() string
	Name() string
	CreatedAt() time.Time
	UpdatedAt() time.Time
}

type TupleReader interface {
	Read(ctx context.Context, object, relation, user string) ([]*openfgav1.Tuple, error)
	ReadWithPaging(ctx context.Context, object, relation, user string, pageSize int32, continuationToken string) ([]*openfgav1.Tuple, string, error)
	Expand(ctx context.Context, object, relation string) (*openfgav1.UsersetTree, error)
	ListObjects(ctx context.Context, typ, relation, user string) ([]string, error)
	ListUsers(ctx context.Context, object, relation, userTyp string, contextKVs ...any) ([]string, error)
	ListRelations(ctx context.Context, object, user string, relations ...string) ([]string, error)

	Check(ctx context.Context, object, relation, user string, contextKVs ...any) (bool, error)
	CheckTuple(ctx context.Context, key *openfgav1.TupleKey, contextKVs ...any) (bool, error)
}

type TupleWriter interface {
	Write(ctx context.Context, object, relation, user string) error
	WriteWithCondition(ctx context.Context, object, relation, user string, condition string, kv ...any) error
	WriteTuples(context.Context, ...*openfgav1.TupleKey) error
	Delete(ctx context.Context, object, relation, user string) error
	DeleteTuples(context.Context, ...*openfgav1.TupleKey) error
}

type Model interface {
	TupleReader
	TupleWriter

	ID() string
	Store() Store
	Show() (string, error)
	// Reload(ctx context.Context) error

	Tx(ctx context.Context, opts ...storage.TxOption) (Tx, error)
	WithTx(tx any) Tx
}

type Tx interface {
	TupleReader
	TupleWriter
	Commit(ctx context.Context) error
	Close() error
}

type fga struct {
	FGA
	s x.OpenFGA
}

func FomClient(c service.Client) FGA {
	return &client{c: wrap(c)}
}

func New(s storage.Datastore, opts ...server.OpenFGAServiceV1Option) (FGA, error) {
	srv, err := x.New(s, opts...)
	if err != nil {
		return nil, err
	}
	return &fga{s: srv, FGA: &client{c: srv}}, nil
}

func (f *fga) Close() {
	f.s.Close()
}
