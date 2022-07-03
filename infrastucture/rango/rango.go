package rango

import (
	"strings"

	"github.com/zsmartex/pkg"
	"github.com/zsmartex/pkg/infrastucture/kafka"
)

type RangoClient struct {
	producer *kafka.KafkaProducer
}

func NewRangoClient(producer *kafka.KafkaProducer) (*RangoClient, error) {
	return &RangoClient{producer: producer}, nil
}

func (k *RangoClient) EnqueueEvent(kind pkg.EnqueueEventKind, id, event string, payload interface{}) {
	key := strings.Join([]string{string(kind), id, event}, ".")

	k.producer.ProduceWithKey("rango.events", key, payload)
}
