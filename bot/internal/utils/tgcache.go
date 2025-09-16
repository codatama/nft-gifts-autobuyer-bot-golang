package utils

import (
	"log"
	"sync"
	"prvbot/internal/tgapi"
	"time"
)

type GiftCacheEntry struct {
	Gifts     []tgapi.TelegramGift
	UpdatedAt int64
}

var giftCache = make(map[int64]GiftCacheEntry)
var mu sync.RWMutex

func SaveGiftList(userID int64, gifts []tgapi.TelegramGift) {
	mu.Lock()
	defer mu.Unlock()
	giftCache[userID] = GiftCacheEntry{
		Gifts:     gifts,
		UpdatedAt: time.Now().Unix(),
	}
}

func GetGiftByIndex(userID int64, index int) (*tgapi.TelegramGift, bool) {
	log.Printf("📦 GetGiftByIndex: userID=%d, index=%d\n", userID, index)
	mu.RLock()
	defer mu.RUnlock()

	cache, ok := giftCache[userID]
	if !ok || index < 0 || index >= len(cache.Gifts) {
		return nil, false
	}

	return &cache.Gifts[index], true
}

func GetFreshGiftList(userID int64) ([]tgapi.TelegramGift, bool) {
	mu.RLock()
	defer mu.RUnlock()

	cache, ok := giftCache[userID]
	if !ok {
		return nil, false
	}

	if time.Since(time.Unix(cache.UpdatedAt, 0)) > 24*time.Hour {
		return nil, false
	}

	return cache.Gifts, true
}