package autobuy

import (
	"gopkg.in/telebot.v4"

	"context"
	"log"
	
	"prvbot/internal/db"
)

func HandleAutoBuySettings(c telebot.Context) error {
	user, err := db.GetUserByChatID(context.Background(), c.Chat().ID)
	if err != nil {
		log.Println("⚠️ Ошибка загрузки пользователя:", err)
		return c.Send("❗ Не удалось загрузить настройки. Сначала пополньте баланс.")
	}

	text := RenderAutoBuySettings(user)
	replyMarkup := RenderAutoBuyKeyboard(user)

	if c.Callback() != nil {
		return c.Edit(text, &telebot.SendOptions{
			ParseMode:   telebot.ModeMarkdown,
			ReplyMarkup: replyMarkup,
		})
	}

	return c.Send(text, &telebot.SendOptions{
		ParseMode:   telebot.ModeMarkdown,
		ReplyMarkup: replyMarkup,
	})
}