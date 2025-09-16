package autobuy

import (
	"gopkg.in/telebot.v4"
	
	"prvbot/internal/messages"
)

func HandleSupplyLimit(c telebot.Context) error {
	return c.Edit(messages.MsgSupplyLimit, &telebot.SendOptions{
		ParseMode:   telebot.ModeMarkdown,
		ReplyMarkup: RenderSupplyLimitKeyboard(),
	})
}