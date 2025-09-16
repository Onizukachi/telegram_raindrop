package telegram

import (
	"errors"
	"fmt"
	"time"

	"github.com/Onizukachi/telegram_raindrop/internal/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	commandStart = "start"
)

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	// Проверить что пользователь существует и что токен валидный
	// если не существует то отправить ссылку на авторизацию
	// если истек срок попробовать рефрешить
	// если все норм то сохраняем
	// err := b.raindropClient.CreateItem(message.Text, user.AccessToken)
	// if err != nil {
	// 	return err
	// }

	// textMsg := "Ссылка успешно сохранена :)"
	// msg := tgbotapi.NewMessage(message.Chat.ID, textMsg)
	// msg.ReplyToMessageID = message.MessageID
	// b.bot.Send(msg)

	return nil
}

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	chatID := message.Chat.ID
	user, err := b.userRepo.GetByChatID(chatID)
	if err != nil {
		if errors.Is(err, storage.ErrNotExist) {
			authLink := b.raindropClient.BuildOAuthLink(chatID)
			textMsg := fmt.Sprintf(b.messages.Start, authLink)
			msg := tgbotapi.NewMessage(chatID, textMsg)
			b.bot.Send(msg)
			return nil
		} else {
			msg := tgbotapi.NewMessage(chatID, b.messages.Errors.Default)
			b.bot.Send(msg)
			return err
		}
	}

	if time.Now().Before(user.ExpriresAt) {
		refreshResponse, err := b.raindropClient.RefreshToken(user.RefreshToken)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, b.messages.Errors.FailAuth)
			b.bot.Send(msg)
			return err
		}

		expriresIn := time.Now().Add(time.Second * time.Duration(refreshResponse.ExpiresIn))
		err = b.userRepo.Update(chatID, refreshResponse.AccessToken, refreshResponse.RefreshToken, expriresIn)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, b.messages.Errors.Default)
			b.bot.Send(msg)
			return err
		}
	} else {
		msg := tgbotapi.NewMessage(chatID, b.messages.Responses.AlreadyAuthorized)
		b.bot.Send(msg)
		return nil
	}

	return nil
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Responses.UnknownCommand)
	_, err := b.bot.Send(msg)

	return err
}
