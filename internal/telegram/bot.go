package telegram

import (
	"log/slog"

	"github.com/Onizukachi/telegram_raindrop/internal/config"
	"github.com/Onizukachi/telegram_raindrop/internal/raindrop"
	"github.com/Onizukachi/telegram_raindrop/internal/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot            *tgbotapi.BotAPI
	raindropClient *raindrop.Client
	userRepo       storage.UserRepository
	messages       config.Messages
	logger         *slog.Logger
}

func NewBot(bot *tgbotapi.BotAPI, raindropClient *raindrop.Client, userRepo storage.UserRepository, messages config.Messages, logger *slog.Logger) *Bot {
	return &Bot{bot: bot, raindropClient: raindropClient, userRepo: userRepo, messages: messages, logger: logger}
}

func (b *Bot) Run() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			if err := b.handleCommand(update.Message); err != nil {
				b.logger.Error("error when handle command", "error", err)
			}

			continue
		}

		if err := b.handleMessage(update.Message); err != nil {
			b.logger.Error("error when handle message", "error", err)
		}

		continue
	}

	return nil
}
