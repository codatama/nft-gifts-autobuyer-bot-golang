package autobuy

import (
	"gopkg.in/telebot.v4"

	"context"
	
	"prvbot/internal/db"
)

func HandleToggleAutoBuy(c telebot.Context) error {
	chatID := c.Chat().ID

	err := db.ToggleAutoBuy(context.Background(), chatID)
	if err != nil {
		return c.Send("⚠️ Не удалось переключить автопокупку")
	}

	return HandleAutoBuySettings(c)
}