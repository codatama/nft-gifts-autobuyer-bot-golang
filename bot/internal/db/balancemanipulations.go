package db

import (
	"context"
	"fmt"
)

var GlobalCommissionRate = 0.02

func UpdateBalance(ctx context.Context, chatID int64, delta int) error {
	commissionRate := GlobalCommissionRate

	commission := int(float64(delta) * commissionRate)
	netAmount := delta - commission

	tag, err := Pool.Exec(ctx, `
		UPDATE users
		SET balance = balance + $1
		WHERE chat_id = $2
	`, netAmount, chatID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении баланса: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("не найден пользователь с chat_id %d", chatID)
	}

	fmt.Printf("✅ Пользователь %d пополнил баланс: +%d⭐️ (−%d⭐️ комиссия @ %.2f%%)\n", chatID, netAmount, commission, commissionRate*100)
	return nil
}

func GetBalance(ctx context.Context, chatID int64) (int, error) {
	var balance int
	err := Pool.QueryRow(ctx, `
		SELECT balance FROM users WHERE chat_id = $1
	`, chatID).Scan(&balance)
	return balance, err
}

func GetTotalBalance(ctx context.Context) (int, error) {
	var total int
	err := Pool.QueryRow(ctx, "SELECT COALESCE(SUM(balance), 0) FROM users").Scan(&total)
	return total, err
}

func DecreaseBalance(ctx context.Context, chatID int64, amount int) error {
	query := `
		UPDATE users
		SET balance = balance - $1
		WHERE chat_id = $2 AND balance >= $1
	`

	cmdTag, err := Pool.Exec(ctx, query, amount, chatID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("недостаточно средств для списания %d звёзд", amount)
	}

	return nil
}

func SetGlobalCommissionRate(rate float64) error {
	if rate < 0.02 {
		return fmt.Errorf("комиссия не может быть ниже 0.02")
	}
	if rate > 1 {
		return fmt.Errorf("комиссия не может быть выше 1.0")
	}
	GlobalCommissionRate = rate
	return nil
}