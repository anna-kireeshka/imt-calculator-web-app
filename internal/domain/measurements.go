package domain

import (
	"context"
	"time"
)

type (
	Measurement struct {
		ID           int64     `json:"id"`
		UserID       int64     `json:"user_id"`
		Height       int       `json:"height"`
		Weight       float64   `json:"weight"`
		BustSize     int       `json:"bust_size"`
		WaistSize    int       `json:"waist_size"`
		HipSize      int       `json:"hip_size"`
		ActivityType string    `json:"activity_type"`
		Goal         string    `json:"goal"`
		CreatedAt    time.Time `json:"created_at"`
	}

	FormulaData struct {
		UserID       int64
		Height       int
		Weight       float64
		ActivityType string
		Goal         string // 'lose','maintain','gain'
		DoB          time.Time
		Gender       bool
	}

	MeasurementsRepository interface {
		GetByID(ctx context.Context, userID int64) (Measurement, error)
		Save(ctx context.Context, userID int64, height int, weight float64, activityType string, goal string) (Measurement, error)
		Update(ctx context.Context, userID int64, m Measurement) (Measurement, error)
		GetDataForCalculation(ctx context.Context, userID int64) (FormulaData, error)
	}
)
