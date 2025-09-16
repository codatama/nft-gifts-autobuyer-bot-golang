package autobuy

import (
	"gopkg.in/telebot.v4"

	"fmt"
	"context"
	
	"prvbot/internal/db"
)

func HandleMinLimitCallback(c telebot.Context, limit int) error {
	err := db.UpdateMinLimit(context.Background(), c.Chat().ID, limit)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "❌ Не удалось сохранить лимит"})
	}

	c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("✅ Минимум: %d ⭐️", limit)})
	return HandleAutoBuySettings(c)
}

func HandleMaxLimitCallback(c telebot.Context, limit int) error {
	err := db.UpdateMaxLimit(context.Background(), c.Chat().ID, limit)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "❌ Не удалось сохранить лимит!"})
	}

	c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("✅ Максимум: %d ⭐️", limit)})
	return HandleAutoBuySettings(c)
}