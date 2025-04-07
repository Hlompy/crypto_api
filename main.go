package main

import (
	"crypto_api/config"
	"crypto_api/internal/handlers"
	"crypto_api/internal/storage"

	// Assuming this is where the db connection is initialized
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Загрузка .env
	config.LoadEnv()

	// 2. Подключение к PostgreSQL
	storage.InitDB() // Assuming storage.Init() initializes and returns the DB connection
	defer storage.Close()

	// 3. Кэш
	//cacheInstance := cache.NewCache(10 * time.Minute) // Renamed variable to avoid conflict

	// 4. Инициализация роутов
	router := gin.Default()
	router.GET("/coins", handlers.GetCoinsListHandler) // Монеты из API
	router.GET("/coins/symbol/:symbol", handlers.GetCoinBySymbol)
	router.GET("/coins/:id", handlers.GetCoinByID)           // По ID (из базы)
	router.GET("/coins/users", handlers.GetUserCoinsHandler) // Монеты пользователя
	router.POST("/coins", handlers.CreateCoinHandler)
	router.PATCH("/coins/:id", handlers.UpdateCoinHandler)
	router.DELETE("/coins/:id", handlers.DeleteCoinHandler)

	// 5. Запуск
	router.Run(":8000")
}
