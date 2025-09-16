package autobuy

import (
	"gopkg.in/telebot.v4"
	
	"fmt"
	"log"
)

func SendAutoGiftToUser(bot *telebot.Bot, recipient string, giftID string) error {
	params := map[string]interface{}{
		"chat_id": recipient,
		"gift_id": giftID,
	}

	resp, err := bot.Raw("sendGift", params)
	if err != nil {
		return fmt.Errorf("ошибка при вызове sendGift: %w", err)
	}

	log.Printf("📤 sendGift: %s", string(resp))
	return nil
}