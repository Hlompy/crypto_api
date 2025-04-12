package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// NewKafkaWriter создает новый writer для Kafka.
func NewKafkaWriter(brokers []string, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

// ProduceMessage отправляет сообщение в Kafka.
func ProduceMessage(writer *kafka.Writer, key, value []byte) error {
	err := writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   key,
			Value: value,
			Time:  time.Now(),
		},
	)
	if err != nil {
		log.Printf("Не удалось отправить сообщение: %v", err)
		return err
	}
	log.Println("Сообщение отправлено успешно")
	return nil
}

// EnsureTopic создает топик с заданной конфигурацией, если он отсутствует.
func EnsureTopic(brokerAddress string, topic string) {
	conn, err := kafka.Dial("tcp", brokerAddress)
	if err != nil {
		log.Fatal("Не удалось подключиться к Kafka:", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		log.Fatal("Не удалось получить контроллер Kafka:", err)
	}

	// Формируем адрес контроллера корректно
	controllerAddr := fmt.Sprintf("%s:%d", controller.Host, controller.Port)
	controllerConn, err := kafka.Dial("tcp", controllerAddr)
	if err != nil {
		log.Fatal("Не удалось подключиться к контроллеру Kafka:", err)
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		// Если топик уже существует, kafka-go вернет ошибку, но это можно проигнорировать.
		log.Printf("Ошибка при создании топика (вероятно, он уже существует): %v", err)
	} else {
		log.Printf("Топик '%s' создан успешно", topic)
	}
}

// WaitForKafka проверяет, доступен ли брокер Kafka (устанавливая TCP-соединение).
func WaitForKafka(brokers []string, retries int, delay time.Duration) {
	for i := 0; i < retries; i++ {
		conn, err := kafka.Dial("tcp", brokers[0])
		if err == nil {
			conn.Close()
			log.Println("✅ Kafka доступна!")
			return
		}
		log.Printf("⚠️ Kafka еще не доступна (попытка %d/%d): %v", i+1, retries, err)
		time.Sleep(delay)
	}
	log.Fatal("❌ Не удалось подключиться к Kafka после нескольких попыток")
}
