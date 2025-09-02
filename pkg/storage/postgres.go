package storage

import (
	"database/sql"
	"time"

	"github.com/Onizukachi/telegram_raindrop/pkg/models"
)

type PostgresUserRepo struct {
	db *sql.DB
}

func NewPostgresUserRepo(db *sql.DB) *PostgresUserRepo {
	return &PostgresUserRepo{db: db}
}

func (r *PostgresUserRepo) GetByChatID(chatID int64) (*models.User, error) {
	var user models.User

	err := r.db.QueryRow("SELECT id, chat_id, access_token, refresh_token, expires_at FROM users WHERE chat_id = $1", chatID).
		Scan(&user.ID, &user.ChatID, &user.AccessToken, &user.RefreshToken, &user.ExpriresAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepo) GetAll() ([]*models.User, error) {
	return nil, nil
}

func (r *PostgresUserRepo) Create(id int64, chatID int64, access, refresh string, expiresAt time.Time) error {
	return nil
}
