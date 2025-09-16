package autobuy

import (
	"gopkg.in/telebot.v4"

	"time"
	"context"
	"log"
	"sort"
	"strconv"
	
	"prvbot/internal/db"
	"prvbot/internal/tgapi"
	"prvbot/internal/models"
)

var TelegramQueue = make(chan func(), 100)

const TelegramWorkerCount = 5

func InitTelegramDispatcher() {
	for i := 0; i < TelegramWorkerCount; i++ {
		go func(id int) {
			for task := range TelegramQueue {
				log.Printf("📮 Worker %d выполняет задачу", id)
				task()
				time.Sleep(100 * time.Millisecond) // защитная задержка
			}
		}(i + 1)
	}
}

var ignoredGiftIDs = map[string]struct{}{
	"5170145012310081615": {},
	"5170233102089322756": {},
	"5170250947678437525": {},
	"5168103777563050263": {},
	"5170144170496491616": {},
	"5170314324215857265": {},
	"5170564780938756245": {},
	"5168043875654172773": {},
	"5170690322832818290": {},
	"5170521118301225164": {},
	"6028601630662853006": {},
}

func AutoBuyTick(ctx context.Context, bot *telebot.Bot) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("[AUTO-BUY] ⏹ Остановка автообновления")
			return

		case <-ticker.C:
			log.Println("[AUTO-BUY] 🔄 Новый тик: проверка подарков")

			gifts, err := tgapi.GetAvailableGifts(bot)
			if err != nil {
				log.Printf("[AUTO-BUY] ❌ Не удалось получить подарки: %v", err)
				continue
			}

			log.Printf("[AUTO-BUY] 🎁 Получено %d подарков", len(gifts))
			for _, g := range gifts {
				log.Printf("[AUTO-BUY]    • %s — %d⭐️", g.ID, g.StarCount)
			}

			var filteredGifts []tgapi.TelegramGift
			for _, g := range gifts {
				if _, banned := ignoredGiftIDs[g.ID]; banned {
					log.Printf("[AUTO-BUY] ⚠️ Подарок %s исключён глобально (игнор)", g.ID)
					continue
				}
				filteredGifts = append(filteredGifts, g)
			}

			users, err := db.GetUsersWithAutoBuy(ctx)
			if err != nil {
				log.Printf("[AUTO-BUY] ❌ Не удалось получить пользователей: %v", err)
				continue
			}

			log.Printf("[AUTO-BUY] 👥 Найдено %d пользователей с включённым автобаeм", len(users))
			for _, u := range users {
				log.Printf("[AUTO-BUY]    → Пользователь %d (лимит %d–%d⭐️, баланс %d⭐️)", u.ID, u.MinCostLimit, u.MaxCostLimit, u.Balance)
			}

			for _, user := range users {
				localUser := user
				localGifts := append([]tgapi.TelegramGift(nil), filteredGifts...)
				TelegramQueue <- func() {
					processAutoBuyForUser(ctx, bot, localUser, localGifts)
				}
			}
		}
	}
}

func processAutoBuyForUser(ctx context.Context, bot *telebot.Bot, user models.User, gifts []tgapi.TelegramGift) {
	var recipient string
	var purchaseMode string
	var pc *models.PurchaseChannels

	if !user.ChannelEnabled {
		recipient = strconv.FormatInt(user.ChatID, 10)
		purchaseMode = "direct"
		log.Printf("[AUTO-BUY] 👤 %d → режим личной покупки", user.ID)
	} else {
		var err error
		pc, err = db.GetOrCreatePurchaseChannels(ctx, user.ChatID)
		if err != nil || (pc.Channel1 == "" && pc.Channel2 == "" && pc.Channel3 == "") {
			recipient = strconv.FormatInt(user.ChatID, 10)
			purchaseMode = "direct"
			log.Printf("[AUTO-BUY] ⚠️ %d → каналы не заданы, fallback на личную покупку", user.ID)
		} else {
			purchaseMode = "cascaded"
			log.Printf("[AUTO-BUY] 📡 %d → каскадная покупка через каналы", user.ID)
		}
	}

	var suitable []struct {
		Gift      tgapi.TelegramGift
		Recipient string
	}
	for _, g := range gifts {
		if _, banned := ignoredGiftIDs[g.ID]; banned || g.TotalCount == nil {
			continue
		}
		if g.StarCount < user.MinCostLimit || g.StarCount > user.MaxCostLimit {
			continue
		}
		if user.SupplyLimit > 0 && *g.TotalCount > user.SupplyLimit {
			continue
		}

		target := recipient
		if purchaseMode == "cascaded" {
			count := *g.TotalCount
			switch {
			case count <= 15000 && pc.Channel1 != "":
				target = pc.Channel1
			case count > 15000 && count <= 50000 && pc.Channel2 != "":
				target = pc.Channel2
			case count > 50000 && count <= 1000000 && pc.Channel3!= "":
				target = pc.Channel3
			default:
				continue
			}
		}

		suitable = append(suitable, struct {
			Gift      tgapi.TelegramGift
			Recipient string
		}{g, target})
	}

	if len(suitable) == 0 {
		log.Printf("[AUTO-BUY] ❌ %d → нет подходящих подарков", user.ID)
		return
	}

	sort.Slice(suitable, func(i, j int) bool {
	priceI := suitable[i].Gift.StarCount
	priceJ := suitable[j].Gift.StarCount

	supplyI := 999999
	supplyJ := 999999
	if suitable[i].Gift.TotalCount != nil {
		supplyI = *suitable[i].Gift.TotalCount
	}
	if suitable[j].Gift.TotalCount != nil {
		supplyJ = *suitable[j].Gift.TotalCount
	}

	if priceI == priceJ {
		return supplyI < supplyJ
	}
	return priceI > priceJ
	})

	for _, item := range suitable {
		if user.Balance < item.Gift.StarCount {
			log.Printf("[AUTO-BUY] 💸 %d → недостаточно звёзд для '%s' (%d⭐️), баланс: %d⭐️", user.ID, item.Gift.ID, item.Gift.StarCount, user.Balance)
			continue
		}

		localUser := user
		localGift := item.Gift
		localRecipient := item.Recipient

		TelegramQueue <- func() {
			err := AutoBuyGift(ctx, bot, localUser.ChatID, localRecipient, localUser.Balance, localGift)
			if err != nil {
				log.Printf("[AUTO-BUY] ❌ %d → ошибка покупки '%s': %v", localUser.ID, localGift.ID, err)
				return
			}
			log.Printf("[AUTO-BUY] ✅ %d → куплен '%s' (%d⭐️) для %s", localUser.ID, localGift.ID, localGift.StarCount, localRecipient)
		}

		user.Balance -= localGift.StarCount
	}
}