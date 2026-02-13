package kafka

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

var writer *kafka.Writer
var dlqWriter *kafka.Writer

// InitKafka initializes kafka connection
func InitKafka() error {
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	topic := os.Getenv("KAFKA_TOPIC")

	if kafkaBroker == "" || topic == "" {
		return nil // skip if env not set
	}

	writer = &kafka.Writer{
		Addr:                   kafka.TCP(kafkaBroker),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	// Retry logic to check Kafka availability
	for i := 0; i < 10; i++ {
		err := writer.WriteMessages(context.Background(), kafka.Message{
			Value: []byte("ping"),
		})
		if err == nil {
			log.Println("Kafka is ready")
			return nil
		}
		log.Println("Waiting for Kafka to be ready...")
		time.Sleep(2 * time.Second)
	}

	log.Println("Kafka producer initialized (topic will be auto-created on first write)")
	return nil
}

func initDLQWriter() {
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	topic := os.Getenv("KAFKA_TOPIC")
	if kafkaBroker == "" || topic == "" {
		return
	}
	dlqWriter = &kafka.Writer{
		Addr:                   kafka.TCP(kafkaBroker),
		Topic:                  topic + "-dlq",
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}
}

func PublishMessage(ctx context.Context, payload []byte) error {
	if writer == nil {
		log.Println("Kafka writer is nil; skipping message publish")
		return fmt.Errorf("kafka writer not initialized")
	}

	msg := kafka.Message{
		Value: payload,
	}

	return writer.WriteMessages(ctx, msg)
}

func PublishToDLQ(ctx context.Context, payload []byte) error {
	if dlqWriter == nil {
		initDLQWriter()
	}
	if dlqWriter == nil {
		log.Println("DLQ writer is nil; cannot publish to DLQ")
		return fmt.Errorf("dlq writer not initialized")
	}
	return dlqWriter.WriteMessages(ctx, kafka.Message{Value: payload})
}

func ConsumeMessage(ctx context.Context, handler func(kafka.Message) error) error {
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	topic := os.Getenv("KAFKA_TOPIC")

	if kafkaBroker == "" || topic == "" {
		log.Println("Kafka consumer config not set; skipping consumer start.")
		return nil
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaBroker},
		Topic:    topic,
		GroupID:  "rebalance-consumer",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	go func() {
		defer reader.Close()
		for {
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Kafka read error: %v\n", err)
				continue
			}

			if err := handler(msg); err != nil {
				log.Printf("Failed to process message: %v\n", err)
				continue
			}

			reader.CommitMessages(ctx, msg)
		}
	}()

	log.Println("Kafka consumer started")
	return nil
}
