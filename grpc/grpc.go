package grpc

import (
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/chapar-rest/mock-server/api/petstorev1"
)

type Grpc struct {
	petstorev1.UnimplementedPetstoreServiceServer
}

func New() *Grpc {
	g := &Grpc{}
	return g
}

func (g *Grpc) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8090))
	if err != nil {
		slog.Error("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	petstorev1.RegisterPetstoreServiceServer(grpcServer, g)
	reflection.Register(grpcServer)
	slog.Info("starting grpc server", "port", 8090)
	if err := grpcServer.Serve(lis); err != nil {
		slog.Error("failed to start server", "error", err)
		return err
	}

	return nil
}
