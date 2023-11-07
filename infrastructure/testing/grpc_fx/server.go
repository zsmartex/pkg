package grpc_fx

import (
	"github.com/zsmartex/pkg/v2/infrastructure/grpc_fx"
	"github.com/zsmartex/pkg/v2/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type grpcServer struct {
	server   *grpc.Server
	listener *bufconn.Listener
}

func NewGrpcServer(listener *bufconn.Listener) grpc_fx.GrpcServer {
	s := grpc.NewServer()

	return &grpcServer{
		server:   s,
		listener: listener,
	}
}

func (s *grpcServer) RunGrpcServer(
	configGrpc ...func(grpcServer *grpc.Server),
) error {
	if len(configGrpc) > 0 {
		grpcFunc := configGrpc[0]
		if grpcFunc != nil {
			grpcFunc(s.server)
		}
	}

	err := s.server.Serve(s.listener)
	if err != nil {
		log.Errorf("gRPC server serve error: %+v", err)
	}

	return err
}

func (s *grpcServer) GetCurrentGrpcServer() *grpc.Server {
	return s.server
}

func (s *grpcServer) GracefulShutdown() {
	s.server.Stop()
	s.server.GracefulStop()
}
