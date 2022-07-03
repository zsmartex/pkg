package kafka

import (
	"context"
	"encoding/json"

	"github.com/sirupsen/logrus"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/zsmartex/pkg/log"
)

type KafkaConsumer struct {
	CommitClient *kadm.Client
	Client       *kgo.Client
	Topics       []string
	Group        string
}

func NewKafkaConsumer(brokers []string, group string, topics []string) (*KafkaConsumer, error) {
	seeds := kgo.SeedBrokers(brokers...)

	cl, err := kgo.NewClient(
		seeds,
	)
	if err != nil {
		return nil, err
	}

	var client *kgo.Client

	adm := kadm.NewClient(cl)
	os, err := adm.FetchOffsetsForTopics(context.Background(), group, topics...)

	if os.Ok() && err == nil {
		client, err = kgo.NewClient(
			seeds,
			kgo.ConsumePartitions(os.KOffsets()),
		)
		if err != nil {
			return nil, err
		}
	} else {
		client, err = kgo.NewClient(
			seeds,
			kgo.AllowAutoTopicCreation(),
			kgo.ConsumerGroup(group),
			kgo.ConsumeTopics(topics...),
			kgo.DisableAutoCommit(),
		)
		if err != nil {
			return nil, err
		}
	}

	return &KafkaConsumer{
		CommitClient: adm,
		Client:       client,
		Topics:       topics,
		Group:        group,
	}, nil
}

func (c *KafkaConsumer) Poll() ([]*kgo.Record, error) {
	fetches := c.Client.PollFetches(context.Background())
	if err := fetches.Err(); err != nil {
		return nil, err
	}

	return fetches.Records(), nil
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

func NewKafkaProducer(brokers []string, logger *logrus.Entry) (*KafkaProducer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.AllowAutoTopicCreation(),
	)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		Client: client,
	}, nil
}

func (k *KafkaProducer) Produce(topic string, payload interface{}) {
	k.produce(topic, "", payload)
}

func (k *KafkaProducer) ProduceWithKey(topic, key string, payload interface{}) {
	k.produce(topic, key, payload)
}

func (p *KafkaProducer) produce(topic, key string, payload interface{}) {
	switch data := payload.(type) {
	case string:
		p.produce(topic, key, []byte(data))
		return
	case []byte:
		log.Debugf("Kafka producer produce to: %s, key: %s, payload: %s", topic, key, payload)

		p.Client.Produce(context.Background(), &kgo.Record{
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
		jdata, err := json.Marshal(payload)
		if err != nil {
			return
		}

		p.produce(topic, key, jdata)
	}
}
