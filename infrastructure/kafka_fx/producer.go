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

func (p *Producer) produce(context context.Context, topic Topic, key []byte, payload interface{}) error {
	switch data := payload.(type) {
	case string:
		return p.produce(context, topic, key, []byte(data))
	case []byte:
		log.Debugf("Kafka producer produce to: %s, key: %s, payload: %s", topic, key, payload)

		res := p.client.ProduceSync(context, &kgo.Record{
			Topic: string(topic),
			Key:   []byte(key),
			Value: data,
		})

		if err := res.FirstErr(); err != nil {
			log.Errorf("Kafka producer produce to: %s, key: %s, payload: %s, error: %s", topic, key, payload, err)
		}

		return nil
	default:
		data, err := json.Marshal(payload)
		if err != nil {
			return err
		}

		return p.produce(context, topic, key, data)
	}
}

func (p *Producer) Health(ctx context.Context) error {
	return p.client.Ping(ctx)
}

func (p *Producer) Close() {
	p.client.Close()
}
