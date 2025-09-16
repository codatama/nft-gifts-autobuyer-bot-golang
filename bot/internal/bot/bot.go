package bot

import (
	"gopkg.in/telebot.v4"

	"log"
	"time"
	"strings"
	"strconv"
	"context"

	"prvbot/internal/topup"
	"prvbot/internal/tgapi"
	"prvbot/internal/buttons"
	"prvbot/internal/gifting"
	"prvbot/internal/autobuy"
	"prvbot/internal/adminpanel"
	"prvbot/internal/userstates"
	"prvbot/internal/db"
)


func New(token string) *telebot.Bot {
	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatalf("Ошибка создания бота: %s", err)
	}

	autobuy.InitTelegramDispatcher()

	userstates.Init(bot)

	if err := tgapi.SyncGiftsWithDatabase(bot); err != nil {
		log.Printf("⚠️ Ошибка при синхронизации подарков: %v", err)
	}

	tgapi.StartStarSyncScheduler()

	bot.Handle("/start", ProtectedStart)

	bot.Handle("/admns", adminpanel.SendAdminPanel)

	bot.Handle(&adminpanel.BtnListOfTransactions, func(c telebot.Context) error {
	return adminpanel.ShowTransactionsPage(c, 0)
	})

	bot.Handle(&buttons.BtnTopUp, topup.HandleTopUp)	

	bot.Handle(&autobuy.BtnGoBackToAutoBuyMenu, autobuy.HandleAutoBuySettings)

	bot.Handle(&buttons.BtnUserProfile, HandleUserProfile)

	bot.Handle(&buttons.BtnSendGift, gifting.HandleSendGift)

	bot.Handle(telebot.OnText, HandleTextInput)

	bot.Handle(telebot.OnPayment, topup.HandleSuccessfulPayment)

	bot.Handle(&buttons.BtnShop, gifting.HandleGiftShop)

	bot.Handle(&autobuy.BtnToggleAutoBuy, autobuy.HandleToggleAutoBuy)
	
	bot.Handle(&buttons.BtnAutoBuyMenu, autobuy.HandleAutoBuySettings)

	bot.Handle(&autobuy.BtnMinLimitChange, autobuy.HandleMinLimitChange)

	bot.Handle(&autobuy.BtnMaxLimitChange, autobuy.HandleMaxLimitChange)

	bot.Handle(&autobuy.BtnCycleCountChange, autobuy.HandleCyclesCount)

	bot.Handle(&autobuy.BtnSupplyLimitChange, autobuy.HandleSupplyLimit)

	bot.Handle(&buttons.BtnTopOfBalance, HandleTopOfBalance)

	bot.Handle(&autobuy.BtnChooseTheBuyer, autobuy.HandleChooseTheBuyer)

	bot.Handle(&adminpanel.BtnGrantTechSupport, adminpanel.GrantTechSupportAccess)

	bot.Handle(&adminpanel.BtnGrantRefundAccess, adminpanel.GrantRefundAccess)

	bot.Handle(&adminpanel.BtnGrantAdminPanel, adminpanel.GrantAdminAccess)

	bot.Handle(&adminpanel.BtnBackAdminPanel, adminpanel.RevokeAdminAccess)

	bot.Handle(&adminpanel.BtnBackRefundAccess, adminpanel.RevokeRefundAccess)

	bot.Handle(&adminpanel.BtnBackTechSupport, adminpanel.RevokeTechnicalSupport)

	bot.Handle(&adminpanel.BtnGivePermissions, adminpanel.GivePermissions)

	bot.Handle(&adminpanel.BtnBackPermissions, adminpanel.BackPermissions)

	bot.Handle(&adminpanel.BtnRefundAdmin, adminpanel.HandleRefundRequest)

	bot.Handle(&buttons.BtnRefundStars, HandleSendTicket(bot))

	bot.Handle(&adminpanel.BtnBroadcastMessage, adminpanel.HandleGlobalBroadcastRequest)

	bot.Handle(&buttons.BtnBuyingHistory, HandleBuyingHistory)

	bot.Handle(BtnGift1, HandleGift1)

	bot.Handle(BtnGift2, HandleGift2)

	bot.Handle(BtnGift3, HandleGift3)

	bot.Handle(telebot.OnPhoto, HandlePhotoInput)

	bot.Handle(telebot.OnCallback, func(c telebot.Context) error {
	data := c.Callback().Data
	cleanData := strings.TrimSpace(data)
	cleanData = strings.Trim(cleanData, "\f\r\n")
	log.Println("📨 Получен callback:", c.Callback().Data)
	switch {
	case cleanData == "channel_settings":
	return autobuy.HandleChannelSettings(c)

	case cleanData == "autobuy:add_channel":
	return autobuy.HandleAddChannelStepOne(c)

	case cleanData == "target:self", cleanData == "target:channel":
	return autobuy.HandleChooseTheBuyerCallback(c)

	case cleanData == "allow_user":
	userstates.Set(c.Chat().ID, "awaiting_grant_chat_id")
	return c.Send("✍️ Введите `chat_id` для добавления пользователя в базу:", &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	})

	case cleanData == "add_channel":
	return autobuy.HandleAddChannelStepOne(c)

	case cleanData == "channel:edit1":
	return autobuy.PromptChannelInput(c, "channel1")
	case cleanData == "channel:edit2":
	return autobuy.PromptChannelInput(c, "channel2")
	case cleanData == "channel:edit3":
	return autobuy.PromptChannelInput(c, "channel3")
	
	case cleanData == "set_comission":
	userstates.Set(c.Chat().ID, "awaiting_commission_value")
	return adminpanel.HandleSetCommissionPrompt(c)
	
	case cleanData == "admin_broadcast":
	return adminpanel.HandleGlobalBroadcastRequest(c)

	case cleanData == "admin_refund":
	return adminpanel.HandleRefundRequest(c)

	case cleanData == "deny_subscription":
	userstates.Set(c.Chat().ID, "awaiting_revoke_access_chat_id")
	return c.Send("✂️ Введите `chat_id`, которому нужно _закрыть_ доступ к боту:", &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	})

	case cleanData == "admin_back":
		return adminpanel.EditAdminPanel(c)

	case strings.HasPrefix(cleanData, "txpage:"):
		pageStr := strings.TrimPrefix(cleanData, "txpage:")
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			return c.Respond(&telebot.CallbackResponse{Text: "❌ Неверный номер страницы"})
		}
		return adminpanel.ShowTransactionsPage(c, page)
	}

	prefixes := map[string]func(telebot.Context, int) error{
		"min:":    autobuy.HandleMinLimitCallback,
		"max:":    autobuy.HandleMaxLimitCallback,
		"cycles:": autobuy.HandleCyclesCountCallback,
		"supply:": autobuy.HandleSupplyLimitCallback,
	}

	for prefix, handler := range prefixes {
		if strings.HasPrefix(cleanData, prefix) {
			valueStr := strings.TrimPrefix(cleanData, prefix)
			value, err := strconv.Atoi(valueStr)
			if err != nil {
				return c.Respond(&telebot.CallbackResponse{Text: "❌ Неверное значение"})
			}
			return handler(c, value)
		}
	}

	return c.Respond()
})

	bot.Handle(&autobuy.BtnGoBackToStartMenu, func(c telebot.Context) error {
	text, markup := tgbStartMessage()
	return c.Edit(text, markup)
	}) 
	
	bot.Handle(&buttons.BtnBackToMenu, func(c telebot.Context) error {
	text, markup := tgbStartMessage()
	return c.Edit(text, markup)
	})

	bot.Handle(telebot.OnCheckout, func(c telebot.Context) error {
	query := c.PreCheckoutQuery()
	if query == nil {
		return nil
	}

	_, err := c.Bot().Raw("answerPreCheckoutQuery", map[string]interface{}{
		"pre_checkout_query_id": query.ID,
		"ok": true,
	})
	return err
})

go func() {
	ctx := context.Background()

	if err := db.CaptureGlobalBalanceSnapshot(ctx); err != nil {
		log.Printf("⛔ Ошибка сохранения общего баланса: %v", err)
	}

	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		<-ticker.C
		if err := db.CaptureGlobalBalanceSnapshot(ctx); err != nil {
			log.Printf("⛔ Ошибка автоснимка: %v", err)
		}
	}
}()

bot.SetCommands([]telebot.Command{
	{
		Text:        "start",
		Description: "Главное меню",
	},
})

	go autobuy.AutoBuyTick(context.Background(), bot)
	return bot
}