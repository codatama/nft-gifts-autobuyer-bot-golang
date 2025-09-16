package tgapi

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"encoding/json"

	"prvbot/config"
)

func SendStarsInvoice(chatID int64, starsAmount float64) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendInvoice", config.Load().TelegramToken)

	replyMarkup := map[string]interface{}{
		"inline_keyboard": [][]map[string]interface{}{
			{
				{
					"text": "💳 Оплатить",
					"pay":  true,
				},
			},
		},
	}

replyMarkupBytes, _ := json.Marshal(replyMarkup)


	form := url.Values{}
	form.Set("chat_id", fmt.Sprintf("%d", chatID))
	form.Set("title", "Пополнение звёзд ✨")
	form.Set("description", fmt.Sprintf("Вы пополняете баланс на %.2f звёзд", starsAmount))
	form.Set("payload", fmt.Sprintf("topup_%.2f", starsAmount))
	form.Set("provider_token", "")
	form.Set("currency", "XTR")
	form.Set("prices", fmt.Sprintf(`[{"label":"Звёзды","amount":%d}]`, int(starsAmount)))
	form.Set("reply_markup", string(replyMarkupBytes))

	resp, err := http.Post(
		apiURL,
		"application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Ответ Telegram:", string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("не удалось отправить invoice: код %d, ответ: %s", resp.StatusCode, string(body))
	}

	var result struct {
	OK     bool `json:"ok"`
	Result struct {
		MessageID int `json:"message_id"`
		Chat      struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"result"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}
	return nil
}