package storage

import (
	"database/sql"
	"errors"
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
	stmt := "SELECT id, chat_id, access_token, refresh_token, expires_at FROM users WHERE chat_id = $1"
	row := r.db.QueryRow(stmt, chatID)

	user := &models.User{}
	err := row.Scan(&user.ID, &user.ChatID, &user.AccessToken, &user.RefreshToken, &user.ExpriresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return user, nil
}

func (r *PostgresUserRepo) GetAll() ([]*models.User, error) {
	stmt := "SELECT id, chat_id, access_token, refresh_token, expires_at FROM users"
	rows, err := r.db.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []*models.User{}

	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(&user.ID, &user.ChatID, &user.AccessToken, &user.RefreshToken, &user.ExpriresAt); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *PostgresUserRepo) Create(chatID int64, access, refresh string, expiresAt time.Time) (int64, error) {
	stmt := `INSERT INTO users (chat_id, access_token, refresh_token, expires_at, created_at, updated_at) 
		VALUES (?, ?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`

	result, err := r.db.Exec(stmt, chatID, access, refresh, expiresAt)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
