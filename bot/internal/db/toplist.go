package db

import (
	"prvbot/internal/models"
	"context"
)

func GetTopUsersByBalance(ctx context.Context, limit int) ([]models.User, error) {
	rows, err := Pool.Query(ctx, `
		SELECT id, chat_id, balance, username
		FROM users
		WHERE balance > 0
		ORDER BY balance DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.ChatID, &u.Balance, &u.Username)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func GetUserBalanceRank(ctx context.Context, chatID int64) (int, error) {
	var rank int
	err := Pool.QueryRow(ctx, `
		SELECT COUNT(*) + 1
		FROM users
		WHERE balance > (
			SELECT balance FROM users WHERE chat_id = $1
		)
	`, chatID).Scan(&rank)

	if err != nil {
		return 0, err
	}
	return rank, nil
}