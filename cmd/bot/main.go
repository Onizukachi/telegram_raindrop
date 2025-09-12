package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Onizukachi/telegram_raindrop/pkg/config"
	"github.com/Onizukachi/telegram_raindrop/pkg/raindrop"
	"github.com/Onizukachi/telegram_raindrop/pkg/server"
	"github.com/Onizukachi/telegram_raindrop/pkg/storage"
	"github.com/Onizukachi/telegram_raindrop/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Init()
	fmt.Printf("%+v\n", cfg)

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
	raindropClient := raindrop.NewClient(cfg.Raindrop.ClientId, cfg.Raindrop.ClientSecret, cfg.Raindrop.RedirectUrl)
	bot := telegram.NewBot(botApi, raindropClient, useRepo)
	server := server.NewServer(cfg.ServerAddr, cfg.BotName, botApi, raindropClient, useRepo)

	go func() {
		if err = server.Run(); err != nil {
			log.Fatal(err)
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
