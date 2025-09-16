package gifting

import (
	"gopkg.in/telebot.v4"

	"prvbot/internal/messages"
	"prvbot/internal/userstates"
)
func HandleSendGift(c telebot.Context) error {
	chatID := c.Chat().ID

	userstates.Set(chatID, "awaiting_gift_number")

	return c.Send(messages.UnCorrectGiftNumberOfGift)
}