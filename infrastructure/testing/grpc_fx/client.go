package grpc_fx

import (
	"context"
	"net"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func AsClientParams(resultTags string) interface{} {
	return fx.Annotate(
		NewGrpcClient,
		fx.ResultTags(resultTags),
	)
}

func NewGrpcClient(listener *bufconn.Listener) (grpc.ClientConnInterface, error) {
	ctx := context.Background()

	return grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
		return listener.Dial()
	}), grpc.WithTransportCredentials(insecure.NewCredentials()))
}
