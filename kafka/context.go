package kafka

import (
	"github.com/margostino/anfield/configuration"
	"github.com/segmentio/kafka-go"
	"log"
)

var kafkaConnection *kafka.Conn
var kafkaReader *kafka.Reader
var kafkaWriter *kafka.Writer

func Initialize() {
	kafkaWriter = NewKafkaWriter()
	kafkaReader = NewKafkaReader()
}

func Close() {
	if err := kafkaReader.Close(); err != nil {
		log.Fatal("failed to close kafka reader:", err)
	}
	if err := kafkaWriter.Close(); err != nil {
		log.Fatal("failed to close kafka writer:", err)
	}
}

func KafkaWriter() *kafka.Writer {
	return kafkaWriter
}

func KafkaReader() *kafka.Reader {
	return kafkaReader
}

func NewKafkaWriter() *kafka.Writer {
	topic := configuration.Kafka().Topic
	address := configuration.Kafka().Address

	// make a writer that produces to topic-A, using the least-bytes distribution
	writer := &kafka.Writer{
		Addr:     kafka.TCP(address),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	//conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	return writer
}

func NewKafkaReader() *kafka.Reader {
	topic := configuration.Kafka().Topic
	address := configuration.Kafka().Address
	consumerGroupId := configuration.Kafka().ConsumerGroupId

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{address},
		GroupID: consumerGroupId,
		Topic:   topic,
		//MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	//conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	return reader
}

