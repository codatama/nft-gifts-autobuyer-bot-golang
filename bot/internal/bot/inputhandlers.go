package bot

import (
	"gopkg.in/telebot.v4"
	
	"context"
	"log"
	"strconv"
	"strings"
	"fmt"
	"time"

	"prvbot/internal/userstates"
	"prvbot/internal/db"
	"prvbot/internal/topup"
	"prvbot/internal/gifting"
	"prvbot/internal/refund"
)

func HandlePhotoInput(c telebot.Context) error {
	return HandleTextInput(c)
}

func HandleTextInput(c telebot.Context) error {
	log.Println("📥 Входящий текст от", c.Chat().ID, ":", c.Text())
	chatID := c.Chat().ID
	state := userstates.Get(chatID)
	log.Println("🔎 Состояние пользователя:", state)


	switch userstates.Get(chatID) {
	case "awaiting_topup":
		return topup.HandleTopUpInput(c)

	case "awaiting_gift_number":
		return gifting.HandleGiftInput(c)

	case "awaiting_ticket":
	userstates.Clear(chatID)

	adminIDs, err := db.GetSupportAdmins(context.Background())
	if err != nil {
		log.Printf("❌ Ошибка получения support-админов: %v", err)
		return c.Send("⚠️ Не удалось передать тикет поддержке. Попробуйте позже.")
	}

	for _, adminID := range adminIDs {
		_, err := c.Bot().Forward(&telebot.Chat{ID: adminID}, c.Message())
		if err != nil {
			log.Printf("[TicketForward] ❌ Ошибка при пересылке %d: %v", adminID, err)
		}
	}

	return c.Send("✅ Ваш тикет был отправлен технической поддержке. Ожидайте ответа.")

	case "awaiting_revoke_access_chat_id":
	userstates.Clear(chatID)

	chatIDStr := strings.TrimSpace(c.Text())
	targetID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		return c.Send("❌ Некорректный chat_id. Введите число.")
	}

	err = db.RevokeUserAccess(targetID)
if err != nil {
	switch err.Error() {
	case "не найден":
		return c.Send("🚫 Указанный chat_id не был найден в базе данных.")
	case "нельзя удалить админа":
		return c.Send("⛔ Вы не можете удалить администратора бота.")
	default:
		log.Printf("❗ Ошибка при отзыве доступа: %v", err)
		return c.Send("⚠️ Не удалось удалить пользователя. Проверь логи.")
	}
	}

	return c.Send(fmt.Sprintf("✅ Пользователь `%d` полностью удалён из базы и потерял доступ к боту.", targetID), &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	})

	case "awaiting_admin_revoke":
	userstates.Clear(chatID)
	chatIDStr := strings.TrimSpace(c.Text())
	targetID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		return c.Send("❌ Некорректный chat_id")
	}

	const protectedAdminID int64 = 6142264859
	if targetID == protectedAdminID {
	return c.Send("⛔ Вы не можете снять права у системного администратора.")
	}

	err = db.RevokePermission(context.Background(), targetID, "admin_access")
	if err != nil {
		return c.Send("❌ Не удалось снять права: " + err.Error())
	}

	return c.Send(fmt.Sprintf("📛 Права на админ-панель сняты у пользователя `%d`", targetID), &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	})

	case "awaiting_commission_value":
	userstates.Clear(chatID)

	text := strings.TrimSpace(c.Text())
	rate, err := strconv.ParseFloat(text, 64)
	if err != nil || rate < 0.02 || rate > 1 {
		return c.Send("❌ Введите корректную комиссию от 0.02 до 1 (например: 0.05)")
	}

	db.GlobalCommissionRate = rate
	return c.Send(fmt.Sprintf("✅ Комиссия успешно обновлена: *%.2f%%*", rate*100), &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	})

	case "awaiting_refund_revoke":
	userstates.Clear(chatID)
	chatIDStr := strings.TrimSpace(c.Text())
	targetID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		return c.Send("❌ Некорректный chat_id")
	}

	const protectedAdminID int64 = 6142264859
	if targetID == protectedAdminID {
	return c.Send("⛔ Вы не можете снять права у системного администратора.")
	}

	err = db.RevokePermission(context.Background(), targetID, "refund_access")
	if err != nil {
		return c.Send("❌ Не удалось снять права: " + err.Error())
	}

	return c.Send(fmt.Sprintf("📛 Права на рефаунд сняты у пользователя `%d`", targetID), &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	})

	case "awaiting_tech_revoke":
	userstates.Clear(chatID)
	chatIDStr := strings.TrimSpace(c.Text())
	targetID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		return c.Send("❌ Некорректный chat_id")
	}

	const protectedAdminID int64 = 6142264859
	if targetID == protectedAdminID {
	return c.Send("⛔ Вы не можете снять права у системного администратора.")
	}

	err = db.RevokePermission(context.Background(), targetID, "technical_support")
	if err != nil {
		return c.Send("❌ Не удалось снять права: " + err.Error())
	}

	return c.Send(fmt.Sprintf("📛 Права на техническую поддержку сняты у пользователя `%d`", targetID), &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	})

	case "awaiting_global_broadcast":
	userstates.Clear(chatID)

	ctx := context.Background()
	chatIDs, err := db.GetAllUsers(ctx)
	if err != nil {
		log.Printf("❌ Не удалось получить список пользователей: %v", err)
		return c.Send("⚠️ Не удалось выполнить рассылку.")
	}

	bot := c.Bot().(*telebot.Bot)

	if c.Message().Photo != nil {
		photo := c.Message().Photo
		caption := strings.TrimSpace(c.Message().Caption)

		go func() {
			for _, id := range chatIDs {
				recipient := &telebot.Chat{ID: id}
				newPhoto := &telebot.Photo{
					File:    photo.File,
					Caption: caption,
				}
				_, err := bot.Send(recipient, newPhoto)
				if err != nil {
					log.Printf("⚠️ Ошибка отправки фото %d: %v", id, err)
				}
				time.Sleep(30 * time.Millisecond)
			}
			log.Printf("✅ Фото + подпись отправлены %d пользователям.", len(chatIDs))
		}()

		return c.Send("✅ Изображение с описанием отправлено всем пользователям.")
	}

	text := strings.TrimSpace(c.Text())
	if text == "" {
		return c.Send("⚠️ Сообщение пустое. Попробуйте снова.")
	}

	go func() {
		for _, id := range chatIDs {
			recipient := &telebot.Chat{ID: id}
			_, err := bot.Send(recipient, text)
			if err != nil {
				log.Printf("⚠️ Ошибка отправки текста %d: %v", id, err)
			}
			time.Sleep(30 * time.Millisecond)
		}
		log.Printf("✅ Текстовое сообщение отправлено %d пользователям.", len(chatIDs))
	}()

	return c.Send("✅ Сообщение отправлено всем пользователям.")

	case "awaiting_techsupport_chat_id":
	userstates.Clear(chatID)

	chatIDStr := strings.TrimSpace(c.Text())
	newChatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		return c.Send("❌ Некорректный chat_id. Попробуйте снова.")
	}

	err = db.GrantPermission(context.Background(), newChatID, "technical_support")
	if err != nil {
		log.Printf("❗ Ошибка выдачи тех-доступа: %v", err)
		return c.Send("❌ Не удалось выдать доступ. Убедитесь, что такой пользователь существует.")
	}

	return c.Send(fmt.Sprintf("✅ Пользователь `%d` получил доступ к технической поддержке!", newChatID), &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	})

	case "awaiting_refund_chat_id":
	userstates.Clear(chatID)

	chatIDStr := strings.TrimSpace(c.Text())
	newChatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		return c.Send("❌ Некорректный chat_id. Попробуйте снова.")
	}

	err = db.GrantPermission(context.Background(), newChatID, "refund_access")
	if err != nil {
		log.Printf("❗ Ошибка выдачи прав на возврат: %v", err)
		return c.Send("❌ Не удалось выдать доступ. Убедитесь, что такой пользователь существует.")
	}

	return c.Send(fmt.Sprintf("✅ Пользователь `%d` получил доступ к возврату звёзд!", newChatID), &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	})

	case "awaiting_refund_input":
	args := strings.Fields(c.Text())
	if len(args) != 2 {
		return c.Send("❗ Использование: `user_id charge_id`", &telebot.SendOptions{ParseMode: telebot.ModeMarkdown})
	}

	userID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return c.Send("❌ Некорректный user_id. Попробуйте снова.")
	}

	chargeID := args[1]

	err = refund.TryRefundTransaction(context.Background(), c.Bot().(*telebot.Bot), userID, chargeID)
	if err != nil {
		return c.Send("❌ " + err.Error() + "\nПопробуйте ещё раз.")
	}

	userstates.Clear(chatID)
	return c.Send("✅ Возврат средств выполнен")

	case "awaiting_grant_chat_id":
	userstates.Clear(chatID)

	chatIDStr := strings.TrimSpace(c.Text())
	newChatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		return c.Send("❌ Некорректный chat_id. Попробуйте снова.")
	}

	err = db.CreateUserIfNotExists(newChatID)
	if err != nil {
		log.Printf("❗ Ошибка добавления пользователя: %v", err)
		return c.Send("⚠️ Не удалось добавить пользователя. Возможно, он уже есть.")
	}

	return c.Send(fmt.Sprintf("✅ Пользователь `%d` добавлен в базу данных!", newChatID), &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	})

	case "awaiting_channel_1", "awaiting_channel_2", "awaiting_channel_3":
	state := userstates.Get(chatID)               
	userstates.Clear(chatID)                      

	channelName := strings.TrimSpace(c.Text())
	if !strings.HasPrefix(channelName, "@") {
		return c.Send("❌ Название должно начинаться с `@`. Попробуйте снова.")
	}

	field := map[string]string{
		"awaiting_channel_1": "channel1",
		"awaiting_channel_2": "channel2",
		"awaiting_channel_3": "channel3",
	}[state]

	if field == "" {
		log.Printf("❌ Не удалось распознать поле: state='%s'", state)
		return c.Send("⚠️ Произошла ошибка. Попробуйте снова.")
	}

	err := db.UpdateChannelField(chatID, field, channelName)
	if err != nil {
		log.Printf("❌ Ошибка сохранения %s: %v", field, err)
		return c.Send("⚠️ Не удалось сохранить канал.")
	}

	return c.Send(fmt.Sprintf("✅ Канал %s успешно установлен: %s", field, channelName))

	case "awaiting_admin_chat_id":
	userstates.Clear(chatID)

	chatIDStr := strings.TrimSpace(c.Text())
	newChatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		return c.Send("❌ Некорректный chat_id. Попробуйте снова.")
	}

	err = db.GrantPermission(context.Background(), newChatID, "admin_access")
	if err != nil {
		log.Printf("❗ Ошибка обновления прав: %v", err)
		return c.Send("❌ Не удалось выдать доступ. Проверьте, существует ли такой пользователь.")
	}

	return c.Send(fmt.Sprintf("✅ Пользователь `%d` получил доступ к админ-панели!", newChatID), &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	})
	}

	return nil
}