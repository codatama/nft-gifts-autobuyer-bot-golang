package autobuy

import (
	"gopkg.in/telebot.v4"
	
	"prvbot/internal/models"
)

func RenderAutoBuyKeyboard(user *models.User) *telebot.ReplyMarkup {
	markup := &telebot.ReplyMarkup{}

	markup.Inline(
		markup.Row(markup.Data("🔁 Изменить режим автопокупки", BtnToggleAutoBuy.Unique)),
		markup.Row(
			markup.Data("🪛 Мин. лимит", BtnMinLimitChange.Unique),
			markup.Data("🪛 Макс. лимит", BtnMaxLimitChange.Unique),
		),
		markup.Row(
			markup.Data("🪛 Кол-во циклов", BtnCycleCountChange.Unique),
		),
		markup.Row(
			markup.Data("🪛 Сапплай лимит", BtnSupplyLimitChange.Unique),
		),
		markup.Row(
			markup.Data("📩 Выбрать получателя", BtnChooseTheBuyer.Unique),
			markup.Data(" Настройки каналов", BtnChannelSettings.Unique),
		),
		markup.Row(
			markup.Data("◀️ Вернуться назад", BtnGoBackToStartMenu.Unique),
		),
	)

	return markup
}