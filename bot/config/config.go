package config

type Config struct {
	TelegramToken string
	DatabaseURL   string
}

func Load() *Config {
	return &Config{
		TelegramToken: "your_token",
		DatabaseURL:   "postgres://username:password@localhost:5432/databasename?sslmode=disable",
	}
}