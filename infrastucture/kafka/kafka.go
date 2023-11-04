package kafka

import (
	"context"
	"encoding/json"

	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/zsmartex/pkg/v2/log"
)

type Consumer struct {
	client *kgo.Client
}

func NewConsumer(opts ...kgo.Opt) (*Consumer, error) {
	client, err := kgo.NewClient()
	if err != nil {
		return nil, err
	}

	return &Consumer{
		client: client,
	}, nil
}

func (c *Consumer) Poll(ctx context.Context) ([]*kgo.Record, error) {
	fetches := c.client.PollFetches(ctx)
	if err := fetches.Err(); err != nil {
		return nil, err
	}

	return fetches.Records(), nil
}

func (c *Consumer) CommitRecords(context context.Context, records ...*kgo.Record) error {
	return c.client.CommitRecords(context, records...)
}

func (c *Consumer) Close() {
	c.client.Close()
}

type Producer struct {
	client *kgo.Client
}

func NewProducer(brokers []string) (*Producer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.AllowAutoTopicCreation(),
	)
	if err != nil {
		return nil, err
	}

	return &Producer{
		client: client,
	}, nil
}

func (k *Producer) Produce(context context.Context, topic string, payload interface{}) {
	k.produce(context, topic, "", payload)
}

func (k *Producer) ProduceWithKey(context context.Context, topic, key string, payload interface{}) {
	k.produce(context, topic, key, payload)
}

func (p *Producer) produce(context context.Context, topic, key string, payload interface{}) {
	switch data := payload.(type) {
	case string:
		p.produce(context, topic, key, []byte(data))
		return
	case []byte:
		log.Debugf("Kafka producer produce to: %s, key: %s, payload: %s", topic, key, payload)

		p.client.Produce(context, &kgo.Record{
			Topic: topic,
			Key:   []byte(key),
			Value: data,
		}, func(r *kgo.Record, err error) {
			if err != nil {
				log.Errorf("Kafka producer produce to: %s, key: %s, payload: %s, error: %s", topic, key, payload, err)
			}
		})
		return
	default:
		data, err := json.Marshal(payload)
		if err != nil {
			return
		}

		p.produce(context, topic, key, data)
	}
}

func (p Producer) Health(ctx context.Context) error {
	return p.client.Ping(ctx)
}
