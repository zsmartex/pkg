package grpc_fx

import (
	"context"

	"github.com/zsmartex/pkg/v2/config"
	"github.com/zsmartex/pkg/v2/log"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

var (
	ServerModule = fx.Module("grpcfx_server", grpcServerProviders, grpcServerInvokes)

	GrpcClientOptions = fx.Options(
		grpcClientInvokes,
	)

	grpcServerProviders = fx.Options(fx.Provide(
		NewGrpcServer,
	))

	grpcServerInvokes = fx.Invoke(registerServerHooks)
	grpcClientInvokes = fx.Invoke(RegisterClientHooks)
)

func registerServerHooks(
	lc fx.Lifecycle,
	grpcServer GrpcServer,
	config config.GRPC,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := grpcServer.RunGrpcServer(nil); err != nil {
					// do a fatal for going to OnStop process
					log.Fatalf("gRPC error in running server: %v", err)
				}
			}()
			log.Infof("Grpc is listening on %s", config.Address())

			return nil
		},
		OnStop: func(ctx context.Context) error {
			grpcServer.GracefulShutdown()
			log.Info("server shutdown gracefully")

			return nil
		},
	})
}

func RegisterClientHooks(
	lc fx.Lifecycle,
	conn *grpc.ClientConn,
) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			if err := conn.Close(); err != nil {
				log.Errorf("error in closing grpc-client: %v", err)
			} else {
				log.Info("grpc-client closed gracefully")
			}

			return nil
		},
	})
}
