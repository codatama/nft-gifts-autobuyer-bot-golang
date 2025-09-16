package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func HasAdminAccess(ctx context.Context, chatID int64) (bool, error) {
	var isAdmin bool
	query := `
		SELECT admin_access 
		FROM permissions 
		WHERE chat_id = $1
	`
	err := Pool.QueryRow(ctx, query, chatID).Scan(&isAdmin)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("ошибка при проверке прав: %w", err)
	}
	return isAdmin, nil
}

func HasRefundAccess(ctx context.Context, chatID int64) (bool, error) {
	var isAdmin bool
	query := `
		SELECT refund_access 
		FROM permissions 
		WHERE chat_id = $1
	`
	err := Pool.QueryRow(ctx, query, chatID).Scan(&isAdmin)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("ошибка при проверке прав: %w", err)
	}
	return isAdmin, nil
}

func GetSupportAdmins(ctx context.Context) ([]int64, error) {
	rows, err := Pool.Query(ctx, `SELECT chat_id FROM permissions WHERE technical_support = true`)
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе technical_support: %w", err)
	}
	defer rows.Close()

	var chatIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("ошибка при чтении chat_id: %w", err)
		}
		chatIDs = append(chatIDs, id)
	}
	return chatIDs, nil
}