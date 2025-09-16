package tgapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"log"
	"time"
	"context"

	"prvbot/internal/models"
	"prvbot/internal/db"
	"prvbot/config"
)

type starTransactionsResponse struct {
	Ok     bool `json:"ok"`
	Result struct {
		Transactions []models.StarTransaction `json:"transactions"`
	} `json:"result"`
}

func GetStarTransactions(botToken string, offset, limit int) ([]models.StarTransaction, error) {
	if limit <= 0 || limit > 100 {
		limit = 100
	}
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getStarTransactions?offset=%d&limit=%d", botToken, offset, limit)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP error: %w", err)
	}
	defer resp.Body.Close()

	var parsed starTransactionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("JSON decode error: %w", err)
	}
	if !parsed.Ok {
		return nil, fmt.Errorf("telegram API returned failure")
	}
	return parsed.Result.Transactions, nil
}

func SyncTransactions(botToken string) error {
	offset := 0
	limit := 100

	for {
		txs, err := GetStarTransactions(botToken, offset, limit)
		if err != nil {
			return err
		}
		if len(txs) == 0 {
			break
		}

		err = db.InsertTransactions(context.Background(), txs)
		if err != nil {
			log.Printf("❌ Ошибка при записи транзакций: %v", err)
		}

		offset += len(txs)
	}

	return nil
}

func StartStarSyncScheduler() {
	go func() {
		botToken := config.Load().TelegramToken
		dbURL := config.Load().DatabaseURL

		if botToken == "" || dbURL == "" {
			log.Println("❌ Не хватает BOT_TOKEN или DATABASE_URL")
			return
		}

		for {
			log.Println("🌟 Запуск синхронизации транзакций...")
			if err := SyncTransactions(botToken); err != nil {
				log.Printf("❌ Ошибка при SyncTransactions: %v", err)
			} else {
				log.Println("✅ Транзакции успешно синхронизированы")
			}

			time.Sleep(12 * time.Hour)
		}
	}()
}