package autobuy

import (
	"fmt"
	
	"prvbot/internal/models"
)

func RenderAutoBuySettings(user *models.User) string {
	status := "🔴 Выключено"
	if user.AutoBuy {
		status = "🟢 Включено"
	}

	return fmt.Sprintf(`⚙️ *Настройки автопокупки*

Статус: %s

Лимит цены:
От %d до %d ⭐️

Количество циклов: %d

Сапплай: %d`,
		status, user.MinCostLimit, user.MaxCostLimit, user.CyclesCount, user.SupplyLimit)
}