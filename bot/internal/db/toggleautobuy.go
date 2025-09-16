package db

import (
	"context"
)

func ToggleAutoBuy(ctx context.Context, chatID int64) error {
	_, err := Pool.Exec(ctx, `
		UPDATE users
		SET auto_buy = NOT auto_buy
		WHERE chat_id = $1
	`, chatID)
	return err
}