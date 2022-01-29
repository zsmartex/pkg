package services

import (
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaClient struct {
	config   *kafka.ConfigMap
	consumer *kafka.Consumer
	producer *kafka.Producer
}

func NewKafka() *KafkaClient {
	config := &kafka.ConfigMap{
		"bootstrap.servers":  os.Getenv("KAFKA_URL"),
		"auto.commit.offset": false,
		"group.id":           os.Getenv("KAFKA_GROUP_ID"),
	}

	return &KafkaClient{
		config: config,
	}
}

func (k *KafkaClient) Subscribe(topic string, callback func(c *kafka.Consumer, e kafka.Event) error) error {
	if k.consumer == nil {
		consumer, err := kafka.NewConsumer(k.config)
		if err != nil {
			panic("Can't create consumer due to error: " + err.Error())
		}

		k.consumer = consumer
	}

	k.consumer.Poll(100)

	return k.consumer.Subscribe(topic, func(c *kafka.Consumer, e kafka.Event) error {
		err := callback(c, e)

		if err != nil {
			c.Commit()
		}

		return err
	})
}

func (k *KafkaClient) Publish(topic string, body []byte) error {
	if k.producer == nil {
		producer, err := kafka.NewProducer(k.config)
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
