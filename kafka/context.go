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

func NewWriter() {
	topic := configuration.Kafka().Topic
	address := configuration.Kafka().Address

	// make a writer that produces to topic-A, using the least-bytes distribution
	writer := &kafka.Writer{
		Addr:     kafka.TCP(address),
		Topic:    topic,
		Balancer: &kafka.RoundRobin{},
	}

	//conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	kafkaWriter = writer
}

func NewReader(consumerGroupId string) {
	topic := configuration.Kafka().Topic
	address := configuration.Kafka().Address
	//consumerGroupId := configuration.Kafka().ConsumerGroupId

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{address},
		GroupID: consumerGroupId,
		Topic:   topic,
		//MinBytes: 10e3, // 10KB
		//MaxBytes: 10e6, // 10MB
	})

	//conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	kafkaReader = reader
}
