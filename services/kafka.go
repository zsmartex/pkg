package services

import (
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaClient struct {
	consumer *kafka.Consumer
	producer *kafka.Producer
}

func NewKafka() *KafkaClient {
	return &KafkaClient{}
}

func (k *KafkaClient) Subscribe(topic string, callback func([]byte) error) {
	if k.consumer == nil {
		consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers":  os.Getenv("KAFKA_URL"),
			"enable.auto.commit": false,
			"group.id":           os.Getenv("KAFKA_GROUP_ID"),
		})
		if err != nil {
			panic("Can't create consumer due to error: " + err.Error())
		}

		k.consumer = consumer
	}

	k.consumer.SubscribeTopics([]string{topic}, nil)

	for {
		msg, err := k.consumer.ReadMessage(-1)
		if err != nil {
			log.Printf("Consumer error: %v (%v)\n", err, msg)
		}

		if err := callback(msg.Value); err == nil {
			k.consumer.Commit()
		}
	}
}

func (k *KafkaClient) Publish(topic string, body []byte) error {
	if k.producer == nil {
		producer, err := kafka.NewProducer(&kafka.ConfigMap{
			"bootstrap.servers": os.Getenv("KAFKA_URL"),
		})
		if err != nil {
			panic("Can't create producer due to error: " + err.Error())
		}

		k.producer = producer
	}

	err := k.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          body,
	}, nil)

	if err != nil {
		return err
	}

	k.producer.Flush(100)

	return nil
}
