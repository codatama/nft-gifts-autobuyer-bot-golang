package autobuy

import (
	"gopkg.in/telebot.v4"
)

func RenderChooseTheBuyerKeyboard() *telebot.ReplyMarkup {
	markup := &telebot.ReplyMarkup{}

	markup.Inline(
		markup.Row(markup.Data("Выбрать себя", "target:self")),
		markup.Row(markup.Data("Выбрать канал", "target:channel")),
		markup.Row(markup.Data("◀️ Вернуться назад", BtnGoBackToAutoBuyMenu.Unique)),
	)
	return markup
}