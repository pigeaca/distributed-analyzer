package app

import (
	"context"
	"google.golang.org/grpc"
	"net"
)

type GrpcComponent struct {
	server   *grpc.Server
	listener net.Listener
}

func NewGrpcComponent(server *grpc.Server, listener net.Listener) *GrpcComponent {
	return &GrpcComponent{server: server, listener: listener}
}

func (g *GrpcComponent) Start(ctx context.Context) error {
	go func() {
		err := g.server.Serve(g.listener)
		if err != nil {
			panic(err)
		}
	}()
	return nil
}

func (g *GrpcComponent) Stop(ctx context.Context) error {
	g.server.GracefulStop()
	return nil
}

func (g *GrpcComponent) Name() string {
	return "gRPC"
}
