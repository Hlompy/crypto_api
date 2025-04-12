package handlers

import (
	"crypto_api/internal/storage"
	"crypto_api/kafka"
	"encoding/json"
	"log"

	kafkaGo "github.com/segmentio/kafka-go"
)

func PublishUserCoinsToKafka(writer *kafkaGo.Writer) bool {
	coins, err := storage.GetUnpublishedUserCoins()
	if err != nil {
		log.Printf("❌ Ошибка при получении монет из БД: %v", err)
		return true // true → продолжаем тикать
	}

	if len(coins) == 0 {
		log.Println("🛑 Нет новых монет для публикации. Останавливаем публикацию.")
		return false // false → можно остановить тикер
	}

	log.Printf("🪙 Найдено монет для публикации: %d", len(coins))

	for _, coin := range coins {
		log.Printf("➡️ Публикуем монету: %+v", coin)

		data, err := json.Marshal(coin)
		if err != nil {
			log.Printf("❌ Не удалось сериализовать монету: %v", err)
			continue
		}

		err = kafka.ProduceMessage(writer, []byte(coin.Symbol), data)
		if err != nil {
			log.Printf("❌ Не удалось отправить монету в Kafka: %v", err)
			continue
		}

		// Обновим published_to_kafka = true
		err = storage.MarkCoinAsPublished(coin.ID)
		if err != nil {
			log.Printf("❌ Не удалось обновить статус монеты: %v", err)
			continue
		}

		log.Printf("📤 Отправлено сообщение в Kafka: %s", string(data))
	}

	return true
}
