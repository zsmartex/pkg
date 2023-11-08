package grpc_fx

import (
	"fmt"
	"net"
	"time"

	"github.com/cockroachdb/errors"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcCtxTags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/zsmartex/pkg/v2/config"
	"github.com/zsmartex/pkg/v2/log"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
)

const (
	maxConnectionIdle = 5
	gRPCTimeout       = 15
	maxConnectionAge  = 5
	gRPCTime          = 10
)

type GrpcServer interface {
	RunGrpcServer(configGrpc ...func(grpcServer *grpc.Server)) error
	GracefulShutdown()
	GetCurrentGrpcServer() *grpc.Server
}

type grpcServer struct {
	server *grpc.Server
	config config.GRPC
}

type grpcServerParams struct {
	fx.In

	Config config.GRPC `name:"grpc_server"`
}

func NewGrpcServer(params grpcServerParams) GrpcServer {
	unaryServerInterceptors := []grpc.UnaryServerInterceptor{
		grpcCtxTags.UnaryServerInterceptor(),
		grpcRecovery.UnaryServerInterceptor(),
	}

	s := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: maxConnectionIdle * time.Minute,
			Timeout:           gRPCTimeout * time.Second,
			MaxConnectionAge:  maxConnectionAge * time.Minute,
			Time:              gRPCTime * time.Minute,
		}),
		// https://github.com/open-telemetry/opentelemetry-go-contrib/tree/00b796d0cdc204fa5d864ec690b2ee9656bb5cfc/instrumentation/google.golang.org/grpc/otelgrpc
		// github.com/grpc-ecosystem/go-grpc-middleware
		grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
			unaryServerInterceptors...,
		)),
	)

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	return &grpcServer{
		server: s,
		config: params.Config,
	}
}

func (s *grpcServer) RunGrpcServer(
	configGrpc ...func(grpcServer *grpc.Server),
) error {
	l, err := net.Listen("tcp", s.config.Address())
	if err != nil {
		return errors.Wrap(err, "net.Listen")
	}

	if len(configGrpc) > 0 {
		grpcFunc := configGrpc[0]
		if grpcFunc != nil {
			grpcFunc(s.server)
		}
	}

	log.Infof("gRPC server is listening on: %s", s.config.Address())

	err = s.server.Serve(l)
	if err != nil {
		log.Error(
			fmt.Sprintf("gRPC server serve error: %+v", err),
		)
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
