package rango

import (
	"context"
	"strings"

	"github.com/zsmartex/pkg/v2"
	"github.com/zsmartex/pkg/v2/infrastucture/kafka"
)

type Client struct {
	producer *kafka.Producer
}

func NewClient(producer *kafka.Producer) (*Client, error) {
	return &Client{producer: producer}, nil
}

func (k *Client) EnqueueEvent(context context.Context, kind pkg.EnqueueEventKind, id, event string, payload interface{}) {
	key := strings.Join([]string{string(kind), id, event}, ".")

	k.producer.ProduceWithKey(context, "rango.events", key, payload)
}
