package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"time"

	"github.com/fullstorydev/grpchan/inprocgrpc"
	protodb2 "go.linka.cloud/protodb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"go.linka.cloud/go-openfga"
	pb "go.linka.cloud/go-openfga/example/pb"
	"go.linka.cloud/go-openfga/interceptors"
	"go.linka.cloud/go-openfga/storage/protodb"
)

//go:embed base.fga
var modelBase string

// userKey is the key used to store the user in the context metadata
const userKey = "user"

// defaultSystem is the default system object
var defaultSystem = pb.FGASystemObject("default")

// userContext returns a new context with the user set in the metadata
func userContext(ctx context.Context, user string) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.Pairs(userKey, user))
}

// contextUser returns the user from the context metadata
func contextUser(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md.Get(userKey)) == 0 {
		return "", status.Errorf(codes.Unauthenticated, "missing user from metadata")
	}
	return "user:" + md.Get(userKey)[0], nil
}

// mustContextUser returns the user from the context metadata or panics
func mustContextUser(ctx context.Context) string {
	user, err := contextUser(ctx)
	if err != nil {
		panic(err)
	}
	return user
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create the in-memory openfga server
	ds, err := protodb.New(ctx, true, protodb2.WithInMemory(true))
	f, err := openfga.New(ds)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// create the store
	s, err := f.CreateStore(ctx, "default")
	if err != nil {
		log.Fatal(err)
	}

	// write the model
	model, err := s.WriteAuthorizationModel(ctx, modelBase, pb.FGAModel)
	if err != nil {
		log.Fatal(err)
	}

	// create the interceptors
	fga, err := interceptors.New(ctx, model, interceptors.WithUserFunc(func(ctx context.Context) (string, map[string]any, error) {
		user, err := contextUser(ctx)
		if err != nil {
			return "", nil, err
		}
		return user, nil, nil
	}))
	if err != nil {
		log.Fatal(err)
	}

	// register some users with system roles
	for _, v := range []string{
		pb.FGASystemResourceReader,
		pb.FGASystemResourceWriter,
		pb.FGASystemResourceAdmin,
		pb.FGASystemResourceWatcher,
	} {
		if err := model.Write(ctx, defaultSystem, v, fmt.Sprintf("user:%s", v)); err != nil {
			log.Fatal(err)
		}
	}

	// create the service
	svc := NewResourceService()

	// register the service permissions
	pb.RegisterFGA(fga)

	// create the in-process grpc channel
	channel := (&inprocgrpc.Channel{}).
		WithServerUnaryInterceptor(fga.UnaryServerInterceptor()).
		WithServerStreamInterceptor(fga.StreamServerInterceptor())

	// register the service as usual
	pb.RegisterResourceServiceServer(channel, svc)

	// create a client
	client := pb.NewResourceServiceClient(channel)

	// validate checks
	if _, err := client.List(userContext(ctx, pb.FGASystemResourceReader), &pb.ListRequest{}); err != nil {
		log.Fatal(err)
	}

	if _, err := client.Create(userContext(ctx, pb.FGASystemResourceReader), &pb.CreateRequest{Resource: &pb.Resource{ID: "0"}}); err == nil {
		log.Fatal("reader should not be able to create")
	}

	if _, err := client.Create(userContext(ctx, pb.FGASystemResourceWriter), &pb.CreateRequest{Resource: &pb.Resource{ID: "0"}}); err != nil {
		log.Fatal(err)
	}

	wctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	ss, err := client.Watch(userContext(wctx, pb.FGASystemResourceWriter), &pb.WatchRequest{})
	if err != nil {
		log.Fatal(err)
	}
	// try to receive an event as the interceptor is not called when creating the stream
	if _, err := ss.Recv(); err == nil {
		log.Fatal("writer should not be able to watch")
	}

	wctx, cancel = context.WithTimeout(ctx, time.Second)
	defer cancel()
	ss, err = client.Watch(userContext(wctx, pb.FGASystemResourceAdmin), &pb.WatchRequest{})
	if err != nil {
		log.Fatal(err)
	}

	// create a resource to trigger an event
	go func() {
		time.Sleep(100 * time.Millisecond)
		if _, err := client.Create(userContext(ctx, pb.FGASystemResourceWriter), &pb.CreateRequest{Resource: &pb.Resource{ID: "1"}}); err != nil {
			log.Fatal(err)
		}
	}()
	if _, err := ss.Recv(); err != nil {
		log.Fatal(err)
	}
}
