package telegram

import (
	"errors"
	"fmt"

	"github.com/Onizukachi/telegram_raindrop/pkg/raindrop"
	"github.com/Onizukachi/telegram_raindrop/pkg/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot            *tgbotapi.BotAPI
	raindropClient *raindrop.Client
	userRepo       storage.UserRepository
}

func NewBot(bot *tgbotapi.BotAPI, raindropClient *raindrop.Client, userRepo storage.UserRepository) *Bot {
	return &Bot{bot: bot, raindropClient: raindropClient, userRepo: userRepo}
}

func (b *Bot) Run() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		user, err := b.userRepo.GetByChatID(update.Message.Chat.ID)
		if err != nil {
			if errors.Is(err, storage.ErrNoRecord) {
				authLink := b.raindropClient.BuildOAuthLink()
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, authLink)
				b.bot.Send(msg)
			} else {
				return err
			}
		}

		fmt.Println(user)
		// проверить что пользователь не с протухшим токеном и если что обновить токен

		if update.Message.IsCommand() {
			continue
		}

		if err := b.handleMessage(update.Message); err != nil {
			return err
		}
	}

	return nil
}
