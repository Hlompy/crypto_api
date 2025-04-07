// internal/handlers/coins.go
package handlers

import (
	"crypto_api/internal/models"
	"crypto_api/internal/services"
	"crypto_api/internal/storage"
	"crypto_api/pkg/cache"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// глобальный кэш с 10-минутным TTL
var coinCache = cache.NewCache(10 * time.Minute)
var coin models.UserCoin

func GetCoinsListHandler(c *gin.Context) {
	var coins []models.UserCoin

	// 1. Получаем данные из кэша
	data, found := coinCache.Get("top_coins")
	if !found {
		// 2. Если в кэше ничего нет — получаем с CoinGecko
		apiCoins := services.FetchTopCoins()

		// Преобразуем []services.Coin в []models.UserCoin
		for _, coin := range apiCoins {
			coins = append(coins, models.UserCoin{
				Name:       coin.Name,
				Symbol:     coin.Symbol,
				UsdPrice:   coin.Price,
				IsUserCoin: false, // Указываем, что это монета из API
			})
		}

		// Сохраняем в кэш
		coinCache.Set("top_coins", coins)
	} else {
		coins = data.([]models.UserCoin)
	}

	// 3. Возвращаем только монеты из API
	c.JSON(http.StatusOK, coins)
}

func GetCoinByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	coin, err := storage.GetUserCoinByID(id)
	if err != nil {
		log.Println("Error fetching user coin by ID:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "coin not found"})
		return
	}

	c.JSON(http.StatusOK, coin)
}

func GetCoinBySymbol(c *gin.Context) {
	symbol := strings.ToLower(c.Param("symbol"))

	// Получаем из кеша
	data, found := coinCache.Get("top_coins")
	if !found {
		apiCoins := services.FetchTopCoins()
		var coins []models.UserCoin
		for _, moneta := range apiCoins {
			coins = append(coins, models.UserCoin{
				Name:     moneta.Name,
				Symbol:   moneta.Symbol,
				UsdPrice: moneta.Price,
			})
		}
		coinCache.Set("top_coins", coins)
		data = coins
	}

	coins := data.([]models.UserCoin)

	for _, coin := range coins {
		if strings.ToLower(coin.Symbol) == symbol {
			c.JSON(http.StatusOK, coin)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "coin not found"})
}

func GetUserCoinsHandler(c *gin.Context) {
	// Получаем монеты, созданные пользователем
	userCoins, err := storage.GetUserCoinsByUserFlag(true) // Флаг true для пользовательских монет
	if err != nil {
		log.Println("Error fetching user coins:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user coins"})
		return
	}

	// Возвращаем только монеты пользователя
	c.JSON(http.StatusOK, userCoins)
}

func CreateCoinHandler(c *gin.Context) {
	if err := c.ShouldBindJSON(&coin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := storage.InsertUserCoin(coin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save"})
		return
	}
	c.JSON(http.StatusCreated, coin)
}

func UpdateCoinHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	// Прочитаем JSON тело
	if err := c.ShouldBindJSON(&coin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Обновим монету по ID
	err = storage.UpdateUserCoin(id, coin)
	if err != nil {
		log.Println("Error updating user coin:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update"})
		return
	}

	// Очистим кэш после обновления
	coinCache.Delete("top_coins")

	c.JSON(http.StatusOK, coin)
}

func DeleteCoinHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	// Удалим монету по ID
	err = storage.DeleteUserCoin(id)
	if err != nil {
		log.Println("Error deleting user coin:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete"})
		return
	}

	// Очистим кэш после удаления
	coinCache.Delete("top_coins")

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
