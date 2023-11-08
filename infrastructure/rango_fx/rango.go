package rango_fx

import (
	"context"
	"strings"

	"github.com/zsmartex/pkg/v2/infrastructure/kafka_fx"
)

type PublishKind string

var (
	PublishKindPublic  PublishKind = "public"
	PublishKindPrivate PublishKind = "private"
	PublishKindSystem  PublishKind = "system"
)

type Client struct {
	producer *kafka_fx.Producer
}

func New(producer *kafka_fx.Producer) (*Client, error) {
	return &Client{producer: producer}, nil
}

func (k *Client) Publish(context context.Context, kind PublishKind, id, event string, payload interface{}) error {
	key := strings.Join([]string{string(kind), id, event}, ".")

	return k.producer.ProduceWithKey(context, "rango.events", []byte(key), payload)
}
