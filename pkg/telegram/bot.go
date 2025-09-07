package telegram

import (
	"errors"
	"fmt"
	"time"

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

		// Доработать тут логику наверно надо какойто сервис аунтентификации чтоб вытащить отсюла логику
		// а уже в хенлдер сообщений передавать если все норм
		user, err := b.userRepo.GetByChatID(update.Message.Chat.ID)
		if err != nil {
			if errors.Is(err, storage.ErrNotExist) {
				authLink := b.raindropClient.BuildOAuthLink(update.Message.Chat.ID)
				textMsg := fmt.Sprintf("Необходимо авторизоваться в Raindrop по данной ссылке: %s", authLink)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, textMsg)
				b.bot.Send(msg)
				continue
			} else {
				return err
			}
		}

		if time.Now().Before(user.ExpriresAt) {
			refreshResponse, err := b.raindropClient.RefreshToken(user.RefreshToken)
			if err != nil {
				textMsg := "Не получилось обновить токен"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, textMsg)
				b.bot.Send(msg)
				continue
			}

			expriresIn := time.Now().Add(time.Second * time.Duration(refreshResponse.ExpiresIn))
			err = b.userRepo.Update(user.ChatID, refreshResponse.AccessToken, refreshResponse.RefreshToken, expriresIn)
			if err != nil {
				textMsg := "Не получилось обновить токен"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, textMsg)
				b.bot.Send(msg)
				continue
			}

		}

		if update.Message.IsCommand() {
			continue
		}

		if err := b.handleMessage(update.Message, user); err != nil {
			return err
		}
	}

	return nil
}
