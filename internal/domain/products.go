package domain

import (
	"context"
	"time"
)

type (
	Product struct {
		ID            int64
		Name          string
		Protein       float64
		Fat           float64
		Carbohydrates float64
		Kcal          float64
		CreatedAt     time.Time
	}

	ProductsRepository interface {
		GetById(ctx context.Context, id int64) (Product, error)
		GetAll(ctx context.Context, id int64) ([]Product, error)
		Search(ctx context.Context, query string) ([]Product, error)
		Create(ctx context.Context, product Product) (Product, error)
	}
)
