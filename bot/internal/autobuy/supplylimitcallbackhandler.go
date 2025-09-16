package autobuy

import (
	"gopkg.in/telebot.v4"

	"context"
	"fmt"
	
	"prvbot/internal/db"
)

func HandleSupplyLimitCallback(c telebot.Context, supply int) error {
	err := db.UpdateSupplyLimit(context.Background(), c.Chat().ID, supply)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "❌ Не удалось сохранить лимит саплая!"})
	}

	c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("✅ Лимит сапплая: %d", supply)})
	return HandleAutoBuySettings(c)
}