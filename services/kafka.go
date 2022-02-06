package services

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/twmb/franz-go/pkg/kgo"
)

type KafkaClient struct {
	Client *kgo.Client
}

func NewKafka() (*KafkaClient, error) {
	client, err := kgo.NewClient()
	if err != nil {
		return nil, err
	}
	return &KafkaClient{
		Client: client,
	}, nil
}

type KafkaConsumer struct {
	Client *kgo.Client
}

func NewKafkaConsumer(topics ...string) (*KafkaConsumer, error) {
	brokers := getBrokers()

	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup("zsmartex-"+uuid.NewString()),
		kgo.ConsumeTopics(topics...),
	)
	if err != nil {
		return nil, err
	}
	return &KafkaConsumer{
		Client: client,
	}, nil
}

func (c *KafkaConsumer) Poll() []*kgo.Record {
	records := make([]*kgo.Record, 0)

	fetches := c.Client.PollRecords(context.Background(), -1)

	fetches.EachRecord(func(r *kgo.Record) {
		records = append(records, r)
	})

	return records
}

func (c *KafkaConsumer) CommitRecords(records ...*kgo.Record) error {
	return c.Client.CommitRecords(context.Background(), records...)
}

type KafkaProducer struct {
	Client *kgo.Client
}

func NewKafkaProducer() (*KafkaProducer, error) {
	brokers := getBrokers()

	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
	)
	if err != nil {
		return nil, err
	}
	return &KafkaProducer{
		Client: client,
	}, nil
}

func (k *KafkaProducer) Produce(topic string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return k.produce(topic, "", data)
}

func (k *KafkaProducer) ProduceWithKey(topic, key string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return k.produce(topic, "", data)
}

func getBrokers() []string {
	return strings.Split(os.Getenv("KAFKA_URL"), ",")
}

func (p *KafkaProducer) produce(topic, key string, payload []byte) error {
	r := p.Client.ProduceSync(context.Background(), &kgo.Record{
		Topic: topic,
	}, nil)

	return r.FirstErr()
}
