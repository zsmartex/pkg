package kafka_fx

import (
	"context"
	"encoding/json"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/fx"

	"github.com/zsmartex/pkg/v2/config"
	"github.com/zsmartex/pkg/v2/log"
)

type Producer struct {
	client *kgo.Client
}

type producerParams struct {
	fx.In

	Config config.Kafka
}

func NewProducer(params producerParams) (*Producer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(params.Config.Brokers...),
		kgo.AllowAutoTopicCreation(),
	)
	if err != nil {
		return nil, err
	}

	return &Producer{
		client: client,
	}, nil
}

func (k *Producer) Produce(context context.Context, topic Topic, payload interface{}) error {
	return k.produce(context, topic, nil, payload)
}

func (k *Producer) ProduceWithKey(context context.Context, topic Topic, key []byte, payload interface{}) error {
	return k.produce(context, topic, key, payload)
}

func (k *Producer) ProduceSync(context context.Context, topic Topic, payload interface{}) error {
	return k.produceSync(context, topic, nil, payload)
}

func (k *Producer) ProduceSyncWithKey(context context.Context, topic Topic, key []byte, payload interface{}) error {
	return k.produceSync(context, topic, key, payload)
}

func (p *Producer) produce(context context.Context, topic Topic, key []byte, payload interface{}) error {
	data, err := p.marshalPayload(payload)
	if err != nil {
		return err
	}

	log.Debugf("Kafka producer produce to: %s, key: %s, payload: %s", topic, key, payload)

	p.client.Produce(context, &kgo.Record{
		Topic: string(topic),
		Key:   []byte(key),
		Value: data,
	}, func(r *kgo.Record, err error) {
		if err != nil {
			log.Errorf("Failed to produce message to topic: %s, key: %s, payload: %s, error: %s", topic, key, payload, err)
		}
	})

	return nil
}

func (p *Producer) produceSync(context context.Context, topic Topic, key []byte, payload interface{}) error {
	data, err := p.marshalPayload(payload)
	if err != nil {
		return err
	}

	log.Debugf("Kafka producer produce to: %s, key: %s, payload: %s", topic, key, payload)

	res := p.client.ProduceSync(context, &kgo.Record{
		Topic: string(topic),
		Key:   []byte(key),
		Value: data,
	})

	if err := res.FirstErr(); err != nil {
		log.Errorf("Failed to produce message to topic: %s, key: %s, payload: %s, error: %s", topic, key, payload, err)
	}

	return nil
}

func (p *Producer) marshalPayload(payload interface{}) ([]byte, error) {
	switch data := payload.(type) {
	case string:
		return []byte(data), nil
	case []byte:
		return data, nil
	default:
		return json.Marshal(payload)
	}
}

func (p *Producer) Health(ctx context.Context) error {
	return p.client.Ping(ctx)
}

func (p *Producer) Close() {
	p.client.Close()
}
