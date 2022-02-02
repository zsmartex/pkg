package kafka

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// ProducerConfig _
type ProducerConfig struct {
	// BrokersList is comma separated: "broker1:9092,broker2:9092,broker3:9092"
	BrokersList string
	// level of acknowledgement reliability, default NoResponse
	RequiredAcks SendMsgReliabilityLevel

	//
	// following configs are optional
	//

	IsCompressed bool // if true, producer will use gzip level BestCompression
	// size before compress, default 1000000,
	// should <= broker's message.max.bytes after compressed
	MaxMsgBytes   int
	LogMaxLineLen int // default no limit (can log a very large message)

	// Logger from application
	Logger *log.Entry
}

// Producer _
type Producer struct {
	samProducer sarama.AsyncProducer
	conf        ProducerConfig
}

// NewProducer returns a connected Producer
func NewProducer(conf ProducerConfig) (*Producer, error) {
	conf.Logger.Infof("creating a producer with %#v", conf)
	// construct sarama config
	samConf := sarama.NewConfig()
	kafkaVersion, _ := sarama.ParseKafkaVersion("1.1.1")
	samConf.Version = kafkaVersion
	samConf.Producer.RequiredAcks = sarama.RequiredAcks(conf.RequiredAcks)
	samConf.Producer.Retry.Max = 5
	samConf.Producer.Retry.BackoffFunc = func(retries, maxRetries int) time.Duration {
		ret := 100 * time.Millisecond
		for retries > 0 {
			ret = 2 * ret
			retries--
		}
		return ret
	}
	samConf.Producer.Return.Successes = true
	if conf.IsCompressed {
		samConf.Producer.Compression = sarama.CompressionGZIP
		samConf.Producer.CompressionLevel = gzip.BestCompression
	}
	if conf.MaxMsgBytes > 0 {
		samConf.Producer.MaxMessageBytes = conf.MaxMsgBytes
	}

	// connect to kafka
	p := &Producer{conf: conf}
	brokers := strings.Split(conf.BrokersList, ",")
	var err error
	p.samProducer, err = sarama.NewAsyncProducer(brokers, samConf)
	if err != nil {
		return nil, fmt.Errorf("error create producer: %v", err)
	}
	conf.Logger.Infof("connected to kafka cluster %v", conf.BrokersList)
	go func() {
		for err := range p.samProducer.Errors() {
			errMsg := err.Err.Error()
			if errMsg == "circuit breaker is open" {
				errMsg = "probably you did not input a topic"
			}
			conf.Logger.Infof("failed to produce msgId %v to topic %v: %v",
				err.Msg.Metadata, err.Msg.Topic, errMsg)
		}
	}()
	go func() {
		for sent := range p.samProducer.Successes() {
			conf.Logger.Debugf(
				"delivered msgId %v to topic %v:%v:%v",
				sent.Metadata, sent.Topic, sent.Partition, sent.Offset)
		}
	}()
	return p, nil
}

func truncateStr(s string, limit int) string {
	if len(s) <= limit {
		return s
	}
	return s[:limit]
}

// SendExplicitMessage _
// Deprecated: use ProduceWithKey instead
func (p Producer) SendExplicitMessage(topic string, value []byte, key string) error {
	msgMeta := MsgMetadata{UniqueId: uuid.New(), SentAt: time.Now()}
	samMsg := &sarama.ProducerMessage{
		Value:     sarama.ByteEncoder(value),
		Topic:     topic,
		Metadata:  msgMeta,
		Timestamp: time.Now(),
	}
	if key != "" {
		samMsg.Key = sarama.StringEncoder(key)
	}
	var err error
	select {
	case p.samProducer.Input() <- samMsg:
		if p.conf.LogMaxLineLen > 0 {
			p.conf.Logger.Debugf(
				"producing msgId %v to %v:%v: len %v, msg: %v",
				msgMeta.UniqueId, samMsg.Topic, key,
				len(string(value)), truncateStr(string(value), p.conf.LogMaxLineLen))
		} else {
			p.conf.Logger.Debugf(
				"producing msgId %v to %v:%v: msg: %v",
				msgMeta.UniqueId, samMsg.Topic, key, string(value))
		}
		err = nil
	case <-time.After(1 * time.Minute):
		err = ErrWriteTimeout
	}
	return err
}

// ProduceJSON do JSON the object then sends JSONed string to Kafka clusters,
// in most cases you only need this func
func (p Producer) ProduceJSON(topic string, object interface{}) error {
	return p.ProduceJSONWithKey(topic, object, "")
}

// ProduceJSON do JSON the object then sends JSONed string to Kafka clusters,
// messages have the same key will be sent to the same partition
func (p Producer) ProduceJSONWithKey(
	topic string, object interface{}, kafkaKey string) error {
	switch v := object.(type) {
	case string:
		return p.ProduceWithKey(topic, []byte(v), kafkaKey)
	case []byte:
		return p.ProduceWithKey(topic, v, kafkaKey)
	default:
		beauty, err := json.Marshal(v)
		if err != nil {
			return err
		}
		return p.ProduceWithKey(topic, beauty, kafkaKey)
	}
}

// Produce sends input message to Kafka clusters.
// This func only return timeout error, other errors will be log by the Producer
func (p Producer) Produce(topic string, msg []byte) error {
	return p.SendExplicitMessage(topic, msg, "")
}

// ProduceWithKey sends messages have a same key to same partition.
func (p Producer) ProduceWithKey(topic string, message []byte, key string) error {
	return p.SendExplicitMessage(topic, message, key)
}

// Deprecated: use Produce instead
func (p Producer) SendMessage(topic string, msg []byte) error {
	return p.SendExplicitMessage(topic, msg, "")
}

// Errors when produce
var (
	ErrWriteTimeout = errors.New("write message timeout")
)

// SendMsgReliabilityLevel is the level of acknowledgement reliability.
// * NoResponse: highest throughput,
// * WaitForLocal: high, but not maximum durability and high but not maximum throughput,
// * WaitForAll: no data loss,
type SendMsgReliabilityLevel sarama.RequiredAcks

// SendMsgReliabilityLevel enum
const (
	NoResponse   = SendMsgReliabilityLevel(sarama.NoResponse)
	WaitForLocal = SendMsgReliabilityLevel(sarama.WaitForLocal)
	WaitForAll   = SendMsgReliabilityLevel(sarama.WaitForAll)
)

type MsgMetadata struct {
	UniqueId uuid.UUID
	SentAt   time.Time
}

func since(msgMetaI interface{}) time.Duration {
	msgMeta, ok := msgMetaI.(MsgMetadata)
	if !ok { // unreachable
		return 0
	}
	return time.Since(msgMeta.SentAt)
}
