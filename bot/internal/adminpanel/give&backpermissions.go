package adminpanel

import (
	"gopkg.in/telebot.v4"
)

func GivePermissions(c telebot.Context) error {
	InitPermissionsButtons()
	return c.Edit("🧾 *Выберите, на что вы хотите выдать права:*", PermissionsInline, telebot.ModeMarkdown)
}

func BackPermissions(c telebot.Context) error {
	InitBackPermissionsButtons()
	return c.Edit("🧾 *Выберите, на что вы хотите забрать права:*", BackPermissionsInline, telebot.ModeMarkdown)
}