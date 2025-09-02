package models

import "time"

type User struct {
	ID           int64
	ChatID       int64
	AccessToken  string
	RefreshToken string
	ExpriresAt   time.Time
}
