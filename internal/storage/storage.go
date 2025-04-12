package storage

import (
	"crypto_api/internal/models"
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

// InitDB initializes the database connection
func InitDB() {
	dsn := os.Getenv("POSTGRES_DSN")
	var err error

	// Retry mechanism for database connection
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Printf("failed to open db, attempt %d: %v", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}

		if err = db.Ping(); err != nil {
			log.Printf("failed to connect to db, attempt %d: %v", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}

		log.Println("Successfully connected to the database")

		// 🧱 Создание таблицы, если не существует
		createTable := `
			CREATE TABLE IF NOT EXISTS user_coins (
				id SERIAL PRIMARY KEY,
				name TEXT NOT NULL,
				symbol TEXT UNIQUE NOT NULL,
				usd_price NUMERIC NOT NULL,
				is_user_coin BOOLEAN DEFAULT true,
				published_to_kafka BOOLEAN DEFAULT false
			)
		`
		_, err = db.Exec(createTable)
		if err != nil {
			log.Fatalf("failed to create table: %v", err)
		}

		return
	}

	log.Fatalf("failed to connect to db after multiple attempts")
}

// Close closes the database connection
func Close() {
	if db != nil {
		db.Close()
	}
}

func GetUserCoinByID(id int) (models.UserCoin, error) {
	var coin models.UserCoin
	err := db.QueryRow("SELECT id, name, symbol, usd_price FROM user_coins WHERE id = $1", id).
		Scan(&coin.ID, &coin.Name, &coin.Symbol, &coin.UsdPrice)

	if err != nil {
		return models.UserCoin{}, err
	}

	return coin, nil
}

func GetUserCoinsByUserFlag(isUserCoin bool) ([]models.UserCoin, error) {
	rows, err := db.Query("SELECT id, name, symbol, usd_price, is_user_coin FROM user_coins WHERE is_user_coin = $1", isUserCoin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var coins []models.UserCoin
	for rows.Next() {
		var coin models.UserCoin
		if err := rows.Scan(&coin.ID, &coin.Name, &coin.Symbol, &coin.UsdPrice, &coin.IsUserCoin); err != nil {
			return nil, err
		}
		coins = append(coins, coin)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return coins, nil
}

func GetAllUserCoins() ([]models.UserCoin, error) {
	rows, err := db.Query("SELECT id, name, symbol, usd_price, is_user_coin FROM user_coins")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var coins []models.UserCoin
	for rows.Next() {
		var coin models.UserCoin
		if err := rows.Scan(&coin.ID, &coin.Name, &coin.Symbol, &coin.UsdPrice, &coin.IsUserCoin); err != nil {
			return nil, err
		}
		coins = append(coins, coin)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return coins, nil
}

func GetUnpublishedUserCoins() ([]models.UserCoin, error) {
	rows, err := db.Query(`
		SELECT id, name, symbol, usd_price, is_user_coin 
		FROM user_coins 
		WHERE is_user_coin = true AND published_to_kafka = false
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var coins []models.UserCoin
	for rows.Next() {
		var coin models.UserCoin
		if err := rows.Scan(&coin.ID, &coin.Name, &coin.Symbol, &coin.UsdPrice, &coin.IsUserCoin); err != nil {
			return nil, err
		}
		coins = append(coins, coin)
	}

	return coins, rows.Err()
}

func MarkCoinAsPublished(id int) error {
	_, err := db.Exec("UPDATE user_coins SET published_to_kafka = true WHERE id = $1", id)
	return err
}

func InsertUserCoin(c models.UserCoin) error {
	// Устанавливаем флаг is_user_coin в true
	_, err := db.Exec("INSERT INTO user_coins (name, symbol, usd_price, is_user_coin) VALUES ($1, $2, $3, $4)", c.Name, c.Symbol, c.UsdPrice, true)
	if err != nil {
		log.Println("Error inserting user coin:", err)
		return err
	}
	return nil
}

func UpdateUserCoin(id int, c models.UserCoin) error {
	_, err := db.Exec("UPDATE user_coins SET name=$1, usd_price=$2 WHERE id=$3", c.Name, c.UsdPrice, id)
	return err
}

func DeleteUserCoin(id int) error {
	_, err := db.Exec("DELETE FROM user_coins WHERE id=$1", id)
	return err
}
