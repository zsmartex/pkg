package kafka_fx

import (
	"context"
	"fmt"
	"time"

	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/zsmartex/pkg/v2/log"
	"go.uber.org/fx"
)

var (
	ConsumerModule = fx.Module("kafka_fx_consumer", consumerProviders, consumerInvokes)
	ProducerModule = fx.Module("kafka_fx_producer", producerProviders, producerInvokes)

	consumerProviders = fx.Provide(NewConsumer)
	producerProviders = fx.Provide(NewProducer)

	consumerInvokes = fx.Options(fx.Invoke(registerConsumerHooks), fx.Invoke(registerSubscriberHooks))
	producerInvokes = fx.Options(fx.Invoke(registerProducerHooks))
)

type registerConsumerHooksParams struct {
	fx.In

	Topic             Topic
	Consumer          *Consumer
	AdminClient       *kadm.Client
	ReplicationFactor ReplicationFactor
}

func registerConsumerHooks(
	params registerConsumerHooksParams,
	lc fx.Lifecycle,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			topicDetails, err := params.AdminClient.ListTopics(ctx)
			if err != nil {
				return err
			}

			if topicDetails.Has(string(params.Topic)) {
				replicationFactor := fmt.Sprintf("%d", params.ReplicationFactor)

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
			} else {
				_, err := params.AdminClient.CreateTopic(ctx, 3, int16(params.ReplicationFactor), nil, string(params.Topic))
				if err != nil {
					return err
				}
			}

			return nil
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
