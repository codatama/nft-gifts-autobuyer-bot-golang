package autobuy

import (
	"context"

	"gopkg.in/telebot.v4"

	"fmt"
	"log"
	"strings"

	"prvbot/internal/db"
	"prvbot/internal/userstates"
)

func HandleChannelSettings(c telebot.Context) error {
	chatID := c.Sender().ID
	ctx := context.Background()

	channels, err := db.GetOrCreatePurchaseChannels(ctx, chatID)
	if err != nil {
		log.Println("❌ Ошибка получения каналов:", err)
		return c.Send("⚠️ Ошибка при загрузке каналов.")
	}

	text := "*Настройки каналов*\n\n"
	text += fmt.Sprintf("▪️ Канал 1: %s\n", emptyOrValue(channels.Channel1))
	text += fmt.Sprintf("▪️ Канал 2: %s\n", emptyOrValue(channels.Channel2))
	text += fmt.Sprintf("▪️ Канал 3: %s\n", emptyOrValue(channels.Channel3))

	text += "\nВНИМАНИЕ!\n"
	text += "Канал 1 покупает подарки с сапплаем не выше 15 тысяч\n"
	text += "Канал 2 — не выше 50 тысяч\n"
	text += "Канал 3 — свыше 50 тысяч\n\n"
	text += "Если хотите, чтобы любой подарок мог быть куплен на канал — укажите один и тот же канал во всех полях."


	menu := &telebot.ReplyMarkup{}
	btnAdd := menu.Data("➕ Добавить канал", "add_channel")
	btnBack := menu.Data("⬅️ Назад", BtnGoBackToAutoBuyMenu.Unique)
	menu.Inline(menu.Row(btnAdd), menu.Row(btnBack))

	return c.Edit(text, &telebot.SendOptions{
		ParseMode:   telebot.ModeMarkdown,
		ReplyMarkup: menu,
	})
}

func emptyOrValue(s string) string {
	if strings.TrimSpace(s) == "" {
		return "_Канал не указан_"
	}
	return s
}

func HandleAddChannelStepOne(c telebot.Context) error {
	text := "📡 Выберите номер канала, который хотите изменить:"

	menu := &telebot.ReplyMarkup{}
	btn1 := menu.Data("1️⃣", "channel:edit1")
	btn2 := menu.Data("2️⃣", "channel:edit2")
	btn3 := menu.Data("3️⃣", "channel:edit3")
	btnb := menu.Data("⬅️ Назад", BtnGoBackToAutoBuyMenu.Unique)
	menu.Inline(menu.Row(btn1, btn2, btn3),
	menu.Row(btnb))
	

	return c.Edit(text, &telebot.SendOptions{
		ReplyMarkup: menu,
	})
}

func PromptChannelInput(c telebot.Context, field string) error {
	chatID := c.Chat().ID

	ctx := context.Background()
	channels, err := db.GetOrCreatePurchaseChannels(ctx, chatID)
	if err != nil {
		log.Printf("❌ Не удалось получить каналы: %v", err)
		return c.Send("⚠️ Ошибка при загрузке каналов.")
	}

	switch field {
	case "channel2":
		if strings.TrimSpace(channels.Channel1) == "" {
			return c.Send("⚠️ У вас ещё нет канала 1. Для начала создайте его.")
		}
	case "channel3":
		if strings.TrimSpace(channels.Channel2) == "" {
			return c.Send("⚠️ У вас ещё нет канала 2. Для начала создайте его.")
		}
	}

	switch field {
	case "channel1":
		userstates.Set(chatID, "awaiting_channel_1")
	case "channel2":
		userstates.Set(chatID, "awaiting_channel_2")
	case "channel3":
		userstates.Set(chatID, "awaiting_channel_3")
	}

	return c.Send("✍️ Введите название канала в формате `@channelname`:", &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	})
}

func HandleChannelNameInput(c telebot.Context) error {
	chatID := c.Chat().ID
	state := userstates.Get(chatID)
	userstates.Clear(chatID)

	field := map[string]string{
		"awaiting_channel_1": "channel1",
		"awaiting_channel_2": "channel2",
		"awaiting_channel_3": "channel3",
	}[state]

	if field == "" {
		return nil
	}

	channelName := strings.TrimSpace(c.Text())
	if !strings.HasPrefix(channelName, "@") {
		return c.Send("❌ Название должно начинаться с `@`. Попробуйте снова.")
	}

	err := db.UpdateChannelField(chatID, field, channelName)
	if err != nil {
		log.Println("❌ Ошибка при сохранении:", err)
		return c.Send("⚠️ Не удалось сохранить канал.")
	}

	return c.Send(fmt.Sprintf("✅ Установлен %s: %s", field, channelName))
}