package bot

import (
	"gopkg.in/telebot.v4"

	"log"
	"context"
	"fmt"

	"prvbot/internal/db"
	"prvbot/internal/buttons"
)

func HandleUserProfile(c telebot.Context) error {
	user, err := db.GetUserByChatID(context.Background(), c.Chat().ID)
	if err != nil {
		log.Println("⚠️ Ошибка загрузки профиля:", err)
		return c.Respond(&telebot.CallbackResponse{Text: "❗ Не удалось загрузить профиль. Сначала пополньте баланс."})
	}

	text := fmt.Sprintf("👤 *Ваш профиль*\n\n💳 Баланс: *%d* ⭐️", user.Balance)

	markup := &telebot.ReplyMarkup{}
	btnTopUp := markup.Data("💸 Пополнить баланс", buttons.BtnTopUp.Unique)
	btnRefundStars := markup.Data("⚒️ Техническая поддержка", buttons.BtnRefundStars.Unique)
	btnBack := markup.Data("◀️ Назад", buttons.BtnBackToMenu.Unique)
	btnTopOfBalance := markup.Data("🏅 Топ по балансу", buttons.BtnTopOfBalance.Unique)

	markup.Inline(
		markup.Row(btnTopUp),
		markup.Row(btnTopOfBalance),
		markup.Row(btnRefundStars),
		markup.Row(btnBack),
	)

	return c.Edit(text, &telebot.SendOptions{
		ParseMode:   telebot.ModeMarkdown,
		ReplyMarkup: markup,
	})
}