package buttons

import "gopkg.in/telebot.v4"

var InlineMenu = &telebot.ReplyMarkup{}

var (
	BtnTopUp   = InlineMenu.Data("💰 Пополнить баланс", "top_up")
	BtnShop    = InlineMenu.Data("🎁 Магазин подарков", "shop")
	BtnSendGift = InlineMenu.Data("🎁 Подарить", "open_gift")
	BtnPayStars = InlineMenu.Data("⭐️ Оплатить", "pay_stars")
	BtnBackToMenu = InlineMenu.Data("◀️ Вернуться", "go_back")
	BtnAutoBuyMenu = InlineMenu.Data("⚙️ Настройки автопокупки", "auto_buy")
	BtnUserProfile = InlineMenu.Data("👤 Ваш профиль", "user_profile")
	BtnTopOfBalance = InlineMenu.Data("🏅 Топ по балансу", "top_balance")
	BtnRefundStars = InlineMenu.Data("💫 Возврат звёзд", "stars_refund")
	BtnReferalLink = InlineMenu.URL("🏦 Покупка/Продажа NFT", "https://t.me/portals/market?startapp=8jv908")
	BtnChannelLink = InlineMenu.URL("📢 Канал бота", "https://t.me/autogift")
	BtnBuyingHistory = InlineMenu.Data("🔎 История покупок", "buying_history")
	BtnSupportUser = InlineMenu.URL("🫥 Помощь", "https://t.me/exotical")
)

func InitButtons() {
	InlineMenu.Inline(
		InlineMenu.Row(BtnAutoBuyMenu),
		InlineMenu.Row(BtnTopUp, BtnTopOfBalance),
		InlineMenu.Row(BtnShop, BtnBuyingHistory),
		InlineMenu.Row(BtnUserProfile),
		InlineMenu.Row(BtnReferalLink, BtnChannelLink),
		InlineMenu.Row(BtnSupportUser),
	)
}