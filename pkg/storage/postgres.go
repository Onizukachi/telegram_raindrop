package storage

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Onizukachi/telegram_raindrop/pkg/models"
	"github.com/lib/pq"
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

	var user models.User

	err := row.Scan(&user.ID, &user.ChatID, &user.AccessToken, &user.RefreshToken, &user.ExpriresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExist
		} else {
			return nil, err
		}
	}

	return &user, nil
}

func (r *PostgresUserRepo) All() ([]*models.User, error) {
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

func (r *PostgresUserRepo) Create(chatID int64, access, refresh string, expiresAt time.Time) error {
	stmt := `INSERT INTO users (chat_id, access_token, refresh_token, expires_at, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, NOW(), NOW())`

	_, err := r.db.Exec(stmt, chatID, access, refresh, expiresAt)
	if err != nil {
		var pqError *pq.Error
		if errors.As(err, &pqError) {
			if pqError.Code == "23505" {
				return ErrDuplicate
			}
		}
		return err
	}

	return nil
}

func (r *PostgresUserRepo) Update(chatID int64, access, refresh string, expiresAt time.Time) error {
	stmt := `UPDATE users SET access_token = $1, refresh_token = $2, expires_at = $3, updated_at = NOW() WHERE chat_id = $4`
	res, err := r.db.Exec(stmt, access, refresh, expiresAt, chatID)
	if err != nil {
		var pqError *pq.Error
		if errors.As(err, &pqError) {
			if pqError.Code == "23505" {
				return ErrDuplicate
			}
		}
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUpdateFailed
	}

	return nil
}

func (r *PostgresUserRepo) Delete(chatID int64) error {
	stmt := `DELETE from users WHERE chat_id = $1`
	res, err := r.db.Exec(stmt, chatID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return nil
}
