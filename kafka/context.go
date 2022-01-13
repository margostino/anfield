package kafka

import (
	"github.com/margostino/anfield/configuration"
	"github.com/segmentio/kafka-go"
	"log"
)

var kafkaConnection *kafka.Conn
var kafkaReader *kafka.Reader
var kafkaWriter *kafka.Writer

//func Initialize() {
//	kafkaWriter = NewKafkaWriter()
//	kafkaReader = NewKafkaReader()
//}

type KafkaParams struct {
	topic           string
	address         string
	consumerGroupId string
}

func Close() {
	if kafkaReader != nil {
		closeReader()
	}

	if kafkaWriter != nil {
		closeWriter()
	}
}

func closeReader() {
	if err := kafkaReader.Close(); err != nil {
		log.Fatal("failed to close kafka reader:", err)
	}
}

func closeWriter() {
	if err := kafkaWriter.Close(); err != nil {
		log.Fatal("failed to close kafka writer:", err)
	}
}

func NewKafkaParams(configuration *configuration.Configuration) *KafkaParams {
	return &KafkaParams{
		address:         configuration.Kafka.Address,
		topic:           configuration.Kafka.Topic,
		consumerGroupId: configuration.Kafka.ConsumerGroupId,
	}
}

func NewWriter(params *KafkaParams) *kafka.Writer {
	//conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	// make a writer that produces to topic-A, using the least-bytes distribution
	return &kafka.Writer{
		Addr:     kafka.TCP(params.address),
		Topic:    params.topic,
		Balancer: &kafka.RoundRobin{},
	}
}

func NewReader(params *KafkaParams) *kafka.Reader {
	//conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{params.address},
		GroupID: params.consumerGroupId,
		Topic:   params.topic,
		//MinBytes: 10e3, // 10KB
		//MaxBytes: 10e6, // 10MB
	})
}
