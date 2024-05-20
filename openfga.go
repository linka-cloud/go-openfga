package openfga

import (
	"context"
	"time"

	"github.com/fullstorydev/grpchan/inprocgrpc"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	"github.com/openfga/openfga/pkg/middleware/requestid"
	"github.com/openfga/openfga/pkg/middleware/storeid"
	"github.com/openfga/openfga/pkg/middleware/validator"
	"github.com/openfga/openfga/pkg/server"
	"google.golang.org/grpc"
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
	Show() (string, error)

	Read(ctx context.Context, object, relation, user string) ([]*openfgav1.Tuple, error)
	Expand(ctx context.Context, object, relation string) (*openfgav1.UsersetTree, error)
	List(ctx context.Context, typ, relation, user string) ([]string, error)
	Tx() Tx
	Check(ctx context.Context, object, relation, user string) (bool, error)
	CheckWithContext(ctx context.Context, object, relation, user string, kv ...any) (bool, error)
	CheckTuple(ctx context.Context, key *openfgav1.TupleKey) (bool, error)
	CheckTupleWithContext(ctx context.Context, key *openfgav1.TupleKey, kv ...any) (bool, error)
	Write(ctx context.Context, object, relation, user string) error
	WriteWithCondition(ctx context.Context, object, relation, user string, condition string, kv ...any) error
	WriteTuples(context.Context, ...*openfgav1.TupleKey) error
	Delete(ctx context.Context, object, relation, user string) error
	DeleteTuples(context.Context, ...*openfgav1.TupleKey) error
}

type Tx interface {
	Write(object, relation, user string) error
	WriteTuples(...*openfgav1.TupleKey) error
	WriteWithCondition(object, relation, user string, condition string, kv ...any) error
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
	ch.WithServerUnaryInterceptor(
		grpcmiddleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(), // needed for logging
			requestid.NewUnaryInterceptor(),       // add request_id to ctxtags
			storeid.NewUnaryInterceptor(),         // if available, add store_id to ctxtags
			// logging.NewLoggingInterceptor(s.Logger), // needed to log invalid requests
			validator.UnaryServerInterceptor(),
		),
	)
	ch.WithServerStreamInterceptor(
		grpcmiddleware.ChainStreamServer(
			[]grpc.StreamServerInterceptor{
				requestid.NewStreamingInterceptor(),
				validator.StreamServerInterceptor(),
				grpc_ctxtags.StreamServerInterceptor(),
			}...,
		),
	)
	openfgav1.RegisterOpenFGAServiceServer(ch, s)
	return &fga{s: s, Client: &client{c: openfgav1.NewOpenFGAServiceClient(ch)}}, nil
}

func (f *fga) Service() *server.Server {
	return f.s
}

func (f *fga) Close() {
	f.s.Close()
}
