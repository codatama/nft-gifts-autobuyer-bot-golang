package autobuy

import (
	"gopkg.in/telebot.v4"
)

func RenderSupplyLimitKeyboard() *telebot.ReplyMarkup {
	markup := &telebot.ReplyMarkup{}

	markup.Inline(
		markup.Row(markup.Data("500", "supply:500"), markup.Data("1000", "supply:1000")),
		markup.Row(markup.Data("1500", "supply:1500"), markup.Data("1999", "supply:1999")),
		markup.Row(markup.Data("2000", "supply:2000"), markup.Data("3000", "supply:3000")),
		markup.Row(markup.Data("5000", "supply:5000"), markup.Data("7500", "supply:7500")),
		markup.Row(markup.Data("10000", "supply:10000"), markup.Data("15000", "supply:15000")),
		markup.Row(markup.Data("25000", "supply:25000"), markup.Data("50000", "supply:50000")),
		markup.Row(markup.Data("100000", "supply:100000"), markup.Data("250000", "supply:250000")),
		markup.Row(markup.Data("500000", "supply:500000"), markup.Data("1000000", "supply:1000000")),
		markup.Row(markup.Data("◀️ Вернуться назад", BtnGoBackToAutoBuyMenu.Unique)),
	)
	return markup
}