package fitness

import (
	"fmt"

	"app/imt-calculator-web-app/internal/domain"
)

const (
	genderShiftFemale = -161
	genderShiftMale   = 5
)

const goalShift = 500

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
