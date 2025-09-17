package telegram

import "net/url"

func (b *Bot) validateURL(text string) error {
	_, err := url.ParseRequestURI(text)
	return err
}
