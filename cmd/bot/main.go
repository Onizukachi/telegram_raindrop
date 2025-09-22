package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Onizukachi/telegram_raindrop/internal/config"
	"github.com/Onizukachi/telegram_raindrop/internal/logger"
	"github.com/Onizukachi/telegram_raindrop/internal/raindrop"
	"github.com/Onizukachi/telegram_raindrop/internal/server"
	"github.com/Onizukachi/telegram_raindrop/internal/storage"
	"github.com/Onizukachi/telegram_raindrop/internal/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("Error during loading config: %v", err)
	}

	logger := logger.SetupLogger(cfg.DebugMode)

	botApi, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		logger.Error("Error during create BotAPI", "error", err)
		os.Exit(1)
	}

	if cfg.DebugMode {
		botApi.Debug = true
	}

	db, err := initDB(cfg)
	if err != nil {
		logger.Error("Error during init DB", "error", err)
		os.Exit(1)
	}
	logger.Info("Database created successfully")

	defer db.Close()

	err = runMigrations(db)
	if err != nil {
		logger.Error("Error during apply migrations", "error", err)
		os.Exit(1)
	}
	logger.Info("Migrations applied successfully")

	useRepo := storage.NewPostgresUserRepo(db)
	raindropClient := raindrop.NewClient(cfg.Raindrop.ClientId, cfg.Raindrop.ClientSecret, cfg.Raindrop.RedirectUrl)
	bot := telegram.NewBot(botApi, raindropClient, useRepo, cfg.Messages, logger)
	server := server.NewServer(cfg.ServerAddr, cfg.BotName, botApi, raindropClient, useRepo, cfg.Messages)

	go func() {
		if err = server.Run(); err != nil {
			logger.Error("Error in server", "error", err)
			os.Exit(1)
		}
	}()

	if err := bot.Run(); err != nil {
		logger.Error("Error in bot", "error", err)
		os.Exit(1)
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

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create postgres driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///internal/storage/migrations",
		"postgres", driver)

	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %v", err)
	}

	return nil
}
