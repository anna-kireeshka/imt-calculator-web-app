package fitness

import (
	"app/imt-calculator-web-app/internal/domain"
	"fmt"
)

var macrosValues = map[string]map[string][2]float64{
	"fiber":         {"maintain": {1, 1.5}, "lose": {1.5, 2}, "gain": {1.8, 2.2}},
	"carbohydrates": {"maintain": {3, 4}, "lose": {2, 3}, "gain": {4, 6}},
	"fat":           {"maintain": {0.8, 1}, "lose": {0.8, 1.1}, "gain": {0.8, 1.2}},
}

func Macros(weight float64, goal string) (domain.Macros, error) {
	var result domain.Macros = make(map[string][2]float64)

	if goal == "" {
		return nil, fmt.Errorf("неизвестная цель: %q", goal)
	}
	for key, values := range macrosValues {
		for g, value := range values {
			if g == goal {
				value[0] *= weight
				value[1] *= weight
				result[key] = [2]float64{value[0], value[1]}
			}

		}
	}

	return result, nil
}
