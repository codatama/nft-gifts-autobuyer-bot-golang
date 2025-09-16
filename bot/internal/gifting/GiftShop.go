package gifting

import (
	"gopkg.in/telebot.v4"

	"fmt"

	"prvbot/internal/buttons"
	"prvbot/internal/messages"
	"prvbot/internal/tgapi"
	"prvbot/internal/utils"
)

func HandleGiftShop(c telebot.Context) error {
	b := c.Bot().(*telebot.Bot)

	var gifts []tgapi.TelegramGift
	var err error

	gifts, ok := utils.GetFreshGiftList(c.Sender().ID)
	if !ok {
		gifts, err = tgapi.GetAvailableGifts(b)
		if err != nil {
			return c.Edit(messages.ErrGiftsLoadFailed)
		}
		if len(gifts) == 0 {
			return c.Edit(messages.NoGiftsAvailable)
		}
		utils.SaveGiftList(c.Sender().ID, gifts)
	}

	var text string
	for i, gift := range gifts {
		text += fmt.Sprintf("%d. *%s* — %d ⭐️\n\n", i+1, gift.Sticker.Emoji, gift.StarCount)
	}

	markup := &telebot.ReplyMarkup{}
	markup.Inline(
		markup.Row(buttons.BtnSendGift),
		markup.Row(buttons.BtnBackToMenu),
	)

	return c.Edit(text, &telebot.SendOptions{
		ParseMode:   telebot.ModeMarkdown,
		ReplyMarkup: markup,
	})
}