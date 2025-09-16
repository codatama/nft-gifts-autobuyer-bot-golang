package db

import (
	"gopkg.in/telebot.v4"

	"context"
	"log"
)

func SetBuyerUser(c telebot.Context) error {
	_, err := Pool.Exec(context.Background(), `
		UPDATE users
		SET channel_enabled = FALSE
		WHERE chat_id = $1
	`, c.Sender().ID)

	if err != nil {
		log.Printf("❌ Ошибка при обновлении channel_enabled для пользователя %d: %v", c.Sender().ID, err)
		return err
	}

	log.Printf("✅ channel_enabled = FALSE установлен для пользователя %d", c.Sender().ID)
	return nil
	
}

func SetBuyerChannel(c telebot.Context) error {
	_, err := Pool.Exec(context.Background(), `
		UPDATE users
		SET channel_enabled = TRUE
		WHERE chat_id = $1
	`, c.Sender().ID)

	if err != nil {
		log.Printf("❌ Ошибка при установке channel_enabled=TRUE для %d: %v", c.Sender().ID, err)
		return err
	}

	log.Printf("✅ Пользователь %d активировал доставку через канал", c.Sender().ID)
	return nil
}