package main

import (
	"context"
	"flag"
	"net"
	"net/http"
	"os"

	server "github.com/casmelad/bootcamp-gateway/server"
	proto "github.com/casmelad/bootcamp-gateway/server/proto"
	implementations "github.com/casmelad/bootcamp-gateway/server/repository"
	"github.com/casmelad/bootcamp-gateway/users"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

var (
	// command-line options:
	// gRPC server endpoint
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:9090", "gRPC server endpoint")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	grpcListener, err := net.Listen("tcp", ":9090")
	if err != nil {
		os.Exit(-1)
	}

	grpcSrv := server.NewUserServer(users.NewUserService(implementations.NewInMemoryUserRepository()))
	baseServer := grpc.NewServer()
	proto.RegisterUsersServer(baseServer, grpcSrv)
	go baseServer.Serve(grpcListener)

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err = proto.RegisterUsersHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	if err != nil {
		return err
	}

	http.Handle("/", mux)
	fs := http.FileServer(http.Dir("./server/swagger"))
	http.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(":8080", nil)
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
