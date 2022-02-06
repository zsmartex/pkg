package services

import (
	"log"
	"strings"

	"github.com/zsmartex/pkg"
)

type RangoClient struct {
	producer *KafkaProducer
}

func NewRangoClient() (*RangoClient, error) {
	producer, err := NewKafkaProducer()
	if err != nil {
		return nil, err
	}

	return &RangoClient{producer: producer}, nil
}

func (k *RangoClient) EnqueueEvent(kind pkg.EnqueueEventKind, id, event string, payload interface{}) error {
	log.Println(k.producer)
	log.Println(k.producer.Client)
	key := strings.Join([]string{string(kind), id, event}, ".")

	return k.producer.ProduceWithKey("rango.events", key, payload)
}
