package openfga

import (
	"context"
	"time"

	"github.com/fullstorydev/grpchan/inprocgrpc"
	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	"github.com/openfga/openfga/pkg/server"
)

type FGA interface {
	Client
	Service() *server.Server
	Close()
}

type Client interface {
	CreateStore(ctx context.Context, name string) (Store, error)
	DeleteStore(ctx context.Context, name string) error
	GetStore(ctx context.Context, name string) (Store, error)
	ListStores(ctx context.Context) ([]Store, error)
}

type Store interface {
	AuthorizationModel(ctx context.Context, id string) (Model, error)
	LastAuthorizationModel(ctx context.Context) (Model, error)
	ListAuthorizationModels(ctx context.Context) ([]Model, error)
	WriteAuthorizationModel(ctx context.Context, dsl string) (Model, error)

	ID() string
	Name() string
	CreatedAt() time.Time
	UpdatedAt() time.Time
}

type Model interface {
	ID() string
	Store() Store

	Write(ctx context.Context, object, relation, user string) error
	WriteTuples(context.Context, ...*openfgav1.TupleKey) error
	Delete(ctx context.Context, object, relation, user string) error
	DeleteTuples(context.Context, ...*openfgav1.TupleKey) error
	Tx() Tx
	Check(ctx context.Context, object, relation, user string) (bool, error)
	CheckTuple(ctx context.Context, key *openfgav1.TupleKey) (bool, error)
}

type Tx interface {
	Write(object, relation, user string) error
	WriteTuples(...*openfgav1.TupleKey) error
	Delete(object, relation, user string) error
	DeleteTuples(...*openfgav1.TupleKey) error
	Commit(ctx context.Context) error
	Close()
}

type fga struct {
	Client
	s *server.Server
}

func New(opts ...server.OpenFGAServiceV1Option) (FGA, error) {
	s, err := server.NewServerWithOpts(opts...)
	if err != nil {
		return nil, err
	}
	ch := &inprocgrpc.Channel{}
	openfgav1.RegisterOpenFGAServiceServer(ch, s)
	return &fga{s: s, Client: &client{c: openfgav1.NewOpenFGAServiceClient(ch)}}, nil
}

func (f *fga) Service() *server.Server {
	return f.s
}

func (f *fga) Close() {
	f.s.Close()
}
