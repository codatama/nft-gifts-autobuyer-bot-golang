package autobuy

import (
	"gopkg.in/telebot.v4"
)

func RenderMinLimitKeyboard() *telebot.ReplyMarkup {
	markup := &telebot.ReplyMarkup{}

	markup.Inline(
		markup.Row(markup.Data("0 ⭐️", "min:0"), markup.Data("15 ⭐️", "min:15")),
		markup.Row(markup.Data("25 ⭐️", "min:25"), markup.Data("50 ⭐️", "min:50")),
		markup.Row(markup.Data("100 ⭐️", "min:100"), markup.Data("200 ⭐️", "min:200")),
		markup.Row(markup.Data("500 ⭐️", "min:500"), markup.Data("1000 ⭐️", "min:1000")),
		markup.Row(markup.Data("1500 ⭐️", "min:1500"), markup.Data("2000 ⭐️", "min:2000")),
		markup.Row(markup.Data("2500 ⭐️", "min:2500"), markup.Data("3000 ⭐️", "min:3000")),
		markup.Row(markup.Data("5000 ⭐️", "min:5000"), markup.Data("10000 ⭐️", "min:10000")),
		markup.Row(markup.Data("20000 ⭐️", "min:20000"), markup.Data("25000 ⭐️", "min:25000")),
		markup.Row(markup.Data("◀️ Назад", BtnGoBackToAutoBuyMenu.Unique)),
	)

	return markup
}

func RenderMaxLimitKeyboard() *telebot.ReplyMarkup {
	markup := &telebot.ReplyMarkup{}

	markup.Inline(
		markup.Row(markup.Data("15 ⭐️", "max:15"), markup.Data("25 ⭐️", "max:25")),
		markup.Row(markup.Data("50 ⭐️", "max:50"), markup.Data("100 ⭐️", "max:100")),
		markup.Row(markup.Data("200 ⭐️", "max:200"), markup.Data("500 ⭐️", "max:500")),
		markup.Row(markup.Data("1000 ⭐️", "max:1000"), markup.Data("1500 ⭐️", "max:1500")),
		markup.Row(markup.Data("2000 ⭐️", "max:2000"), markup.Data("2500 ⭐️", "max:2500")),
		markup.Row(markup.Data("3000 ⭐️", "max:3000"), markup.Data("5000 ⭐️", "max:5000")),
		markup.Row(markup.Data("7500 ⭐️", "max:7500"), markup.Data("10000 ⭐️", "max:10000")),
		markup.Row(markup.Data("15000 ⭐️", "max:15000"), markup.Data("20000 ⭐️", "max:20000")),
		markup.Row(markup.Data("25000 ⭐️", "max:25000")),
		markup.Row(markup.Data("◀️ Назад", BtnGoBackToAutoBuyMenu.Unique)),
	)

	return markup
}