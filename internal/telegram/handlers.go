package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	commandStart = "start"
)

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	chatID := message.Chat.ID
	authResult, err := b.AuthenticateUser(chatID)
	if err != nil {
		return b.sendMessage(chatID, b.messages.Errors.FailAuth, message.MessageID)
	}

	if authResult.NeedsAuth {
		textMsg := fmt.Sprintf(b.messages.Start, authResult.AuthLink)
		return b.sendMessage(chatID, textMsg, message.MessageID)
	}

	if err := b.validateURL(message.Text); err != nil {
		return b.sendMessage(chatID, b.messages.InvalidURL, message.MessageID)
	}

	err = b.raindropClient.CreateItem(message.Text, authResult.User.AccessToken)
	if err != nil {
		return b.sendMessage(chatID, b.messages.UnableToSave, message.MessageID)
	}

	return b.sendMessage(chatID, b.messages.Responses.LinkSaved, message.MessageID)
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
	authResult, err := b.AuthenticateUser(chatID)
	if err != nil {
		return b.sendMessage(chatID, b.messages.Errors.FailAuth)
	}

	if authResult.NeedsAuth {
		textMsg := fmt.Sprintf(b.messages.Start, authResult.AuthLink)
		return b.sendMessage(chatID, textMsg)
	}

	return b.sendMessage(chatID, b.messages.Responses.AlreadyAuthorized)
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	return b.sendMessage(message.Chat.ID, b.messages.Responses.UnknownCommand)
}

func (b *Bot) sendMessage(chatID int64, message string, replyTo ...int) error {
	msg := tgbotapi.NewMessage(chatID, message)
	if len(replyTo) > 0 && replyTo[0] > 0 {
		msg.ReplyToMessageID = replyTo[0]
	}
	_, err := b.bot.Send(msg)
	return err
}
