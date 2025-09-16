package bot

import (
	"gopkg.in/telebot.v4"
	
	"time"

	"prvbot/internal/userstates"
)

func HandleSendTicket(bot *telebot.Bot) func(c telebot.Context) error {
	return func(c telebot.Context) error {
		chatID := c.Chat().ID
		userstates.UserState[chatID] = "awaiting_ticket"

		msg, err := bot.Send(c.Recipient(), "Пожалуйста, введите свой вопрос и мы ответим на него в течении 24 часов. Возврат звёзд также происходит по этому запросу. ВАЖНО!!! Для полноценной обратной связи и получения ответа на интересующий вас вопрос, оставляйте @(ваше имя)")
		if err != nil {
			return err
		}

		go func() {
			time.Sleep(60 * time.Second)
			bot.Delete(msg)
		}()

		return nil
	}
}