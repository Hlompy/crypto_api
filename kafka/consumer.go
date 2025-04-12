package kafka

import (
	"context"
	"log"
	"time"

	kafkaGo "github.com/segmentio/kafka-go"
)

func StartConsumer(brokers []string, topic string, groupID string) {
	r := kafkaGo.NewReader(kafkaGo.ReaderConfig{
		Brokers:         brokers,
		GroupID:         groupID,
		Topic:           topic,
		MinBytes:        10e3, // 10KB
		MaxBytes:        10e6, // 10MB
		CommitInterval:  time.Second,
		ReadLagInterval: -1, // disables lag metrics for performance
	})
	defer r.Close()

	log.Println("👂 Kafka Consumer запущен и слушает топик:", topic)

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("❌ Ошибка при чтении сообщения из Kafka: %v", err)
			continue
		}

		log.Printf("✅ [CONSUMER] Получено сообщение из Kafka:\n🔹 Ключ: %s\n📦 Данные: %s", string(msg.Key), string(msg.Value))
	}
}
