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

	"go.linka.cloud/go-openfga/storage"
	"go.linka.cloud/go-openfga/x"
	pbv1 "go.linka.cloud/go-openfga/x/pb/v1"
)

type FGA interface {
	Client
	Service() x.Server
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
}

type Tx interface {
	TupleReader
	TupleWriter
	Commit(ctx context.Context) error
	Close() error
}

type fga struct {
	Client
	s x.Server
}

func FomClient(c x.Client) Client {
	return &client{c: c}
}

func New[T any](s storage.Datastore[T], opts ...server.OpenFGAServiceV1Option) (FGA, error) {
	svc, err := x.Wrap(s, opts...)
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
	openfgav1.RegisterOpenFGAServiceServer(ch, svc)
	pbv1.RegisterOpenFGAXServiceServer(ch, svc)
	return &fga{s: svc, Client: &client{c: x.NewClient(ch)}}, nil
}

func (f *fga) Service() x.Server {
	return f.s
}

func (f *fga) Close() {
	f.s.Close()
}
