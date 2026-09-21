package domain

import (
	"context"
	"time"
)

type (
	Meal struct {
		ID            int64
		UserID        int64
		Name          string
		Ingredients   []Product
		CookedWeight  float64
		Servings      int64
		Protein       float64
		Fat           float64
		Carbohydrates float64
		Kcal          float64
		CreatedAt     time.Time
		UpdatedAt     time.Time
	}

	MealRepository interface {
		GetById(ctx context.Context, id int64) (Meal, error)
		GetAll(ctx context.Context, id int64) ([]Meal, error)
		Search(ctx context.Context, query string) ([]Meal, error)
		Create(ctx context.Context, ingredients Product) (Meal, error)
	}
)
