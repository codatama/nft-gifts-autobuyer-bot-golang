package db

import (
	"context"
)

func UpdateMinLimit(ctx context.Context, chatID int64, limit int) error {
	_, err := Pool.Exec(ctx, `
		UPDATE users
		SET min_cost_limit = $1
		WHERE chat_id = $2
	`, limit, chatID)
	return err
}

func UpdateMaxLimit(ctx context.Context, chatID int64, limit int) error {
	_, err := Pool.Exec(ctx, `
		UPDATE users
		SET max_cost_limit = $1
		WHERE chat_id = $2
	`, limit, chatID)
	return err
}

func UpdateCyclesCount(ctx context.Context, chatID int64, cycles int) error {
	_, err := Pool.Exec(ctx, `
		UPDATE users
		SET cycles_count = $1
		WHERE chat_id = $2
	`, cycles, chatID)
	return err
}

func UpdateSupplyLimit(ctx context.Context, chatID int64, supply int) error {
	_, err := Pool.Exec(ctx, `
	UPDATE users
	SET supply_limit = $1
	WHERE chat_id = $2
	`, supply, chatID)
	return err
}