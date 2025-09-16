package db

import (
	"context"
	"fmt"
	"log"
)

func CaptureGlobalBalanceSnapshot(ctx context.Context) error {
	var total int64
	err := Pool.QueryRow(ctx, `
		SELECT SUM(balance) FROM users
	`).Scan(&total)
	if err != nil {
		return fmt.Errorf("ошибка при суммировании: %w", err)
	}

	_, err = Pool.Exec(ctx, `
		INSERT INTO star_balance_log (balance) VALUES ($1)
	`, total)
	if err != nil {
		return fmt.Errorf("ошибка при сохранении логов: %w", err)
	}

	log.Printf("💾 Общий баланс сохранён: %d⭐️", total)
	return nil
}

func GetLatestBalanceDifference(ctx context.Context) (int64, error) {
	var diff int64
	err := Pool.QueryRow(ctx, `
		SELECT difference
		FROM star_balance_log
		ORDER BY recorded_at DESC
		LIMIT 1;
	`).Scan(&diff)
	if err != nil {
		return 0, fmt.Errorf("ошибка при получении difference: %w", err)
	}
	return diff, nil
}