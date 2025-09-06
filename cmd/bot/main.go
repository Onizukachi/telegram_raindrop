package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/Onizukachi/telegram_raindrop/pkg/config"
	"github.com/Onizukachi/telegram_raindrop/pkg/raindrop"
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

	raindropClient := raindrop.NewClient(cfg.ClientId, cfg.ClientSecret, cfg.RedirectUrl)

	bot := telegram.NewBot(botApi, raindropClient, useRepo)

	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			botUrl := fmt.Sprintf("https://t.me/%s?start=code", cfg.BotName)

			code := r.URL.Query().Get("code")
			echangeResponse, err := raindropClient.ExchangeToken(code)
			if err != nil {
				log.Println(err)
				log.Println("ERRRRORRRR")
				http.Redirect(w, r, botUrl, http.StatusMovedPermanently)
				return
			}

			// &{AccessToken:af9edd9c-67e7-4101-bb07-5db226fa492b RefreshToken:a6fcf1c1-8073-4bee-8f04-ebe700d4812a ExpiresIn:1209599 TokenType:Bearer}
			log.Printf("%+v\n", echangeResponse)

			http.Redirect(w, r, botUrl, http.StatusMovedPermanently)
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
	db, err := sql.Open("postgres", cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
