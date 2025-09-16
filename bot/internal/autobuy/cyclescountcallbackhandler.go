package autobuy

import (
	"gopkg.in/telebot.v4"

	"context"
	"fmt"
	
	"prvbot/internal/db"
)

func HandleCyclesCountCallback(c telebot.Context, cycles int) error {
	err := db.UpdateCyclesCount(context.Background(), c.Chat().ID, cycles)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "❌ Не удалось сохранить количество циклов!"})
	}

	c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("✅ Циклов: %d", cycles)})
	return HandleAutoBuySettings(c)
}