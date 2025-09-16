package gifting

import (
	"gopkg.in/telebot.v4"

	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"prvbot/internal/db"
	"prvbot/internal/messages"
	"prvbot/internal/userstates"
	"prvbot/internal/utils"
)

func HandleGiftInput(c telebot.Context) error {
	userID := c.Sender().ID
	input := c.Text()
	chatID := c.Chat().ID

	userstates.Set(chatID, "awaiting_gift_number")

	parts := strings.Split(input, ",")
	if len(parts) == 0 {
    return c.Send(messages.ErrInvalidGiftNumber)
	}

	indexStr := strings.TrimSpace(parts[0])
	index, err := strconv.Atoi(indexStr)
	log.Printf("🧮 Введён индекс: %d (input: %q)\n", index, input)
	if err != nil || index < 1 {
    return c.Send(messages.ErrInvalidGiftNumber)
	}

	var target interface{} = c.Chat().ID
	if len(parts) > 1 {
    targetStr := strings.TrimSpace(parts[1])
    if strings.HasPrefix(targetStr, "@") {
        target = targetStr
    } else {
        return c.Send("❌ Неверный формат канала. Укажите @channelname")
    }
	}


	gift, found := utils.GetGiftByIndex(userID, index-1)
	if !found {
		return c.Send(messages.ErrGiftNotFound)
	}
	log.Printf("📥 Из кэша: index=%d | ID=%q | Emoji=%s\n", index-1, gift.ID, gift.Sticker.Emoji)

	userstates.Clear(chatID)

	balance, err := db.GetBalance(context.Background(), userID)
	if err != nil {
		return c.Send(messages.ErrBalanceFetchFailed)
	}

	if balance < gift.StarCount {
		return c.Send(fmt.Sprintf("😢 У вас недостаточно звёзд. Необходимо %d ⭐️, у вас %d ⭐️", gift.StarCount, balance))
	}

	fmt.Printf("📦 Gift ID = %q, Emoji = %s, Цена = %d ⭐️\n", gift.ID, gift.Sticker.Emoji, gift.StarCount)

	return ExecuteGiftPurchaseAndSend(c, gift, balance, target)
}