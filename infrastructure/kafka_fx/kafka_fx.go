package kafka_fx

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kmsg"
	"go.uber.org/fx"

	"github.com/zsmartex/pkg/v2/config"
	"github.com/zsmartex/pkg/v2/log"
)

var (
	ConsumerModule = fx.Module("kafka_fx_consumer", consumerProviders, consumerInvokes)
	ProducerModule = fx.Module("kafka_fx_producer", producerProviders, producerInvokes)

	consumerProviders = fx.Provide(NewConsumer)
	producerProviders = fx.Provide(NewProducer)

	consumerInvokes = fx.Options(fx.Invoke(registerConsumerHooks), fx.Invoke(registerSubscriberHooks))
	producerInvokes = fx.Options(fx.Invoke(registerProducerHooks))
)

type alterTopicParams struct {
	Topic       Topic
	Config      config.Kafka
	Consumer    *Consumer
	AdminClient *kadm.Client
}

func alterTopic(ctx context.Context, params alterTopicParams) error {
	topicDetails, err := params.AdminClient.ListTopics(ctx)
	if err != nil {
		return err
	}

	if topicDetails.Has(string(params.Topic)) {
		replicationFactor := fmt.Sprintf("%d", params.Config.ReplicationFactor)

		_, err := params.AdminClient.AlterTopicConfigs(ctx, []kadm.AlterConfig{
			{
				Op:    kadm.SetConfig,
				Name:  "replication.factor",
				Value: &replicationFactor,
			},
		}, string(params.Topic))
		if err != nil {
			return err
		}

		req := kmsg.NewPtrMetadataRequest()
		reqTopic := kmsg.NewMetadataRequestTopic()
		reqTopic.Topic = kmsg.StringPtr(string(params.Topic))
		req.Topics = append(req.Topics, reqTopic)
		resp, err := req.RequestWith(context.Background(), params.Consumer.client)
		if err != nil {
			return err
		}

		if len(resp.Topics) == 0 {
			return errors.New("no topics found")
		}

		t := resp.Topics[0]

		if len(t.Partitions) < int(params.Config.Partitions) {
			remainTopics := int(params.Config.Partitions) - len(t.Partitions)
			_, err := params.AdminClient.CreatePartitions(ctx, remainTopics, string(params.Topic))
			if err != nil {
				return err
			}
		}
	} else {
		_, err := params.AdminClient.CreateTopic(
			ctx,
			params.Config.Partitions,
			params.Config.ReplicationFactor,
			nil,
			string(params.Topic),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

type registerConsumerHooksParams struct {
	fx.In

	Topic       Topic `optional:"true"`
	Config      config.Kafka
	Consumer    *Consumer
	AdminClient *kadm.Client
}

func registerConsumerHooks(
	params registerConsumerHooksParams,
	lc fx.Lifecycle,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if params.Topic == "" {
				return nil
			}

			return alterTopic(ctx, alterTopicParams{
				Topic:       params.Topic,
				Config:      params.Config,
				Consumer:    params.Consumer,
				AdminClient: params.AdminClient,
			})
		},
		OnStop: func(ctx context.Context) error {
			params.Consumer.Close()

			return nil
		},
	})
}

type SubscriberHooksParams struct {
	fx.In

	Consumer   *Consumer
	Subscriber ConsumerSubscriber
	Ticker     *time.Ticker
}

func registerSubscriberHooks(lc fx.Lifecycle, params SubscriberHooksParams) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Kafka consumer subscriber started listening")

			go params.Consumer.Subscribe(params.Subscriber, params.Ticker)
			return nil
		},
	})
}

func registerProducerHooks(
	lc fx.Lifecycle,
	kafkaProducer *Producer,
) {
	lc.Append(fx.StopHook(func(ctx context.Context) error {
		kafkaProducer.Close()

		return nil
	}))
}
