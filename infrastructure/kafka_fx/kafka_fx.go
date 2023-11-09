package kafka_fx

import (
	"context"
	"time"

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

func registerConsumerHooks(
	lc fx.Lifecycle,
	kafkaConsumer *Consumer,
) {
	lc.Append(fx.StopHook(func(ctx context.Context) error {
		kafkaConsumer.Close()

		return nil
	}))
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
