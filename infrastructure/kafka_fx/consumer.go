package kafka_fx

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/fx"

	"github.com/zsmartex/pkg/v2/config"
	"github.com/zsmartex/pkg/v2/log"
)

type Topic string
type Group string

type Consumer struct {
	client       *kgo.Client
	config       config.Kafka
	adminClient  *kadm.Client
	manualCommit bool
}

type consumerParams struct {
	fx.In

	Config       config.Kafka
	Topics       []Topic `optional:"true"`
	Group        Group   `optional:"true"`
	AtEnd        bool    `name:"at_end" optional:"true"`
	ManualCommit bool    `name:"manual_commit" optional:"true"`
}

func NewConsumer(params consumerParams) (*Consumer, *kadm.Client, error) {
	topics := make([]string, len(params.Topics))
	for i, topic := range params.Topics {
		topics[i] = string(topic)
	}

	options := []kgo.Opt{
		kgo.SeedBrokers(params.Config.Brokers...),
		kgo.AllowAutoTopicCreation(),
		kgo.ConsumeTopics(topics...),
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
		client:       client,
		config:       params.Config,
		adminClient:  adminClient,
		manualCommit: params.ManualCommit,
	}, adminClient, nil
}

func (c *Consumer) AddConsumeTopics(ctx context.Context, topics ...Topic) error {
	var strTopics []string
	for _, topic := range topics {
		strTopics = append(strTopics, string(topic))
	}

	for _, topic := range topics {
		err := alterTopic(ctx, alterTopicParams{
			Topic:       topic,
			Config:      c.config,
			Consumer:    c,
			AdminClient: c.adminClient,
		})
		if err != nil {
			return err
		}
	}

	c.client.AddConsumeTopics(strTopics...)
	return nil
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
				log.Errorf("kafka consumer handle message error: %s", err)
			} else if c.manualCommit {
				err := c.CommitRecords(context.Background(), record)
				if err != nil {
					log.Errorf("kafka consumer commit error: %s", err)
				}
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
