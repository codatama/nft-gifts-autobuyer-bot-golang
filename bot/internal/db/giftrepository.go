package db

import (
	"context"

	"prvbot/internal/models"
)

func GetGiftByID(id int) (*models.GiftEntity, error) {
	query := `
		SELECT id, api_id, star_count, emoji, file_id, thumbnail
		FROM telegram_gifts
		WHERE id = $1
	`

	var gift models.GiftEntity
	err := Pool.QueryRow(context.Background(), query, id).Scan(
		&gift.ID,
		&gift.ApiID,
		&gift.StarCount,
		&gift.Emoji,
		&gift.FileID,
		&gift.Thumbnail,
	)
	if err != nil {
		return nil, err
	}

	return &gift, nil
}

func UpsertTelegramGift(ctx context.Context, gift *models.GiftEntity) error {
	query := `
	INSERT INTO telegram_gifts (api_id, star_count, emoji, file_id, thumbnail)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (api_id) DO UPDATE SET
		star_count = EXCLUDED.star_count,
		emoji = EXCLUDED.emoji,
		file_id = EXCLUDED.file_id,
		thumbnail = EXCLUDED.thumbnail;
	`

	_, err := Pool.Exec(ctx, query,
		gift.ApiID,
		gift.StarCount,
		gift.Emoji,
		gift.FileID,
		gift.Thumbnail,
	)

	return err
}