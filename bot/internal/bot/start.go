package bot

import (
	"gopkg.in/telebot.v4"

	"context"
	"log"
	"fmt"

	"prvbot/internal/db"
	"prvbot/internal/buttons"
)

func ProtectedStart(c telebot.Context) error {
	ctx := context.Background()
	chatID := c.Chat().ID

	user, err := db.GetUser(ctx, chatID)
	if err != nil {
		log.Printf("⚠️ Ошибка поиска пользователя: %v", err)
		return c.Send("❌ Возникла ошибка при попытке авторизации")
	}

	if user == nil {
		return c.Send("🚫 У вас нет доступа к приватному боту. Для подробной информации обратитесь к @exotical")
	}

	return Start(c)
}

func Start(c telebot.Context) error {
	buttons.InitButtons()

	total, err := db.GetTotalBalance(context.Background())
	if err != nil {
		log.Println("❌ Ошибка при получении общего баланса:", err)
		total = 0
	}

	text := fmt.Sprintf(
		"Привет! Перед тобой *самый быстрый* и *новый* Telegram-бот автоматической покупки подарков.\n\n"+
			"С ним ты *быстрее всех* подберёшь и приобретёшь новые подарки, а значит — избавишься от стресса и сможешь спокойно выспаться.\n\n"+
			"Бот работает *уже почти год* и стал *одним из главных инструментов* на рынке подарков.\n\n"+
			"На сегодняшний день в боте лежит сумма звёзд в более чем *%d ⭐️*!", total)

	return c.Send(text, &telebot.SendOptions{
		ParseMode:   telebot.ModeMarkdown,
		ReplyMarkup: buttons.InlineMenu,
	})
}

func tgbStartMessage() (string, *telebot.SendOptions) {
	buttons.InitButtons()

	total, err := db.GetTotalBalance(context.Background())
	if err != nil {
		log.Println("❌ Ошибка при получении общего баланса:", err)
		total = 0
	}

	text := fmt.Sprintf(
		"Привет! Перед тобой *самый быстрый* и *новый* Telegram-бот автоматической покупки подарков.\n\n"+
			"С ним ты *быстрее всех* подберёшь и приобретёшь новые подарки, а значит — избавишься от стресса и сможешь спокойно выспаться.\n\n"+
			"Бот работает *уже почти год* и стал *одним из главных инструментов* на рынке подарков.\n\n"+
			"На сегодняшний день в боте лежит сумма звёзд в более чем *%d ⭐️*!", total)

	return text, &telebot.SendOptions{
		ParseMode:   telebot.ModeMarkdown,
		ReplyMarkup: buttons.InlineMenu,
	}
}