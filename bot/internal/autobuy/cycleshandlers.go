package autobuy

import (
	"gopkg.in/telebot.v4"
	
	"prvbot/internal/messages"
)

func HandleCyclesCount(c telebot.Context) error {
	return c.Edit(messages.MsgAutoBuyCycles, &telebot.SendOptions{
		ParseMode:   telebot.ModeMarkdown,
		ReplyMarkup: RenderCyclesCountKeyboard(),
	})
}