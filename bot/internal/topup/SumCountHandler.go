package topup

import (
	"gopkg.in/telebot.v4"
	
	"fmt"
	"strconv"

	"prvbot/internal/messages"
	"prvbot/internal/tgapi"
	"prvbot/internal/userstates"
)

func HandleTopUpInput(c telebot.Context) error {
	chatID := c.Chat().ID

	userstates.Set(chatID, "awaiting_topup")

	amountText := c.Text()
	amount, err := strconv.Atoi(amountText)
	if err != nil || amount <= 0 {
		return c.Send(messages.NotCorrectNumberOfStars)
	}
	if amount < 50 {
	return c.Send("❌ Минимальная сумма пополнения — 50⭐️. Попробуйте ввести больше.")
	}
	if amount > 100_000 {
		return c.Send(messages.TooBigTopUp)
	}

	userstates.Clear(chatID)

	if err := tgapi.SendStarsInvoice(chatID, float64(amount)); err != nil {
		fmt.Println("Ошибка отправки инвойса:", err)
		return c.Send("❌ Ошибка при создании инвойса. Попробуйте позже.")
	}

	return nil
}