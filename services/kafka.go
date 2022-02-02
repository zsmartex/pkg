package services

import (
	"context"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/zsmartex/pkg"
	"github.com/zsmartex/pkg/wrap/kafka"
)

type KafkaClient struct {
	Consumer     *kafka.Consumer
	Producer     *kafka.Producer
	publishMutex sync.Mutex
	logger       *logrus.Entry
}

func NewKafka(logger *logrus.Entry) *KafkaClient {
	return &KafkaClient{
		logger: logger,
	}
}

func (k *KafkaClient) CreateConsumer(topics []string) (*kafka.Consumer, error) {
	return kafka.NewConsumer(kafka.ConsumerConfig{
		BootstrapServers: os.Getenv("KAFKA_URL"),
		Offset:           kafka.OffsetEarliest,
		GroupId:          "zsmartex",
		Topics:           strings.Join(topics, ","),
		Logger:           k.logger,
	})
}

func (k *KafkaClient) Subscribe(topics []string, callback func(msg kafka.Message) error) {
	if k.Consumer == nil {
		consumer, err := k.CreateConsumer(topics)
		if err != nil {
			panic("Can't create consumer due to error: " + err.Error())
		}

		k.Consumer = consumer
	}

	for {
		messages, err := k.Consumer.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Consumer error: %v (%v)\n", err, messages)
		}

		for _, msg := range messages {
			if err := callback(msg); err == nil {
				msg.Session.MarkMessage(msg.SamMsg, "")
			}
		}
	}
}

func (k *KafkaClient) publish(topic string, key string, body []byte) error {
	k.publishMutex.Lock()
	defer k.publishMutex.Unlock()

	if k.Producer == nil {
		producer, err := kafka.NewProducer(kafka.ProducerConfig{
			BrokersList:  os.Getenv("KAFKA_URL"),
			RequiredAcks: kafka.WaitForAll,
			IsCompressed: true,
			Logger:       k.logger,
		})
		if err != nil {
			panic("Can't create producer due to error: " + err.Error())
		}

		k.Producer = producer
	}

	if len(key) > 0 {
		return k.Producer.Produce(topic, body)
	} else {
		return k.Producer.ProduceWithKey(topic, body, key)
	}
}

func (k *KafkaClient) Publish(topic string, body []byte) error {
	return k.publish(topic, "", body)
}

func (k *KafkaClient) EnqueueEvent(kind pkg.EnqueueEventKind, id, event string, payload []byte) error {
	key := strings.Join([]string{string(kind), id, event}, ".")

	return k.publish("rango:events", key, payload)
}
