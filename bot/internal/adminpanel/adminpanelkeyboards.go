package adminpanel

import (
	"gopkg.in/telebot.v4"
)

var AdminPanelInline = &telebot.ReplyMarkup{}

var (
	BtnAllowSubscription = AdminPanelInline.Data("✍️ Выдать доступ к боту", "allow_user")
	BtnDenySubscription = AdminPanelInline.Data("🔐 Забрать доступ к боту", "deny_subscription")
	BtnListOfTransactions = AdminPanelInline.Data("🧾 Список транзакций", "trans_list_admin")
	BtnGivePermissions = AdminPanelInline.Data("🌊 Выдать права на...", "give_permissions")
	BtnBackPermissions = AdminPanelInline.Data("📛 Забрать права на...", "back_permissions")
	BtnBroadcastMessage = AdminPanelInline.Data("📢 Отправить сообщение всем", "admin_broadcast")
	BtnSetComission = AdminPanelInline.Data("💳 Установить комиссию", "set_comission")
	BtnRefundAdmin = AdminPanelInline.Data("💫 Возврат звезд", "admin_refund")
)

var PermissionsInline = &telebot.ReplyMarkup{}

var (
	BtnGrantAdminPanel = PermissionsInline.Data("🛠 Права на админ панель", "grant_admin_access")
	BtnGrantRefundAccess = PermissionsInline.Data("💫 Права на возврат звёзд", "grant_refund_access")
	BtnGrantTechSupport = PermissionsInline.Data("🧰 Права на тех. поддержку", "grant_tech_support")
)

var BackPermissionsInline = & telebot.ReplyMarkup{}

var (
	BtnBackAdminPanel = BackPermissionsInline.Data("🛠 Забрать права на админ панель", "back_admin_access")
	BtnBackRefundAccess = BackPermissionsInline.Data("💫 Забрать права на возврат звёзд", "back_refund_access")
	BtnBackTechSupport = BackPermissionsInline.Data("🧰 Забрать права на тех. поддержку", "back_tech_support")
)