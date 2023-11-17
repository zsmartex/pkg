package grpc_fx

import (
	"context"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/zsmartex/pkg/v2/config"
)

func AsClientParams(paramsTags, resultTags string) interface{} {
	return fx.Annotate(
		NewGrpcClient,
		fx.ParamTags(paramsTags),
		fx.ResultTags(resultTags),
	)
}

func NewGrpcClient(ctx context.Context, config config.GRPC) (grpc.ClientConnInterface, error) {
	conn, err := grpc.DialContext(ctx, config.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	return conn, nil
}
