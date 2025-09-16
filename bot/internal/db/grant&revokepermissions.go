package db

import	(
	"context"
	"fmt"
)

func GrantPermission(ctx context.Context, chatID int64, column string) error {
	if column != "admin_access" && column != "refund_access" && column != "technical_support" {
		return fmt.Errorf("❌ недопустимая колонка: %s", column)
	}

	query := fmt.Sprintf(`UPDATE permissions SET %s = true WHERE chat_id = $1`, column)
	cmd, err := Pool.Exec(ctx, query, chatID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("такого chat_id нет в таблице permissions")
	}
	return nil
}

func RevokePermission(ctx context.Context, chatID int64, column string) error {
	const protectedAdminID int64 = 6142264859
	if chatID == protectedAdminID {
		return fmt.Errorf("нельзя отозвать права у администратора")
	}

	if column != "admin_access" && column != "refund_access" && column != "technical_support" {
		return fmt.Errorf("❌ недопустимая колонка: %s", column)
	}

	query := fmt.Sprintf(`UPDATE permissions SET %s = false WHERE chat_id = $1`, column)
	cmd, err := Pool.Exec(ctx, query, chatID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("chat_id %d не найден в permissions", chatID)
	}
	return nil
}