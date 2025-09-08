package telegram

import "errors"

var (
	invalidUrlError   = errors.New("url is invalid")
	unableToSaveError = errors.New("unable to save link to Raindrop")
)

// func (b *Bot) handleError(chatID int64, err error) {
// 	var messageText string

// 	switch err {
// 	case invalidUrlError:
// 		"url is invalid"
// 	}
// }
