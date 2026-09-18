package domain

type ActivityType string

// Значения уезжают в БД и приходят из <select> на клиенте, поэтому типизированы:
// нетипизированные константы не мешали передать в Coef любую строку.
const (
	ActivityTypeSedentary ActivityType = "sedentary"
	ActivityTypeLight     ActivityType = "light"
	ActivityTypeModerate  ActivityType = "moderate"
	ActivityTypeHigh      ActivityType = "high"
	ActivityTypeExtreme   ActivityType = "extreme"
)

var activityCoef = map[ActivityType]float64{
	ActivityTypeSedentary: 1.2,
	ActivityTypeLight:     1.375,
	ActivityTypeModerate:  1.55,
	ActivityTypeHigh:      1.725,
	ActivityTypeExtreme:   1.9,
}

// Coef возвращает множитель к базовому обмену. Второе значение — false для
// неизвестной активности: молчаливый ноль обнулял бы всю норму калорий.
func (a ActivityType) Coef() (float64, bool) {
	var coef, ok = activityCoef[a]
	return coef, ok
}

func (a ActivityType) IsValid() bool {
	var _, ok = activityCoef[a]
	return ok
}
