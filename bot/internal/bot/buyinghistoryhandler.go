package bot

import (
	"gopkg.in/telebot.v4"
	"prvbot/internal/buttons"
)

var BtnGift1 = &telebot.Btn{
	Text:   "№1",
	Unique: "gift_1",
}

var BtnGift2 = &telebot.Btn{
	Text:   "№2",
	Unique: "gift_2",
}

var BtnGift3 = &telebot.Btn{
	Text:   "№3",
	Unique: "gift_3",
}

func HandleBuyingHistory(c telebot.Context) error {
	text := "🎁 *Недавно купленные подарки в боте:*"

	markup := &telebot.ReplyMarkup{}
	btnBack := markup.Data("◀️ Назад", buttons.BtnBackToMenu.Unique)

	markup.Inline(
		markup.Row(*BtnGift1),
		markup.Row(*BtnGift2),
		markup.Row(*BtnGift3),
		markup.Row(btnBack),
	)

	return c.Edit(text, &telebot.SendOptions{
		ParseMode:   telebot.ModeMarkdown,
		ReplyMarkup: markup,
	})
}

func HandleGift1(c telebot.Context) error {
	text := `🎁 *Покупка #1*

👤 Куплен: *Имя скрыто*  
💰 Стоимость: *50* ⭐️  
🕓 Время покупки: *15:11 17-07-2025*

[🔗 Содержимое подарка](https://t.me/giftsids/19)`

	return c.Send(text, &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
		DisableWebPagePreview: false,
	})
}

func HandleGift2(c telebot.Context) error {
	text := `🎁 *Покупка #2*

👤 Куплен: *Имя скрыто*  
💰 Стоимость: *50* ⭐️  
🕓 Время покупки: *11:55 12-07-2025*

[🔗 Содержимое подарка](https://t.me/giftsids/18)`

	return c.Send(text, &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
		DisableWebPagePreview: false,
	})
}

func HandleGift3(c telebot.Context) error {
	text := `🎁 *Покупка #3*

👤 Куплен: *Имя скрыто*  
💰 Стоимость: *100* ⭐️  
🕓 Время покупки: *23:05 30-06-2025*

[🔗 Содержимое подарка](https://t.me/giftsids/21)`

	return c.Send(text, &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
		DisableWebPagePreview: false,
	})
}