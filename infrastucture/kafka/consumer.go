package kafka

import (
	"context"

	"github.com/cockroachdb/errors"
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
		return nil, errors.Wrap(err, "kafka consumer poll")
	}

	return fetches.Records(), nil
}

type ConsumerSubscriber interface {
	OnMessage(key []byte, message []byte) error
}

func (c *Consumer) Subscribe(subscriber ConsumerSubscriber) error {
	for {
		records, err := c.Poll(context.Background())
		if err != nil {
			return err
		}

		for _, record := range records {
			if err := subscriber.OnMessage(record.Key, record.Value); err != nil {
				log.Errorf("kafka consumer error: %s", err)
			}
		}
	}
}

func (c *Consumer) CommitRecords(context context.Context, records ...*kgo.Record) error {
	return c.client.CommitRecords(context, records...)
}

func (c *Consumer) Close() {
	c.client.Close()
}
