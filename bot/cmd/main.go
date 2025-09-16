package main

import (
	"log"
	"prvbot/config"
	"prvbot/internal/bot"
	"prvbot/internal/db"
)

func main() {
	cfg := config.Load()

	db.Init(cfg.DatabaseURL)
	telegramBot := bot.New(cfg.TelegramToken)

	log.Println("Бот успешно запущен...")
	telegramBot.Start()
}