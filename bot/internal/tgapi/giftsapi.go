package tgapi

import (
	"gopkg.in/telebot.v4"
	
	"encoding/json"
	"fmt"
	"context"

	"prvbot/internal/db"
	"prvbot/internal/models"
)

type TelegramGift struct {
	ID        string `json:"id"`
	StarCount int    `json:"star_count"`
	TotalCount       *int `json:"total_count,omitempty"`

	Sticker struct {
		Emoji     string `json:"emoji"`
		FileID    string `json:"file_id"`
		Thumbnail struct {
			FileID string `json:"file_id"`
		} `json:"thumbnail"`
	} `json:"sticker"`
}

type BuyGiftResponse struct {
	OK     bool   `json:"ok"`
	Result struct {
		Success bool `json:"success"`
	} `json:"result"`
}

type getAvailableGiftsResponse struct {
	OK     bool `json:"ok"`
	Result struct {
		Gifts []TelegramGift `json:"gifts"`
	} `json:"result"`
}

func GetAvailableGifts(bot *telebot.Bot) ([]TelegramGift, error) {
	var parsed getAvailableGiftsResponse

	resp, err := bot.Raw("getAvailableGifts", nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка при вызове Raw: %w", err)
	}

	err = json.Unmarshal(resp, &parsed)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	if !parsed.OK {
		return nil, fmt.Errorf("ответ Telegram содержит ok = false")
	}

	return parsed.Result.Gifts, nil
}

func SendGiftToUser(bot *telebot.Bot, chatID interface{}, giftID string) error {
	params := map[string]interface{}{
		"chat_id":  chatID,
		"gift_id":  giftID,
	}
	resp, err := bot.Raw("sendGift", params)
	if err != nil {
		return fmt.Errorf("ошибка при вызове sendGift: %w", err)
	}

	fmt.Println("✅ sendGift:", string(resp))
	return nil
}

func SyncGiftsWithDatabase(bot *telebot.Bot) error {
	gifts, err := GetAvailableGifts(bot)
	if err != nil {
		return fmt.Errorf("не удалось получить список подарков: %w", err)
	}

	for _, g := range gifts {
		giftModel := &models.GiftEntity{
			ApiID:     g.ID,
			StarCount: g.StarCount,
			Emoji:     g.Sticker.Emoji,
			FileID:    g.Sticker.FileID,
			Thumbnail: g.Sticker.Thumbnail.FileID,
		}

		fmt.Printf("🐞 g.ID: %s | Emoji: %s | FileID: %s\n", g.ID, g.Sticker.Emoji, g.Sticker.FileID)

		err := db.UpsertTelegramGift(context.Background(), giftModel)
		if err != nil {
			fmt.Printf("⚠️ Ошибка при сохранении %s: %v\n", g.Sticker.Emoji, err)
		}
	}

	fmt.Printf("✅ Синхронизировано %d подарков\n", len(gifts))
	return nil
}