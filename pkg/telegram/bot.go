package telegram

import (
	"github.com/Onizukachi/telegram_raindrop/pkg/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot         *tgbotapi.BotAPI
	userRepo    storage.UserRepository
	redirectUrl string
}

func NewBot(bot *tgbotapi.BotAPI, userRepo storage.UserRepository, redirectUrl string) *Bot {
	return &Bot{bot: bot, userRepo: userRepo, redirectUrl: redirectUrl}
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
			continue
		}

		if err := b.handleMessage(update.Message); err != nil {
			return err
		}
	}

	return nil
}
