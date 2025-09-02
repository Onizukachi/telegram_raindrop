package storage

import (
	"time"

	"github.com/Onizukachi/telegram_raindrop/pkg/models"
)

type UserRepository interface {
	GetByChatID(chatID int64) (*models.User, error)
	GetAll() ([]*models.User, error)
	Create(id int64, chatID int64, access, refresh string, expiresAt time.Time) error
}
