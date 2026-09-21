package domain

import "time"

type (
	BMR struct {
		ID        int64
		IMT       float64
		Ccal      int
		Macros    Macros // macros: {Fiber: [],Fat: [], Carbohydrates[] }
		CreatedAt time.Time
	}

	Macros map[string][2]float64
)
