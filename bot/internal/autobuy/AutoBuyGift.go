package autobuy

import (
	"gopkg.in/telebot.v4"

	"context"
	"fmt"
	"log"
	
	"prvbot/internal/db"
	"prvbot/internal/tgapi"
)

func AutoBuyGift(ctx context.Context, bot *telebot.Bot, chatID int64, recipient string, balance int, gift tgapi.TelegramGift) error {
	if balance < gift.StarCount {
		return fmt.Errorf("недостаточно звёзд: нужно %d⭐️, у пользователя %d⭐️", gift.StarCount, balance)
	}

	err := SendAutoGiftToUser(bot, recipient, gift.ID)
	if err != nil {
		log.Printf("[AUTO-BUY] ❌ Не удалось отправить подарок %s → %s: %v", gift.ID, recipient, err)
		return fmt.Errorf("не удалось отправить подарок: %w", err)
	}

	if err := db.DecreaseBalance(ctx, chatID, gift.StarCount); err != nil {
		return fmt.Errorf("не удалось списать %d⭐️: %w", gift.StarCount, err)
	}

	log.Printf("[AUTO-BUY] ✅ %d → %s: подарок %s за %d⭐️", chatID, recipient, gift.ID, gift.StarCount)
	return nil
}