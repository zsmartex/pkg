package services

import (
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaClient struct {
	Consumer *kafka.Consumer
	Producer *kafka.Producer
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

func (k *KafkaClient) Subscribe(topic string, callback func([]byte) error) {
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

		if err := callback(msg.Value); err == nil {
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

func (k *KafkaClient) Publish(topic string, body []byte) error {
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
		Value:          body,
	}, nil)

	if err != nil {
		return err
	}

	k.Producer.Flush(100)

	return nil
}
