package fitness

import (
	"fmt"

	"app/imt-calculator-web-app/internal/domain"
)

// Поправка на пол в формуле Миффлина — Сан Жеора.
// Gender: true — женщина, false — мужчина (см. domain.User).
const (
	genderShiftFemale = -161
	genderShiftMale   = 5
)

// Поправка на цель: дефицит и профицит в килокалориях.
const goalShift = 500

// Callories считает суточную норму по Миффлину — Сан Жеору с учётом активности
// и цели. Таблица коэффициентов одна на приложение — в domain, иначе локальная
// копия разъезжается со значениями, которые реально приходят из <select> и БД.
func Callories(weight float64, height int, age int, activityType string, goal string, gender bool) (float64, error) {
	var coef, ok = domain.ActivityType(activityType).Coef()
	if !ok {
		return 0, fmt.Errorf("неизвестный тип активности: %q", activityType)
	}

	var genderShift float64 = genderShiftMale
	if gender {
		genderShift = genderShiftFemale
	}

	var bmr = (10 * weight) + (6.25 * float64(height)) - (5 * float64(age)) + genderShift
	var ccal = bmr * coef

	switch goal {
	case "lose":
		ccal -= goalShift
	case "gain":
		ccal += goalShift
	}

	return ccal, nil
}
