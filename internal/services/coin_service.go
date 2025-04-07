package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"crypto_api/pkg/cache"
)

type Coin struct {
	ID     string  `json:"id"`
	Symbol string  `json:"symbol"`
	Name   string  `json:"name"`
	Price  float64 `json:"current_price"`
}

var coinCache = cache.NewCache(10 * time.Minute)

func FetchTopCoins() []Coin {
	if data, found := coinCache.Get("top_coins"); found {
		return data.([]Coin)
	}

	url := "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=100&page=1&sparkline=false"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Failed to fetch from CoinGecko:", err)
		return nil
	}
	defer resp.Body.Close()

	var coins []Coin
	if err := json.NewDecoder(resp.Body).Decode(&coins); err != nil {
		fmt.Println("Failed to decode CoinGecko response:", err)
		return nil
	}

	coinCache.Set("top_coins", coins)
	return coins
}

func FetchCoinByID(id string) *Coin {
	key := fmt.Sprintf("coin_%s", id)
	if data, found := coinCache.Get(key); found {
		coin := data.(Coin)
		return &coin
	}

	url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s?localization=false&tickers=false&market_data=true&community_data=false&developer_data=false&sparkline=false", id)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Failed to fetch coin from CoinGecko:", err)
		return nil
	}
	defer resp.Body.Close()

	var result struct {
		ID         string `json:"id"`
		Symbol     string `json:"symbol"`
		Name       string `json:"name"`
		MarketData struct {
			CurrentPrice struct {
				USD float64 `json:"usd"`
			} `json:"current_price"`
		} `json:"market_data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Failed to decode coin by ID:", err)
		return nil
	}

	coin := Coin{
		ID:     result.ID,
		Symbol: result.Symbol,
		Name:   result.Name,
		Price:  result.MarketData.CurrentPrice.USD,
	}

	coinCache.Set(key, coin)
	return &coin
}
