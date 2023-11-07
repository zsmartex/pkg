package grpc_fx

import (
	"github.com/zsmartex/pkg/v2/config"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func AsClientParams(paramsTags, resultTags string) interface{} {
	return fx.Annotate(
		NewGrpcClient,
		fx.ParamTags(paramsTags),
		fx.ResultTags(resultTags),
	)
}

func NewGrpcClient(config config.GRPC) (grpc.ClientConnInterface, error) {
	conn, err := grpc.Dial(config.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return conn, nil
}
