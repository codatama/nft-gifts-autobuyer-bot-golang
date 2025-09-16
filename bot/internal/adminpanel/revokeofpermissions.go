package adminpanel

import (
	"gopkg.in/telebot.v4"

	"time"
	
	"prvbot/internal/userstates"
)

func RevokeAdminAccess(c telebot.Context) error {
	chatID := c.Chat().ID
	userstates.Set(chatID, "awaiting_admin_revoke")

	go func(chatID int64, bot *telebot.Bot) {
		time.Sleep(30 * time.Second)
		if userstates.Get(chatID) == "awaiting_admin_revoke" {
			userstates.Clear(chatID)
			bot.Send(&telebot.Chat{ID: chatID}, "⏳ Время на ввод истекло.")
		}
	}(chatID, c.Bot().(*telebot.Bot))

	return c.Send("🛠 Введите `chat_id`, у которого хотите _забрать_ права на админ-панель:", &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	})
}

func RevokeRefundAccess(c telebot.Context) error {
	chatID := c.Chat().ID
	userstates.Set(chatID, "awaiting_refund_revoke")

	go func(chatID int64, bot *telebot.Bot) {
		time.Sleep(30 * time.Second)
		if userstates.Get(chatID) == "awaiting_refund_revoke" {
			userstates.Clear(chatID)
			bot.Send(&telebot.Chat{ID: chatID}, "⏳ Время на ввод истекло.")
		}
	}(chatID, c.Bot().(*telebot.Bot))

	return c.Send("🛠 Введите `chat_id`, у которого хотите _забрать_ права на refund:", &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	})
}

func RevokeTechnicalSupport(c telebot.Context) error {
	chatID := c.Chat().ID
	userstates.Set(chatID, "awaiting_tech_revoke")

	go func(chatID int64, bot *telebot.Bot) {
		time.Sleep(30 * time.Second)
		if userstates.Get(chatID) == "awaiting_tech_revoke" {
			userstates.Clear(chatID)
			bot.Send(&telebot.Chat{ID: chatID}, "⏳ Время на ввод истекло.")
		}
	}(chatID, c.Bot().(*telebot.Bot))

	return c.Send("🛠 Введите `chat_id`, у которого хотите _забрать_ права на техническую поддержку:", &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	})
}