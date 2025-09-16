package db

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Init(databaseURL string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	Pool, err = pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("❌ Ошибка подключения к базе данных через пул: %v", err)
	}

	log.Println("✅ Подключение к базе данных установлено через пул pgxpool")
}

func Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return Pool.Query(ctx, query, args...)
}

func Close() {
	if Pool != nil {
		Pool.Close()
		log.Println("🔒 Пул соединений закрыт")
	}
}