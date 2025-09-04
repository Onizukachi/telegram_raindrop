package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/Onizukachi/telegram_raindrop/pkg/config"
	"github.com/Onizukachi/telegram_raindrop/pkg/storage"
	"github.com/Onizukachi/telegram_raindrop/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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
		log.Fatalf("Error during create BotAPI: %v", err)
	}

	if cfg.DebugMode {
		botApi.Debug = true
	}

	db, err := initDB(cfg)
	if err != nil {
		log.Fatalf("Error during init DB: %v", err)
	}
	defer db.Close()

	useRepo := storage.NewPostgresUserRepo(db)
	bot := telegram.NewBot(botApi, useRepo, cfg.RedirectUrl)

	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "hello world")
		})

		fmt.Println("HTTP server is running on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	if err := bot.Run(); err != nil {
		log.Fatal(err)
	}
}

func initDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DbDSN)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
