package services

import (
	"strings"

	"github.com/zsmartex/pkg"
)

type RangoClient struct {
	producer *KafkaProducer
}

func NewRangoClient(producer *KafkaProducer) (*RangoClient, error) {
	return &RangoClient{producer: producer}, nil
}

func (k *RangoClient) EnqueueEvent(kind pkg.EnqueueEventKind, id, event string, payload interface{}) {
	key := strings.Join([]string{string(kind), id, event}, ".")

	k.producer.ProduceWithKey("rango.events", key, payload)
}
