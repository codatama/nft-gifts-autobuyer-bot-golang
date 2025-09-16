package userstates

import (
	"gopkg.in/telebot.v4"
	"time"
)

var Bot *telebot.Bot

func Init(bot *telebot.Bot) {
	Bot = bot
}


var UserState = map[int64]string{}

func Set(chatID int64, state string) {
	UserState[chatID] = state

	go func(id int64, expected string) {
	time.Sleep(30 * time.Second)
	if UserState[id] == expected {
		Clear(id)
		if Bot != nil {
			Bot.Send(&telebot.Chat{ID: id}, "⏳ Время на ввод истекло.")
		}
	}
	}(chatID, state)
}

func Get(chatID int64) string {
	return UserState[chatID]
}

func Clear(chatID int64) {
	delete(UserState, chatID)
}