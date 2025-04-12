package main

import (
	"crypto_api/internal/handlers"
	"crypto_api/internal/storage"
	"crypto_api/kafka"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	brokers := []string{"kafka:9092"}
	topic := "user-coins"
	groupID := "crypto-api-group"

	kafka.WaitForKafka(brokers, 10, 2*time.Second)
	kafka.EnsureTopic(brokers[0], topic)
	storage.InitDB()

	// Запуск Kafka Consumer в отдельной горутине
	go kafka.StartConsumer(brokers, topic, groupID)

	writer := kafka.NewKafkaWriter(brokers, topic)

	// ⏱ Ticker будет каждые 10 секунд проверять базу и слать в Kafka
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("🔁 Проверка базы и публикация монет в Kafka")

		continuePublishing := handlers.PublishUserCoinsToKafka(writer)
		if !continuePublishing {
			log.Println("✅ Все монеты опубликованы. Завершаем тикер.")
			break
		}
	}

	// Инициализация роутов
	router := gin.Default()
	router.GET("/coins", handlers.GetCoinsListHandler)
	router.GET("/coins/symbol/:symbol", handlers.GetCoinBySymbol)
	router.GET("/coins/:id", handlers.GetCoinByID)
	router.GET("/coins/users", handlers.GetUserCoinsHandler)
	router.POST("/coins", handlers.CreateCoinHandler)
	router.PATCH("/coins/:id", handlers.UpdateCoinHandler)
	router.DELETE("/coins/:id", handlers.DeleteCoinHandler)

	// Запуск сервера
	router.Run(":8000")
}
