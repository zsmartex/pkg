package services

import (
	"context"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/google/uuid"
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
		GroupId:          "zsmartex-" + uuid.NewString(),
		Topics:           strings.Join(topics, ", "),
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

func (k *KafkaClient) CreateProducer() (*kafka.Producer, error) {
	return kafka.NewProducer(kafka.ProducerConfig{
		BrokersList:  os.Getenv("KAFKA_URL"),
		RequiredAcks: kafka.WaitForLocal,
		IsCompressed: true,
		Logger:       k.logger,
	})
}

func (k *KafkaClient) publishJSON(topic, key string, payload interface{}) error {
	k.publishMutex.Lock()
	defer k.publishMutex.Unlock()

	if k.Producer == nil {
		producer, err := k.CreateProducer()
		if err != nil {
			panic("Can't create producer due to error: " + err.Error())
		}

		k.Producer = producer
	}

	if len(key) > 0 {
		return k.Producer.ProduceJSON(topic, payload)
	} else {
		return k.Producer.ProduceJSONWithKey(topic, payload, key)
	}
}

func (k *KafkaClient) publish(topic string, key string, payload []byte) error {
	k.publishMutex.Lock()
	defer k.publishMutex.Unlock()

	if k.Producer == nil {
		producer, err := k.CreateProducer()
		if err != nil {
			panic("Can't create producer due to error: " + err.Error())
		}

		k.Producer = producer
	}

	if len(key) > 0 {
		return k.Producer.Produce(topic, payload)
	} else {
		return k.Producer.ProduceWithKey(topic, payload, key)
	}
}

func (k *KafkaClient) Publish(topic string, payload []byte) error {
	return k.publish(topic, "", payload)
}

func (k *KafkaClient) PublishWithKey(topic, key string, payload []byte) error {
	return k.publish(topic, "", payload)
}

func (k *KafkaClient) PublishJSON(topic string, payload interface{}) error {
	return k.publishJSON(topic, "", payload)
}

func (k *KafkaClient) PublishJSONWithKey(topic, key string, payload interface{}) error {
	return k.publishJSON(topic, key, payload)
}

func (k *KafkaClient) EnqueueEvent(kind pkg.EnqueueEventKind, id, event string, payload interface{}) error {
	key := strings.Join([]string{string(kind), id, event}, ".")

	return k.publishJSON("rango:events", key, payload)
}
