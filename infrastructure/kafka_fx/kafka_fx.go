package kafka_fx

import (
	"context"

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
}

func registerSubscriberHooks(lc fx.Lifecycle, params SubscriberHooksParams) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go params.Consumer.Subscribe(params.Subscriber)
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
