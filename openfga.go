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

type FGA[T any] interface {
	CreateStore(ctx context.Context, name string) (Store[T], error)
	DeleteStore(ctx context.Context, name string) error
	GetStore(ctx context.Context, name string) (Store[T], error)
	ListStores(ctx context.Context) ([]Store[T], error)
	Close()
}

type Store[T any] interface {
	AuthorizationModel(ctx context.Context, id string) (Model[T], error)
	LastAuthorizationModel(ctx context.Context) (Model[T], error)
	ListAuthorizationModels(ctx context.Context) ([]Model[T], error)
	WriteAuthorizationModel(ctx context.Context, dsl ...string) (Model[T], error)

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

type Model[T any] interface {
	TupleReader
	TupleWriter

	ID() string
	Store() Store[T]
	Show() (string, error)
	// Reload(ctx context.Context) error

	Tx(ctx context.Context, opts ...storage.TxOption) (Tx, error)
	WithTx(tx T) Tx
}

type Tx interface {
	TupleReader
	TupleWriter
	Commit(ctx context.Context) error
	Close() error
}

type fga[T any] struct {
	FGA[T]
	s x.OpenFGA[T]
}

func FomClient(c service.Client) FGA[none] {
	return &client[none]{c: wrap(c)}
}

func New[T any](s storage.Datastore[T], opts ...server.OpenFGAServiceV1Option) (FGA[T], error) {
	srv, err := x.New(s, opts...)
	if err != nil {
		return nil, err
	}
	return &fga[T]{s: srv, FGA: &client[T]{c: srv}}, nil
}

func (f *fga[T]) Close() {
	f.s.Close()
}
