package storage

import (
	"errors"
	"time"

	"github.com/Onizukachi/telegram_raindrop/internal/models"
)

var (
	ErrDuplicate    = errors.New("repository: record already exists")
	ErrNotExist     = errors.New("repository: row does not exist")
	ErrUpdateFailed = errors.New("repository: update failed")
	ErrDeleteFailed = errors.New("repository: delete failed")
)

type UserRepository interface {
	GetByChatID(chatID int64) (*models.User, error)
	All() ([]*models.User, error)
	Create(chatID int64, access, refresh string, expiresAt time.Time) error
	Update(chatID int64, access, refresh string, expiresAt time.Time) error
	Delete(chatID int64) error
}
