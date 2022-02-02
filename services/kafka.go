package services

import (
	"log"
	"os"
	"strings"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/zsmartex/pkg"
)

type KafkaClient struct {
	Consumer     *kafka.Consumer
	Producer     *kafka.Producer
	publishMutex sync.Mutex
}

func NewKafka() *KafkaClient {
	return &KafkaClient{}
}

func (k *KafkaClient) CreateConsumer() (*kafka.Consumer, error) {
	return kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  os.Getenv("KAFKA_URL"),
		"enable.auto.commit": false,
		"group.id":           os.Getenv("KAFKA_GROUP_ID"),
	})
}

func (k *KafkaClient) Subscribe(topic string, callback func(msg *kafka.Message) error) {
	if k.Consumer == nil {
		consumer, err := k.CreateConsumer()
		if err != nil {
			panic("Can't create consumer due to error: " + err.Error())
		}

		k.Consumer = consumer
	}

	k.SubscribeTopics([]string{topic}, nil)

	for {
		msg, err := k.Consumer.ReadMessage(-1)
		if err != nil {
			log.Printf("Consumer error: %v (%v)\n", err, msg)
		}

		if err := callback(msg); err == nil {
			k.Consumer.CommitMessage(msg)
		}
	}
}

func (k *KafkaClient) SubscribeTopics(topics []string, rebalanceCb kafka.RebalanceCb) error {
	if k.Consumer == nil {
		consumer, err := k.CreateConsumer()
		if err != nil {
			panic("Can't create consumer due to error: " + err.Error())
		}

		k.Consumer = consumer
	}

	return k.Consumer.SubscribeTopics(topics, rebalanceCb)
}

func (k *KafkaClient) publish(topic string, key []byte, body []byte) error {
	k.publishMutex.Lock()
	defer k.publishMutex.Unlock()

	if k.Producer == nil {
		producer, err := kafka.NewProducer(&kafka.ConfigMap{
			"bootstrap.servers": os.Getenv("KAFKA_URL"),
		})
		if err != nil {
			panic("Can't create producer due to error: " + err.Error())
		}

		k.Producer = producer

		k.Producer.GetFatalError()
	}

	err := k.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
		Value:          body,
	}, nil)

	if err != nil {
		return err
	}

	k.Producer.Flush(5)

	return nil
}

func (k *KafkaClient) Publish(topic string, body []byte) error {
	return k.publish(topic, nil, body)
}

func (k *KafkaClient) EnqueueEvent(kind pkg.EnqueueEventKind, id, event string, payload []byte) error {
	key := strings.Join([]string{string(kind), id, event}, ".")

	return k.publish("rango:events", []byte(key), payload)
}
