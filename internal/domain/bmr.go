package domain

import "time"

type (
	// Macros — нормы по нутриентам: для каждого нижняя и верхняя граница в граммах.
	Macros map[string][2]float64

	BMR struct {
		ID        int64     `json:"id"`
		IMT       float64   `json:"imt"`
		Age       int       `json:"age"`
		Ccal      int       `json:"ccal"`
		Macros    Macros    `json:"macros"`
		CreatedAt time.Time `json:"created_at"`
	}
)
