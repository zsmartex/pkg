package middlewares

import (
	"context"

	"github.com/zsmartex/pkg/v2/log"
	"google.golang.org/grpc"
)

func ErrorHandler(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Call the gRPC handler function
	resp, err := handler(ctx, req)
	if err != nil {
		log.Errorf("%+v\n", err)
	}

	// Return the response
	return resp, nil
}
