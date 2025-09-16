package autobuy

import (
	"gopkg.in/telebot.v4"
)

func HandleChooseTheBuyer(c telebot.Context) error {
	return c.Edit("Выберите, кто будет получать подарки при автопокупке:", &telebot.SendOptions{
		ParseMode:   telebot.ModeMarkdown,
		ReplyMarkup: RenderChooseTheBuyerKeyboard(),
	})
}