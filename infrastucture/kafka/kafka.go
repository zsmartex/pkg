package kafka

import (
	"context"
	"encoding/json"

	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/zsmartex/pkg/v3/log"
)

type Consumer struct {
	CommitClient *kadm.Client
	Client       *kgo.Client
	Topics       []string
	Group        string
}

func NewConsumer(context context.Context, brokers []string, group string, topics ...string) (*Consumer, error) {
	seeds := kgo.SeedBrokers(brokers...)

	cl, err := kgo.NewClient(
		seeds,
		kgo.AllowAutoTopicCreation(),
		kgo.ConsumerGroup(group),
		kgo.ConsumeTopics(topics...),
		kgo.DisableAutoCommit(),
	)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		CommitClient: kadm.NewClient(cl),
		Client:       cl,
		Topics:       topics,
		Group:        group,
	}, nil
}

func (c *Consumer) Poll(ctx context.Context) ([]*kgo.Record, error) {
	fetches := c.Client.PollFetches(ctx)
	if err := fetches.Err(); err != nil {
		return nil, err
	}

	return fetches.Records(), nil
}

func (c *Consumer) CommitRecords(context context.Context, records ...*kgo.Record) error {
	return c.Client.CommitRecords(context, records...)
}

func (c *Consumer) Close() {
	c.Client.Close()
	c.CommitClient.Close()
}

type Producer struct {
	Client *kgo.Client
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
		Client: client,
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

		p.Client.Produce(context, &kgo.Record{
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
	return p.Client.Ping(ctx)
}
