package database

import (
	"app/imt-calculator-web-app/internal/domain"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type userRepo struct {
	conn *pgx.Conn
}

func NewUserRepository(conn *pgx.Conn) domain.UserRepository {
	return &userRepo{
		conn: conn,
	}
}

func (r *userRepo) Save(ctx context.Context, userID int64, dob time.Time, gender bool) (domain.User, error) {
	var result domain.User

	if err := r.conn.QueryRow(ctx, `
		INSERT INTO users (id, telegram_id, dob, gender)
		VALUES ($1, $1, $2, $3)
		ON CONFLICT (id) DO UPDATE SET
			telegram_id = EXCLUDED.telegram_id,
			dob         = EXCLUDED.dob,
			gender      = EXCLUDED.gender,
			updated_at  = now()
		RETURNING id, telegram_id, COALESCE(name, ''), dob, gender, created_at, updated_at
	`, userID, dob, gender).
		Scan(&result.ID, &result.TelegramID, &result.Name, &result.DoB, &result.Gender,
			&result.CreatedAt, &result.UpdatedAt); err != nil {
		return domain.User{}, fmt.Errorf("save user: %w", err)
	}
	return result, nil
}

func (r *userRepo) GetByID(ctx context.Context, userID int64) (domain.User, error) {
	var result domain.User

	if err := r.conn.QueryRow(ctx, `
		SELECT id, COALESCE(telegram_id, id), COALESCE(name, ''), dob, gender, created_at, updated_at
		FROM users
		WHERE id = $1
	`, userID).
		Scan(&result.ID, &result.TelegramID, &result.Name, &result.DoB, &result.Gender,
			&result.CreatedAt, &result.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("%w: пользователь не заведён", domain.ErrNotFound)
		}
		return domain.User{}, fmt.Errorf("read user: %w", err)
	}

	return result, nil
}
