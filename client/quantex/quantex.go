package clientQuantex

import (
	"context"
	"log"
	"os"
	"time"

	GrpcQuantex "github.com/zsmartex/pkg/Grpc/quantex"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GrpcQuantexClient struct {
	connect *grpc.ClientConn
	client  GrpcQuantex.QuantexServiceClient
}

func NewQuantexClient() *GrpcQuantexClient {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(os.Getenv("QUANTEX_ENGINE_URL"), grpc.WithInsecure(), grpc.WithBackoffMaxDelay(5*time.Second))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	return &GrpcQuantexClient{
		connect: conn,
		client:  GrpcQuantex.NewQuantexServiceClient(conn),
	}
}

func (c *GrpcQuantexClient) UpdateOrder(in *GrpcQuantex.UpdateOrderRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return c.client.UpdateOrder(context.Background(), in, opts...)
}

func (c *GrpcQuantexClient) ReloadStrategy(in *GrpcQuantex.ReloadStrategyRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return c.client.ReloadStrategy(context.Background(), in, opts...)
}

func (c *GrpcQuantexClient) Close() error {
	return c.connect.Close()
}
