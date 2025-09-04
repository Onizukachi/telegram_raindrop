package storage

import (
	"errors"
	"time"

	"github.com/Onizukachi/telegram_raindrop/pkg/models"
)

var (
	ErrNoRecord = errors.New("repository: no matching record")
)

type UserRepository interface {
	GetByChatID(chatID int64) (*models.User, error)
	GetAll() ([]*models.User, error)
	Create(chatID int64, access, refresh string, expiresAt time.Time) (int64, error)
}
