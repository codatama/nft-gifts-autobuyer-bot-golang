package db

import (
	"context"
	"log"
	"fmt"

	"prvbot/internal/models"
)

func InsertTransaction(ctx context.Context, transaction models.StarTransaction) error {
	query := `
		INSERT INTO transactions (id, amount, nanostar_amount, date, source, receiver)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id, date) DO NOTHING;
	`
	_, err := Pool.Exec(ctx, query,
		transaction.ID,
		transaction.Amount,
		transaction.NanostarAmount,
		transaction.Date,
		transaction.Source.User.ID,
		transaction.Receiver.User.ID,
	)

	return err
}

func InsertTransactions(ctx context.Context, txs []models.StarTransaction) error {
	for _, t := range txs {
		err := InsertTransaction(ctx, t)
		if err != nil {
			log.Printf("❌ Не удалось вставить транзакцию %s: %v", t.ID, err)
		}
	}
	return nil
}

func GetTransactionsByChargeAndUser(ctx context.Context, userID int64, chargeID string) ([]models.StarTransaction, error) {
	rows, err := Pool.Query(ctx, `
		SELECT id, amount, source
		FROM transactions
		WHERE source = $1 AND id = $2
	`, userID, chargeID)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к базе: %w", err)
	}
	defer rows.Close()

	var result []models.StarTransaction
	for rows.Next() {
		var txn models.StarTransaction
		var sourceID int64

		if err := rows.Scan(&txn.ID, &txn.Amount, &sourceID); err != nil {
			return nil, fmt.Errorf("ошибка scan строки: %w", err)
		}

		txn.Source = models.TransactionPartnerUser{
			Type:            "user",
			TransactionType: "payment",
			User: models.TelegramUser{
				ID:       sourceID,
				IsBot:    false,
				Username: "",
			},
		}

		result = append(result, txn)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка итерации по строкам: %w", err)
	}

	return result, nil
}