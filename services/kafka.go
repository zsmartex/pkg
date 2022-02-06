package services

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

type KafkaConsumer struct {
	CommitClient *kadm.Client
	Client       *kgo.Client
	Topics       []string
	Group        string
}

func NewKafkaConsumer(group string, topics []string) (*KafkaConsumer, error) {
	brokers := getBrokers()
	seeds := kgo.SeedBrokers(brokers...)

	cl, err := kgo.NewClient(
		seeds,
	)
	if err != nil {
		return nil, err
	}

	adm := kadm.NewClient(cl)
	os, err := adm.FetchOffsetsForTopics(context.Background(), group, topics...)
	if err != nil {
		return nil, err
	}

	client, err := kgo.NewClient(seeds, kgo.ConsumePartitions(os.Into().Into()))
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		CommitClient: adm,
		Client:       client,
		Topics:       topics,
		Group:        group,
	}, nil
}

func (c *KafkaConsumer) Poll() ([]*kgo.Record, error) {
	records := make([]*kgo.Record, 0)
	errors := make([]error, 0)

	fetches := c.Client.PollRecords(context.Background(), -1)

	fetches.EachError(func(s string, i int32, e error) {
		errors = append(errors, e)
	})

	if len(errors) > 0 {
		return records, errors[0]
	}

	fetches.EachRecord(func(r *kgo.Record) {
		records = append(records, r)
	})

	return records, nil
}

func (c *KafkaConsumer) CommitRecords(records ...kgo.Record) error {
	return c.CommitClient.CommitAllOffsets(context.Background(), c.Group, kadm.OffsetsFromRecords(records...))
}

func (c *KafkaConsumer) Close() {
	c.Client.Close()
	c.CommitClient.Close()
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

	return k.produce(topic, key, data)
}

func getBrokers() []string {
	return strings.Split(os.Getenv("KAFKA_URL"), ",")
}

func (p *KafkaProducer) produce(topic, key string, payload []byte) error {
	var bkey []byte

	if len(key) > 0 {
		bkey = []byte(key)
	}

	r := p.Client.ProduceSync(context.Background(), &kgo.Record{
		Topic: topic,
		Key:   bkey,
		Value: payload,
	})

	return r.FirstErr()
}
