package main

import (
	"database/sql"
	"log"

	"github.com/Onizukachi/telegram_raindrop/pkg/config"
	"github.com/Onizukachi/telegram_raindrop/pkg/storage"
	"github.com/Onizukachi/telegram_raindrop/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("Error during loading config: %v", err)
	}

	botApi, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.DebugMode {
		botApi.Debug = true
	}

	db, err := sql.Open("postgres", "postgres://hikaru:potolok149@localhost/mydb?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	useRepo := storage.NewPostgresUserRepo(db)

	bot := telegram.NewBot(botApi, cfg.RedirectUrl)
	if err := bot.Run(); err != nil {
		log.Fatal(err)
	}
}
