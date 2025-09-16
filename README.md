# 🧠 NFT Auto-Buyer Bot (Golang + Telegram API)

This Telegram bot was originally developed for commercial use to automate the purchase of NFT gifts within the Telegram messenger. It is written in Go and leverages Telegram Bot API, concurrent processing, and buffered channels for efficient task handling.

> ⚠️ **Disclaimer:**  
> This bot is **not intended as a ready-to-use solution**. It lacks built-in database migration logic — tables were created manually. If needed, you can generate them based on the models located in the `models` directory.

## 🚀 How It Works

- Each user’s NFT gift parameters are stored in the database.
- If the user has sufficient balance, the bot initiates a parsing loop every 10 seconds to check for available gifts.
- Auto-purchase is triggered for users with the feature enabled.
- The bot uses **worker pools**, **goroutines**, and **buffered channels** to handle concurrent purchases efficiently.

## 📊 Performance

- Used by over **2,000 users**  
- Managed a total balance of **200,000 Telegram Stars** (~2,000 USDT)  
- At peak load, handled simultaneous purchases for **200 users** without errors or downtime

## ⚙️ Setup

1. Clone the repository
2. Fill in your configuration in `bot/config/config.go`:
   - Telegram Bot Token (from BotFather)
   - Database connection details

## 🛠️ Tech Stack

- **Language:** Golang  
- **Database:** PostgreSQL  
- **Concurrency:** Goroutines, Channels  
- **Bot Framework:** Telegram Bot API

---

Feel free to explore the code, but use it with caution and adapt it to your own infrastructure.
