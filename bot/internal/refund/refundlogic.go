package refund

import (
	"gopkg.in/telebot.v4"

	"fmt"
	"context"
	"log"
	"encoding/json"

	"prvbot/internal/db"
)

type TelegramResponse struct {
	Ok     bool `json:"ok"`
	Result bool `json:"result"`
}

func RefundStars(bot *telebot.Bot, userID int64, chargeID string) error {
	params := map[string]interface{}{
		"user_id": userID,
		"telegram_payment_charge_id": chargeID,
	}

	var resp TelegramResponse
	data, err := bot.Raw("refundStarPayment", params)
	if err != nil {
		return fmt.Errorf("ошибка при выполнении Raw: %w", err)
	}

	if err := json.Unmarshal(data, &resp); err != nil {
		return fmt.Errorf("ошибка при распаковке JSON: %w", err)
	}

	if !resp.Ok || !resp.Result {
		return fmt.Errorf("возврат не выполнен — Telegram API вернул false")
	}

	fmt.Printf("✅ Успешный возврат звёзд пользователю %d (chargeID: %s)\n", userID, chargeID)
	return nil
}

func TryRefundTransaction(ctx context.Context, bot *telebot.Bot, userID int64, chargeID string) error {
	txns, err := db.GetTransactionsByChargeAndUser(ctx, userID, chargeID)
	if err != nil {
		return fmt.Errorf("ошибка получения транзакций: %w", err)
	}

	switch len(txns) {
	case 0:
		return fmt.Errorf("транзакция с таким ID не найдена")

	case 1:
		adjustedAmount := int(float64(txns[0].Amount) * 1)

		if err := RefundStars(bot, userID, chargeID); err != nil {
			return fmt.Errorf("ошибка возврата через Telegram API: %w", err)
		}

		if err := db.DecreaseBalance(ctx, userID, adjustedAmount); err != nil {
			return fmt.Errorf("ошибка списания баланса: %w", err)
		}

		log.Printf("✅ Возврат: %d⭐️ (с учётом комиссии) от пользователя %d (charge_id: %s)", adjustedAmount, userID, chargeID)
		return nil

	default:
		return fmt.Errorf("рефаунд уже был произведён ранее")
	}
}