package db

import (
	"context"
	"log"
	"fmt"
)

func CreatePurchaseChannels(ctx context.Context, chatID int64, totalAmount int) error {
	_, err := Pool.Exec(ctx, `
		INSERT INTO purchase_channels (chat_id, channel1, channel2, channel3, total_amount)
		VALUES ($1, '', '', '', $2)
	`, chatID, totalAmount)

	if err != nil {
		return fmt.Errorf("ошибка при создании purchase_channels для chat_id %d: %w", chatID, err)
	}

	log.Printf("🛍️ Добавлена запись в purchase_channels для chat_id %d (total_amount: %d)", chatID, totalAmount)
	return nil
}