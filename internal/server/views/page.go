package views

import (
	"app/imt-calculator-web-app/internal/domain"

	"github.com/a-h/templ"
)

// Page — данные, общие для любого экрана. Tab подсвечивает пункт таб-бара;
// пустой Tab убирает таб-бар совсем (экраны ошибок и «вне Telegram»).
type Page struct {
	Title    string
	Tab      string
	Name     string
	Initials string
	Detail   string
	Centered bool
	Scripts  []string
}

// tabs перечисляет пункты нижней навигации в порядке показа.
var tabs = []struct {
	ID    string
	Href  string
	Icon  string
	Label string
}{
	{ID: "bmi", Href: "/", Icon: "ti-scale", Label: "ИМТ"},
	{ID: "history", Href: "/history", Icon: "ti-history", Label: "История"},
	// Дневник и настройки временно скрыты — экраны ещё не на живых данных.
	// Сами страницы остаются доступны по /diary и /settings.
	// {ID: "calories", Href: "/diary", Icon: "ti-flame", Label: "Дневник"},
	// {ID: "settings", Href: "/settings", Icon: "ti-settings", Label: "Ещё"},
}

// current помечает активный пункт: атрибут либо присутствует, либо его нет
// вовсе — aria-current="false" скринридеры трактуют иначе, чем отсутствие.
func current(active bool) templ.Attributes {
	if active {
		return templ.Attributes{"aria-current": "page"}
	}
	return nil
}

// Списки для онбординга живут здесь, а не в разметке: значения уходят в API,
// и разъехаться с domain.ActivityType они не должны.
var activities = []struct {
	Value domain.ActivityType
	Label string
}{
	{domain.ActivityTypeSedentary, "Сидячий, без тренировок"},
	{domain.ActivityTypeLight, "Лёгкая, 1–3 тренировки в неделю"},
	{domain.ActivityTypeModerate, "Умеренная, 3–5 в неделю"},
	{domain.ActivityTypeHigh, "Высокая, 6–7 в неделю"},
	{domain.ActivityTypeExtreme, "Очень высокая, тяжёлая работа"},
}

const defaultActivity = domain.ActivityTypeModerate

var goals = []struct {
	Value string
	Label string
}{
	{"lose", "Похудение"},
	{"gain", "Набор"},
	{"maintain", "Удержание"},
}

const defaultGoal = "maintain"

// boolAttr — aria-pressed принимает строку, а не булево.
func boolAttr(value bool) string {
	if value {
		return "true"
	}
	return "false"
}

// months — значения совпадают с номером месяца в ISO-дате, которую
// собирает onboarding.js: "01".."12".
var months = []struct {
	Value string
	Label string
}{
	{"01", "Январь"}, {"02", "Февраль"}, {"03", "Март"},
	{"04", "Апрель"}, {"05", "Май"}, {"06", "Июнь"},
	{"07", "Июль"}, {"08", "Август"}, {"09", "Сентябрь"},
	{"10", "Октябрь"}, {"11", "Ноябрь"}, {"12", "Декабрь"},
}
