package views

import (
	"testing"

	"app/imt-calculator-web-app/internal/fitness"
)

// Значения из <select> уходят в API и попадают прямо в формулы. Если списки
// здесь разъедутся с тем, что понимает fitness, пользователь молча получит
// норму без учёта цели.
func TestGoalsFromOnboardingAreKnownToFitness(t *testing.T) {
	for _, goal := range append(goals, struct {
		Value string
		Label string
	}{defaultGoal, "по умолчанию"}) {
		var macros, err = fitness.Macros(60, goal.Value)
		if err != nil {
			t.Errorf("цель %q (%s) неизвестна fitness.Macros: %v", goal.Value, goal.Label, err)
			continue
		}
		if len(macros) == 0 {
			t.Errorf("цель %q (%s): пустые нормы по нутриентам", goal.Value, goal.Label)
		}
	}
}

func TestActivitiesFromOnboardingAreKnownToDomain(t *testing.T) {
	for _, activity := range activities {
		if !activity.Value.IsValid() {
			t.Errorf("активность %q (%s) нет в таблице коэффициентов", activity.Value, activity.Label)
		}
	}
	if !defaultActivity.IsValid() {
		t.Errorf("активность по умолчанию %q нет в таблице коэффициентов", defaultActivity)
	}
}
