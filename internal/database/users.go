package database

import (
	"app/imt-calculator-web-app/internal/domain"
	"context"
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
		INSERT INTO users (id, dob, gender)
		VALUES ($1, $2, $3)
		RETURNING id, telegram_id, name, dob, gender, created_at, updated_at
	`, userID, dob, gender).
		Scan(&result.ID, &result.TelegramID, &result.Name, &result.DoB, &result.Gender,
			&result.CreatedAt, &result.UpdatedAt); err != nil {
		return domain.User{}, fmt.Errorf("commit tx: %w", err)
	}
	return result, nil
}
