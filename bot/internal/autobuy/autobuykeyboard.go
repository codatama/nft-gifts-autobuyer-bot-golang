package autobuy

import (
	"gopkg.in/telebot.v4"
)

var (
	BtnToggleAutoBuy = telebot.Btn{Unique: "autobuy_toggle"}
	BtnCycleCountChange = telebot.Btn{Unique: "autobuy_cycles_change"}
	BtnMinLimitChange = telebot.Btn{Unique: "autobuy_minlimit_change"}
	BtnMaxLimitChange = telebot.Btn{Unique: "autobuy_maxlimit_change"}
	BtnSupplyLimitChange = telebot.Btn{Unique: "autobuy_supplylimit_change"}
	BtnGoBackToStartMenu = telebot.Btn{Unique: "back_to_start_menu"}
	BtnChooseTheBuyer = telebot.Btn{Unique: "choose_the_buyer"}
	BtnChannelSettings = telebot.Btn{Unique: "channel_settings"}
	BtnGoBackToAutoBuyMenu = telebot.Btn{Unique: "back_to_auto_buy_menu"}
)