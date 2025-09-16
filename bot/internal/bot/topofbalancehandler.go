package bot

import (
	"gopkg.in/telebot.v4"

	"context"
	"fmt"
	"strings"

	"prvbot/internal/db"
	"prvbot/internal/buttons"
)

func HandleTopOfBalance(c telebot.Context) error {
	ctx := context.Background()
	userID := c.Chat().ID

	users, err := db.GetTopUsersByBalance(ctx, 10) // Получаем топ-10
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "❗ Ошибка при загрузке топа."})
	}

	currentUser, err := db.GetUserByChatID(ctx, userID)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "❗ Не удалось загрузить ваш профиль."})
	}

	position, err := db.GetUserBalanceRank(ctx, currentUser.ChatID)
	if err != nil {
		position = -1
	}

	var list string
	for i, u := range users {
		name := "Имя скрыто"
		if u.Username != "" {
			name = EscapeMarkdown("@" + u.Username)
		}
		list += fmt.Sprintf("%d. %s — %d⭐️\n", i+1, name, u.Balance)
	}

	youLine := "❌ Не найдено в рейтинге"
	if position > 0 {
		selfName := "Имя скрыто"
		if currentUser.Username != "" {
			selfName = EscapeMarkdown("@" + currentUser.Username)
		}
		youLine = fmt.Sprintf("№%d — %s (%d⭐️)", position, selfName, currentUser.Balance)
	}

	text := fmt.Sprintf("*🏅 Топ по балансу:*\n\n%s\n👤 *Вы:*\n%s", list, youLine)

	markup := &telebot.ReplyMarkup{}
	backBtn := markup.Data("◀️ Назад", buttons.BtnBackToMenu.Unique)
	markup.Inline(markup.Row(backBtn))

	return c.Edit(text, &telebot.SendOptions{
		ParseMode:   telebot.ModeMarkdown,
		ReplyMarkup: markup,
	})
}

func EscapeMarkdown(s string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"~", "\\~",
		"`", "\\`",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		".", "\\.",
		"!", "\\!",
	)
	return replacer.Replace(s)
}