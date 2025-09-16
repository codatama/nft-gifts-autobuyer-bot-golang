package adminpanel

import (
	"gopkg.in/telebot.v4"

	"time"
	"context"
	
	"prvbot/internal/db"
	"prvbot/internal/userstates"
)

func HandleGlobalBroadcastRequest(c telebot.Context) error {
	chatID := c.Chat().ID
	userstates.Set(chatID, "awaiting_global_broadcast")

	go func(chatID int64, bot *telebot.Bot) {
		time.Sleep(30 * time.Second)
		if userstates.Get(chatID) == "awaiting_global_broadcast" {
			userstates.Clear(chatID)
			bot.Send(&telebot.Chat{ID: chatID}, "⏳ Время на ввод истекло.")
		}
	}(chatID, c.Bot().(*telebot.Bot))

	return c.Send("📨 Введите сообщение (текст или фото с подписью), которое будет отправлено всем пользователям:",
		&telebot.SendOptions{ParseMode: telebot.ModeMarkdown})
}

func HandleRefundRequest(c telebot.Context) error {
	chatID := c.Chat().ID

	ok, err := db.HasRefundAccess(context.Background(), chatID)
	if err != nil || !ok {
		return c.Send("⛔ У вас нет прав на выполнение возврата.")
	}

	userstates.Set(chatID, "awaiting_refund_input")

	go func(chatID int64, bot *telebot.Bot) {
		time.Sleep(30 * time.Second)
		if userstates.Get(chatID) == "awaiting_refund_input" {
			userstates.Clear(chatID)
			bot.Send(&telebot.Chat{ID: chatID}, "⏳ Время на ввод истекло.")
		}
	}(chatID, c.Bot().(*telebot.Bot))

	return c.Send("💸 Введите `user_id` и `charge_id`, через пробел. Например:\n`123456789 9876-chrge`", &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	})
}