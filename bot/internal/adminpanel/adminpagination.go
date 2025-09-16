package adminpanel

import (
	"gopkg.in/telebot.v4"
)

func PaginationButtons(backData, nextData, menuData string) *telebot.ReplyMarkup {
	markup := &telebot.ReplyMarkup{}
	var row []telebot.InlineButton

	if backData != "" {
		back := telebot.InlineButton{Text: "⬅️ Назад", Data: backData}
		row = append(row, back)
	}
	if nextData != "" {
		next := telebot.InlineButton{Text: "➡️ Вперёд", Data: nextData}
		row = append(row, next)
	}
	if len(row) > 0 {
		markup.InlineKeyboard = append(markup.InlineKeyboard, row)
	}

	if menuData != "" {
		menu := telebot.InlineButton{Text: "🔙 В админку", Data: menuData}
		markup.InlineKeyboard = append(markup.InlineKeyboard,
			[]telebot.InlineButton{menu},
		)
	}

	return markup
}