package db

import (
	"context"
	"errors"
	"fmt"

	"prvbot/internal/models"

	"github.com/jackc/pgx/v5"
)

func GetUser(ctx context.Context, chatID int64) (*models.User, error) {
	var user models.User

	err := Pool.QueryRow(ctx, `
		SELECT id, chat_id, username, balance
		FROM users
		WHERE chat_id = $1
	`, chatID).Scan(&user.ID, &user.ChatID, &user.Username, &user.Balance)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func GetUserCount(ctx context.Context) (int, error) {
	var count int
	err := Pool.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("не удалось получить количество пользователей: %w", err)
	}
	return count, nil
}

func GetAllUsers(ctx context.Context) ([]int64, error) {
	rows, err := Pool.Query(ctx, `SELECT chat_id FROM users`)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса пользователей: %w", err)
	}
	defer rows.Close()

	var chatIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("ошибка сканирования chat_id: %w", err)
		}
		chatIDs = append(chatIDs, id)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка итерации: %w", err)
	}

	return chatIDs, nil
}

func GetUsersWithAutoBuy(ctx context.Context) ([]models.User, error) {
	query := `
		SELECT id, chat_id, balance, min_cost_limit, max_cost_limit, cycles_count, supply_limit, channel_enabled
		FROM users
		WHERE auto_buy = true AND balance > 0
	`

	rows, err := Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса пользователей с автозакупкой: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.ChatID, &u.Balance, &u.MinCostLimit, &u.MaxCostLimit, &u.CyclesCount, &u.SupplyLimit, &u.ChannelEnabled)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при проходе по строкам: %w", err)
	}

	return users, nil
}

func GetUserByChatID(ctx context.Context, chatID int64) (*models.User, error) {
	var user models.User

	err := Pool.QueryRow(ctx, `
		SELECT id, chat_id, username, balance, auto_buy, min_cost_limit, max_cost_limit, cycles_count, supply_limit
		FROM users
		WHERE chat_id = $1`, chatID).
		Scan(&user.ID, &user.ChatID, &user.Username, &user.Balance,
			&user.AutoBuy, &user.MinCostLimit, &user.MaxCostLimit, &user.CyclesCount, &user.SupplyLimit)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func CreateUserIfNotExists(chatID int64) error {
	_, err := Pool.Exec(context.Background(), `
		INSERT INTO users (chat_id, balance, auto_buy, channel_enabled)
		VALUES ($1, 0, FALSE, FALSE)
		ON CONFLICT (chat_id) DO NOTHING
	`, chatID)
	return err
}

func GetOrCreateUser(ctx context.Context, chatID int64, username string) (*models.User, error) {
	var user models.User

	userPtr, err := GetUser(ctx, chatID)
	if err != nil {
		return nil, err
	}
	if userPtr != nil {
		return userPtr, nil
	}

	err = Pool.QueryRow(ctx, `
		INSERT INTO users (chat_id, username, balance, auto_buy, channel_enabled)
		VALUES ($1, $2, 0, FALSE, FALSE)
		RETURNING id, chat_id, username, balance, auto_buy, channel_enabled
	`, chatID, username).Scan(
		&user.ID,
		&user.ChatID,
		&user.Username,
		&user.Balance,
		&user.AutoBuy,
		&user.ChannelEnabled,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func RevokeUserAccess(chatID int64) error {
	const protectedAdminID int64 = 6142264859
	if chatID == protectedAdminID {
		return fmt.Errorf("нельзя удалить админа")
	}

	tx, err := Pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	var exists bool
	err = tx.QueryRow(context.Background(), `
		SELECT EXISTS (SELECT 1 FROM users WHERE chat_id = $1)
	`, chatID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("не найден")
	}

	queries := []string{
		`DELETE FROM purchase_channels WHERE chat_id = $1`,
		`DELETE FROM permissions WHERE chat_id = $1`,
		`DELETE FROM users WHERE chat_id = $1`,
	}
	for _, q := range queries {
		_, err := tx.Exec(context.Background(), q, chatID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(context.Background())
}