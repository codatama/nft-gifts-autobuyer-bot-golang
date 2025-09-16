package adminpanel

import (
	"gopkg.in/telebot.v4"

	"time"
	
	"prvbot/internal/userstates"
)

func GrantAdminAccess(c telebot.Context) error {
	chatID := c.Chat().ID

	userstates.Set(chatID, "awaiting_admin_chat_id")

	go func(chatID int64, bot *telebot.Bot) {
		time.Sleep(30 * time.Second)
		if userstates.Get(chatID) == "awaiting_admin_chat_id" {
			userstates.Clear(chatID)
			bot.Send(&telebot.Chat{ID: chatID}, "⏳ Время на ввод истекло.")
		}
	}(chatID, c.Bot().(*telebot.Bot))

	return c.Send("🆔 Введите `chat_id` пользователя, которому хотите выдать доступ к админ-панели:", &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	})
}

func GrantRefundAccess(c telebot.Context) error {
	chatID := c.Chat().ID
	userstates.Set(chatID, "awaiting_refund_chat_id")

	go func(chatID int64, bot *telebot.Bot) {
		time.Sleep(30 * time.Second)
		if userstates.Get(chatID) == "awaiting_refund_chat_id" {
			userstates.Clear(chatID)
			bot.Send(&telebot.Chat{ID: chatID}, "⏳ Время на ввод истекло.")
		}
	}(chatID, c.Bot().(*telebot.Bot))

	return c.Send("💫 Введите `chat_id` пользователя, которому хотите выдать доступ к возврату звёзд:", &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	})
}

func GrantTechSupportAccess(c telebot.Context) error {
	chatID := c.Chat().ID
	userstates.Set(chatID, "awaiting_techsupport_chat_id")

	go func(chatID int64, bot *telebot.Bot) {
		time.Sleep(30 * time.Second)
		if userstates.Get(chatID) == "awaiting_techsupport_chat_id" {
			userstates.Clear(chatID)
			bot.Send(&telebot.Chat{ID: chatID}, "⏳ Время на ввод истекло.")
		}
	}(chatID, c.Bot().(*telebot.Bot))

	return c.Send("🧰 Введите `chat_id` пользователя, которому хотите выдать доступ к тех. поддержке:", &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	})
}