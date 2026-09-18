package fitness

import (
	"fmt"

	"app/imt-calculator-web-app/internal/domain"
)

// Нормы на килограмм веса: нижняя и верхняя граница для каждой цели.
var macrosValues = map[string]map[string][2]float64{
	"protein":       {"maintain": {1, 1.5}, "lose": {1.5, 2}, "gain": {1.8, 2.2}},
	"carbohydrates": {"maintain": {3, 4}, "lose": {2, 3}, "gain": {4, 6}},
	"fat":           {"maintain": {0.8, 1}, "lose": {0.8, 1.1}, "gain": {0.8, 1.2}},
}

// Macros переводит нормы на килограмм в граммы для конкретного веса.
// Неизвестная цель — ошибка, а не пустая карта: пустую клиент показал бы
// как «нормы не заданы», не отличив от сбоя.
func Macros(weight float64, goal string) (domain.Macros, error) {
	var result = make(domain.Macros, len(macrosValues))

	for nutrient, byGoal := range macrosValues {
		var bounds, ok = byGoal[goal]
		if !ok {
			return nil, fmt.Errorf("неизвестная цель: %q", goal)
		}
		result[nutrient] = [2]float64{bounds[0] * weight, bounds[1] * weight}
	}

	return result, nil
}
