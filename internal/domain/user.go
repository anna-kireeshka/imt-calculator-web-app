package domain

import (
	"context"
	"time"
)

type (
	User struct {
		ID         int64     `json:"id"`
		TelegramID int64     `json:"telegram_id"`
		Name       string    `json:"name"`
		DoB        time.Time `json:"dob"`
		Gender     bool      `json:"gender"` // 1 female 0 male
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
	}

	UserRepository interface {
		Save(ctx context.Context, userID int64, dob time.Time, gender bool) (User, error)
		GetByID(ctx context.Context, userID int64) (User, error)
	}
)
