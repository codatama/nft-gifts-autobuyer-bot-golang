package autobuy

import (
	"gopkg.in/telebot.v4"
)

func RenderCyclesCountKeyboard() *telebot.ReplyMarkup {
	markup := &telebot.ReplyMarkup{}

	markup.Inline(
		markup.Row(markup.Data("1", "cycles:1"), markup.Data("2", "cycles:2")),
		markup.Row(markup.Data("3", "cycles:3"), markup.Data("5", "cycles:5")),
		markup.Row(markup.Data("10", "cycles:10"), markup.Data("20", "cycles:20")),
		markup.Row(markup.Data("30", "cycles:30"), markup.Data("50", "cycles:30")),
		markup.Row(markup.Data("75", "cycles:75"), markup.Data("100", "cycles:100")),
		markup.Row(markup.Data("◀️ Вернуться назад", BtnGoBackToAutoBuyMenu.Unique)),
	)
	return markup
}