package autobuy

import (
	"gopkg.in/telebot.v4"
)

func HandleMinLimitChange(c telebot.Context) error {
	return c.Edit("Выберите новый минимум цены для автопокупки:", &telebot.SendOptions{
		ParseMode:   telebot.ModeMarkdown,
		ReplyMarkup: RenderMinLimitKeyboard(),
	})
}

func HandleMaxLimitChange(c telebot.Context) error {
	return c.Edit("Выберите новый максимум цены для автопокупки:", &telebot.SendOptions{
		ParseMode:   telebot.ModeMarkdown,
		ReplyMarkup: RenderMaxLimitKeyboard(),
	})
}