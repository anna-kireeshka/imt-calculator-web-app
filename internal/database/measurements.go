package database

import (
	"app/imt-calculator-web-app/internal/domain"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type measurementsRepo struct {
	conn *pgx.Conn
}

func NewMeasurementsRepository(conn *pgx.Conn) domain.MeasurementsRepository {
	return &measurementsRepo{
		conn: conn,
	}
}

func (r *measurementsRepo) GetByID(ctx context.Context, userID int64) (domain.Measurement, error) {
	var result domain.Measurement

	if err := r.conn.QueryRow(ctx, `
		SELECT user_id, height, weight, created_at
		FROM measurements
		WHERE user_id = $1
		ORDER BY created_at DESC 
		LIMIT 1


	`, userID).Scan(&result.Height, &result.Weight); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Measurement{}, errors.New("у пользователя нет замеров")
		}
		return result, fmt.Errorf("read last measurement: %w", err)
	}

	return result, nil
}

func (r *measurementsRepo) Save(ctx context.Context, userID int64, height int, weight float64, activityType string, goal string) (domain.Measurement, error) {
	var result domain.Measurement

	if err := r.conn.QueryRow(ctx, `
		INSERT INTO measurements (user_id, height, weight, activity_type, goal)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, height, weight, bust_size, waist_size, hip_size, activity_type, goal, created_at
	`, userID, height, weight, activityType, goal).
		Scan(&result.ID, &result.UserID, &result.Height, &result.Weight, &result.BustSize, &result.WaistSize,
			&result.HipSize, &result.ActivityType, &result.Goal, &result.CreatedAt); err != nil {
		return result, fmt.Errorf("commit tx: %w", err)
	}

	return result, nil
}

func (r *measurementsRepo) GetHistory(ctx context.Context, userId int64) ([]domain.History, error) {
	var result []domain.History

	var rows, err = r.conn.Query(ctx, `
		SELECT height, weight, created_at
		FROM measurements
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var el domain.History
		if err := rows.Scan(&el.Height, &el.Weight, &el.CreatedAt); err != nil {
			return result, err
		}

		result = append(result, el)
	}

	if err = rows.Err(); err != nil {
		return result, err
	}

	return result, nil
}

func (r *measurementsRepo) GetDataForCalculation(ctx context.Context, userId int64) (domain.FormulaData, error) {
	var result domain.FormulaData

	if err := r.conn.QueryRow(ctx, `
		SELECT m.user_id, m.height, m.weight, m.activity_type, m.goal, u.dob, u.gender
		FROM measurements m
		INNER JOIN users u ON u.id = m.user_id
		ORDER BY u.created_at
	`, userId).
		Scan(&result.UserID, &result.Height, &result.Weight, &result.ActivityType, &result.Goal, &result.DoB, &result.Gender); err != nil {
		return result, fmt.Errorf("commit tx: %w", err)
	}
	return result, nil
}
