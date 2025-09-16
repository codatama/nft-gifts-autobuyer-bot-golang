package adminpanel

import (
	"gopkg.in/telebot.v4"

	"context"
	"fmt"
	"log"
	
	"prvbot/internal/db"
)

func SendAdminPanel(c telebot.Context) error {
	hasAccess, err := db.HasAdminAccess(context.Background(), c.Chat().ID)
	if err != nil {
		log.Println("❗ Ошибка проверки прав доступа:", err)
		return c.Send("⚠️ Ошибка при проверке прав доступа.")
	}
	if !hasAccess {
		return c.Send("⛔ У тебя нет доступа к админ-панели.")
	}

	userCount, err := db.GetUserCount(context.Background())
	if err != nil {
		log.Printf("❗ Не удалось получить количество пользователей: %v", err)
		userCount = 0
	}

	increase, err := db.GetLatestBalanceDifference(context.Background())
	if err != nil {
	log.Printf("⚠️ Ошибка прироста звёзд: %v", err)
	}

	InitButtonsAdmin()
	text := fmt.Sprintf("🔧 *Админ-панель*\n" + "Количество пользователей: %d\n" + "Прирост звёзд за сутки: %d\n", userCount, increase)
	return c.Send(text, AdminPanelInline, telebot.ModeMarkdown)
}

func EditAdminPanel(c telebot.Context) error {
	hasAccess, err := db.HasAdminAccess(context.Background(), c.Chat().ID)
	if err != nil {
		log.Println("❗ Ошибка проверки прав доступа:", err)
		return c.Send("⚠️ Ошибка при проверке прав доступа.")
	}
	if !hasAccess {
		return c.Send("⛔ У тебя нет доступа к админ-панели.")
	}

	userCount, err := db.GetUserCount(context.Background())
	if err != nil {
		log.Printf("❗ Не удалось получить количество пользователей: %v", err)
		userCount = 0
	}

	increase, err := db.GetLatestBalanceDifference(context.Background())
	if err != nil {
	log.Printf("⚠️ Ошибка прироста звёзд: %v", err)
	}

	InitButtonsAdmin()
	text := fmt.Sprintf("🔧 *Админ-панель*\n" + "Количество пользователей: %d\n" + "Прирост звёзд за сутки: %d\n", userCount, increase)
	return c.Edit(text, AdminPanelInline, telebot.ModeMarkdown)
}