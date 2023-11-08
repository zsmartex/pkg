package grpc_fx

import (
	"context"

	"github.com/zsmartex/pkg/v2/infrastructure/grpc_fx"
	"github.com/zsmartex/pkg/v2/log"
	"go.uber.org/fx"
	"google.golang.org/grpc/test/bufconn"
)

var ServerModule = fx.Module("grpc_fx_testing",
	fx.Provide(NewListener),
	fx.Decorate(NewGrpcServer),
	fx.Invoke(registerServerHooks),
)

func NewListener() *bufconn.Listener {
	return bufconn.Listen(1024 * 1024)
}

func registerServerHooks(
	lc fx.Lifecycle,
	grpcServer grpc_fx.GrpcServer,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := grpcServer.RunGrpcServer(nil); err != nil {
					// do a fatal for going to OnStop process
					log.Fatalf("gRPC error in running server: %v", err)
				}
			}()
			log.Info("Grpc test is listening now")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			grpcServer.GracefulShutdown()
			log.Info("server shutdown gracefully")

			return nil
		},
	})
}
