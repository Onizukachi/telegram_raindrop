package telegram

import (
	"github.com/Onizukachi/telegram_raindrop/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleMessage(message *tgbotapi.Message, user *models.User) error {
	err := b.raindropClient.CreateItem(message.Text, user.AccessToken)
	if err != nil {
		return err
	}

	textMsg := "Ссылка успешно сохранена :)"
	msg := tgbotapi.NewMessage(message.Chat.ID, textMsg)
	msg.ReplyToMessageID = message.MessageID
	b.bot.Send(msg)

	return nil
}
