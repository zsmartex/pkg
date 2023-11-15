package kafka_fx

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/zsmartex/pkg/v2/config"
	"github.com/zsmartex/pkg/v2/log"
	"go.uber.org/fx"
)

type Topic string
type Group string

type Consumer struct {
	client *kgo.Client
}

type consumerParams struct {
	fx.In

	Config       config.Kafka
	Topic        Topic
	Group        Group `optional:"true"`
	AtEnd        bool  `name:"at_end" optional:"true"`
	ManualCommit bool  `name:"manual_commit" optional:"true"`
}

func NewConsumer(params consumerParams) (*Consumer, *kadm.Client, error) {
	options := []kgo.Opt{
		kgo.SeedBrokers(params.Config.Brokers...),
		kgo.AllowAutoTopicCreation(),
		kgo.ConsumeTopics(string(params.Topic)),
	}

	if params.ManualCommit {
		options = append(options, kgo.DisableAutoCommit())
	}

	if params.Group == "" {
		options = append(options, kgo.ConsumeResetOffset(kgo.NewOffset().AtEnd()))
	} else {
		options = append(options, kgo.ConsumerGroup(string(params.Group)))

		if params.AtEnd {
			options = append(options, kgo.ConsumeResetOffset(kgo.NewOffset().AtEnd()))
		}
	}

	client, err := kgo.NewClient(options...)
	if err != nil {
		return nil, nil, err
	}

	adminClient := kadm.NewClient(client)

	return &Consumer{
		client: client,
	}, adminClient, nil
}

func (c *Consumer) Poll(ctx context.Context) ([]*kgo.Record, error) {
	fetches := c.client.PollFetches(ctx)
	if err := fetches.Err(); err != nil {
		return nil, errors.Wrap(err, "kafka consumer poll")
	}

	return fetches.Records(), nil
}

type ConsumerSubscriber interface {
	OnMessage(*kgo.Record) error
}

func (c *Consumer) Subscribe(subscriber ConsumerSubscriber, ticker *time.Ticker) error {
	for range ticker.C {
		records, err := c.Poll(context.Background())
		if err != nil {
			return err
		}

		for _, record := range records {
			if record.Key != nil {
				log.Debugf("kafka consumer received record with key: %s, value: %s", record.Key, record.Value)
			} else {
				log.Debugf("kafka consumer received record with value: %s", record.Value)
			}

			if err := subscriber.OnMessage(record); err != nil {
				log.Errorf("kafka consumer error: %s", err)
			}
		}
	}

	return nil
}

func (c *Consumer) CommitRecords(context context.Context, records ...*kgo.Record) error {
	return c.client.CommitRecords(context, records...)
}

func (c *Consumer) Close() {
	c.client.Close()
}
