package db

import (
	"context"
	"fmt"

	"prvbot/internal/models"
	"github.com/jackc/pgx/v5"
)

func GetOrCreatePurchaseChannels(ctx context.Context, chatID int64) (*models.PurchaseChannels, error) {
	var pc models.PurchaseChannels

	err := Pool.QueryRow(ctx, `
		SELECT id, chat_id, channel1, channel2, channel3
		FROM purchase_channels
		WHERE chat_id = $1
	`, chatID).Scan(&pc.ID, &pc.ChatID, &pc.Channel1, &pc.Channel2, &pc.Channel3)

	if err == nil {
		return &pc, nil
	}

	err = Pool.QueryRow(ctx, `
	INSERT INTO purchase_channels (chat_id, channel1, channel2, channel3)
	VALUES ($1, '', '', '')
	RETURNING id, chat_id, channel1, channel2, channel3
	`, chatID).Scan(&pc.ID, &pc.ChatID, &pc.Channel1, &pc.Channel2, &pc.Channel3)

	if err != nil {
		return nil, fmt.Errorf("❌ не удалось создать purchase_channels для chat_id %d: %w", chatID, err)
	}

	return &pc, nil
}

func UpdateChannelField(chatID int64, field string, value string) error {
	column := ""
	switch field {
	case "channel1", "channel2", "channel3":
		column = field
	default:
		return fmt.Errorf("❌ Недопустимое поле канала: %s", field)
	}

	query := fmt.Sprintf("UPDATE purchase_channels SET %s = $1 WHERE chat_id = $2", column)
	_, err := Pool.Exec(context.Background(), query, value, chatID)
	if err != nil {
		return fmt.Errorf("⚠️ Ошибка обновления %s для chat_id %d: %w", column, chatID, err)
	}

	return nil
}

func GetPurchaseChannelsByChatID(ctx context.Context, chatID int64) (*models.PurchaseChannels, error) {
	query := `
		SELECT id, chat_id, channel1, channel2, channel3
		FROM purchase_channels
		WHERE chat_id = $1
		LIMIT 1;
	`

	row := Pool.QueryRow(ctx, query, chatID)

	var channels models.PurchaseChannels
	err := row.Scan(
		&channels.ID,
		&channels.ChatID,
		&channels.Channel1,
		&channels.Channel2,
		&channels.Channel3,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &channels, nil
}
