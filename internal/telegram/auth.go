package telegram

import (
	"errors"
	"time"

	"github.com/Onizukachi/telegram_raindrop/internal/storage"
)

type AuthResult struct {
	User      *UserWithValidToken
	NeedsAuth bool
	AuthLink  string
}

type UserWithValidToken struct {
	ChatID      int64
	AccessToken string
}

func (b *Bot) AuthenticateUser(chatID int64) (*AuthResult, error) {
	user, err := b.userRepo.GetByChatID(chatID)
	if err != nil {
		if errors.Is(err, storage.ErrNotExist) {
			authLink := b.raindropClient.BuildOAuthLink(chatID)
			return &AuthResult{
				NeedsAuth: true,
				AuthLink:  authLink,
			}, nil
		} else {
			return nil, err
		}
	}

	if time.Now().After(user.ExpriresAt) {
		refreshResponse, err := b.raindropClient.RefreshToken(user.RefreshToken)
		if err != nil {
			return nil, err
		}

		expriresIn := time.Now().Add(time.Second * time.Duration(refreshResponse.ExpiresIn))
		err = b.userRepo.Update(chatID, refreshResponse.AccessToken, refreshResponse.RefreshToken, expriresIn)
		if err != nil {
			return nil, err
		}

		return &AuthResult{
			User: &UserWithValidToken{
				ChatID:      chatID,
				AccessToken: refreshResponse.AccessToken,
			},
		}, nil
	}

	return &AuthResult{
		User: &UserWithValidToken{
			ChatID:      chatID,
			AccessToken: user.AccessToken,
		},
	}, nil
}
