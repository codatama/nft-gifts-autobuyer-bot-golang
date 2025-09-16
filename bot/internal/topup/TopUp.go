package topup

import (
	"gopkg.in/telebot.v4"

	"context"
	"fmt"

	"prvbot/internal/db"
	"prvbot/internal/messages"
	"prvbot/internal/userstates"
)

func HandleTopUp(c telebot.Context) error {
	ctx := context.Background()
	chatID := c.Chat().ID

	user, err := db.GetUser(ctx, chatID)
	if err != nil {
		fmt.Println("❌ Ошибка при поиске пользователя:", err)
		return c.Send(messages.ErrInvoiceSend)
	}
	if user == nil {
		return c.Send("🚫 У вас нет доступа к приватному боту. Для подробной информации обратитесь к @exotical")
	}

	userstates.Set(chatID, "awaiting_topup")
	return c.Send(messages.TopUpCount)
}