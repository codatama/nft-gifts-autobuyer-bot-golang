package topup

import (
	"gopkg.in/telebot.v4"
	
	"context"
	"fmt"
	"math"

	"prvbot/internal/db"
	"prvbot/internal/messages"
)

func HandleSuccessfulPayment(c telebot.Context) error {
	payment := c.Message().Payment
	if payment == nil {
		return nil
	}

	chatID := c.Chat().ID

	var stars float64
	_, err := fmt.Sscanf(payment.Payload, "topup_%f", &stars)
	if err != nil {
		return c.Send(messages.ErrStarsCounter)
	}

	netStars := stars * (1 - db.GlobalCommissionRate)
	starCount := int(math.Round(netStars))

	_, err = db.GetUser(context.Background(), chatID)
	if err != nil {
		return c.Send(messages.ErrDataBaseRegister)
	}

	err = db.UpdateBalance(context.Background(), chatID, starCount)
	if err != nil {
		return c.Send(messages.ErrStarsUpdater)
	}

	return c.Send(fmt.Sprintf("✅ Оплата успешна! Вам начислено %.0f звёзд 🌟", stars))
}