package gifting

import (
	"gopkg.in/telebot.v4"

	"context"
	"fmt"
	"log"

	"prvbot/internal/db"
	"prvbot/internal/messages"
	"prvbot/internal/tgapi"
)

func ExecuteGiftPurchaseAndSend(c telebot.Context, gift *tgapi.TelegramGift, balance int, target interface{}) error {
	userID := c.Sender().ID
	bot := c.Bot().(*telebot.Bot)

	err := tgapi.SendGiftToUser(bot, target, gift.ID)
	if err != nil {
		fmt.Println("❌ SendGift ошибка:", err)
		return c.Send(messages.ErrGiftSendFailed)
	}

	err = db.DecreaseBalance(context.Background(), userID, gift.StarCount)
	if err != nil {
		log.Printf("❗ Ошибка обновления баланса: %v\n", err)
		return c.Send(messages.ErrBalanceFetchFailed)
	}

	newBalance := balance - gift.StarCount
	return c.Send(fmt.Sprintf("🎉 Подарок %s успешно отправлен!\n💰 Остаток: %d ⭐️", gift.Sticker.Emoji, newBalance))
}