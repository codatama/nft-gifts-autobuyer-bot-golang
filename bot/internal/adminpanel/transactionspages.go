package adminpanel

import (
	"gopkg.in/telebot.v4"

	"context"
	"log"
	"time"
	"fmt"
	
	"prvbot/internal/db"
)

func ShowTransactionsPage(c telebot.Context, page int) error {
	const pageSize = 10
	offset := page * pageSize

	rows, err := db.Query(context.Background(),
		`SELECT id, amount, nanostar_amount, date, source, receiver 
		 FROM transactions ORDER BY date DESC LIMIT $1 OFFSET $2`, pageSize, offset)
	if err != nil {
		return c.Edit("❌ Ошибка при загрузке транзакций")
	}
	defer rows.Close()

	var text string
	count := 0
	for rows.Next() {
		var id string
		var amount, nano int
		var date, source, receiver int64

		rows.Scan(&id, &amount, &nano, &date, &source, &receiver)
		t := time.Unix(date, 0)
		text += fmt.Sprintf("🧾 *ID:*`%s`\n🪙 *%d* звёзд (%d nano)\n👤 %d → %d\n📅 %s\n\n",
			id, amount, nano, source, receiver, t.Format("02.01.2006 15:04"))
		count++
	}

	if count == 0 {
		return c.Edit("⚠️ Записей пока нет")
	}

	back := ""
	if page > 0 {
	back = fmt.Sprintf("txpage:%d", page-1)
	}
	log.Printf("👉 Back: %s | Next: %s | Menu: %s", back, fmt.Sprintf("txpage:%d", page+1), "admin_back")
	keyboard := PaginationButtons(back, fmt.Sprintf("txpage:%d", page+1), "admin_back")
	log.Printf("🧩 InlineKeyboard: %#v", keyboard.InlineKeyboard)

	return c.Edit(text, &telebot.SendOptions{
	ParseMode:   telebot.ModeMarkdown,
	ReplyMarkup: keyboard,
})
}