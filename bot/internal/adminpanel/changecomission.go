package adminpanel

import (
	"gopkg.in/telebot.v4"

	"strings"
	"strconv"
	"fmt"

	"prvbot/internal/db"
)

func HandleSetCommissionPrompt(c telebot.Context) error {
	return c.Send("💬 Введите желаемую комиссию\n(значение от *0.02* до *1*)", &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	})
}

func HandleCommissionInput(c telebot.Context) error {
	text := strings.TrimSpace(c.Text())
	rate, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return c.Send("❌ Неверный формат. Введите число от 0.02 до 1 (например: 0.05)")
	}
	if rate < 0.02 || rate > 1 {
		return c.Send("⚠️ Комиссия должна быть от 0.02 до 1")
	}

	db.GlobalCommissionRate = rate

	return c.Send(fmt.Sprintf("✅ Комиссия успешно установлена: %.2f%%", rate*100))
}