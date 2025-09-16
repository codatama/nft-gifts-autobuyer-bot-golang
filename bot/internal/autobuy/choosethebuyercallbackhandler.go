package autobuy

import (
	"gopkg.in/telebot.v4"

	"log"
	"context"
	"strings"

	"prvbot/internal/db"
	"prvbot/internal/models"
)

func HandleChooseTheBuyerCallback(c telebot.Context) error {
	log.Println("🔥 Вызван HandleChooseTheBuyerCallback")
	data := strings.TrimSpace(strings.Trim(c.Callback().Data, "\f\r\n"))
	log.Println("👉 Внутри хендлера, data =", data)

	var (
		err      error
		channels *models.PurchaseChannels
	)

	switch data {
	case "target:channel":
		channels, err = db.GetPurchaseChannelsByChatID(context.Background(), c.Chat().ID)
		if err != nil {
			log.Println("⚠️ Ошибка при получении каналов:", err)
			return c.Respond(&telebot.CallbackResponse{Text: "❌ Ошибка при проверке каналов!"})
		}

		if channels == nil ||
			strings.TrimSpace(channels.Channel1) == "" &&
				strings.TrimSpace(channels.Channel2) == "" &&
				strings.TrimSpace(channels.Channel3) == "" {
			return c.Respond(&telebot.CallbackResponse{Text: "⚠️ У вас ещё не зарегистрировано ни одного канала!"})
		}

		err = db.SetBuyerChannel(c)

	case "target:self":
		err = db.SetBuyerUser(c)

	default:
		return c.Respond(&telebot.CallbackResponse{Text: "❌ Неизвестный выбор получателя"})
	}

	if err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "❌ Не удалось изменить получателя!"})
	}

	_ = c.Respond(&telebot.CallbackResponse{Text: "✅ Получатель успешно установлен!"})
	return HandleAutoBuySettings(c)
}