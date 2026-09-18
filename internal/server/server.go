package server

import (
	"app/imt-calculator-web-app/internal/config"
	"app/imt-calculator-web-app/internal/database"
	"app/imt-calculator-web-app/internal/server/views"
	"io/fs"
	"net/http"
)

type Server struct {
	config *config.Config
	repo   *database.Repository
}

func New(cfg *config.Config, repo *database.Repository) *Server {
	return &Server{
		config: cfg,
		repo:   repo,
	}
}

func (s *Server) Router() http.Handler {
	var mux = http.NewServeMux()

	mux.HandleFunc("GET /{$}", s.handleResult)
	mux.HandleFunc("GET /onboarding", s.handleOnboarding)
	mux.HandleFunc("GET /diary", s.handleDiary)

	// Расчёт калорий переехал на главную. Старый адрес мог остаться
	// в закладках и в истории Telegram, поэтому не 404, а редирект.
	mux.HandleFunc("GET /calories", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	})
	mux.HandleFunc("GET /history", s.handleHistory)
	mux.HandleFunc("GET /goal", s.handleGoal)
	mux.HandleFunc("GET /settings", s.handleSettings)

	mux.HandleFunc("GET /loading", s.handleLoading)
	mux.HandleFunc("GET /offline", s.handleOffline)
	mux.HandleFunc("GET /failure", s.handleFailure)
	mux.HandleFunc("GET /outside", s.handleOutside)

	static, err := fs.Sub(staticFS, "static")
	if err != nil {
		panic(err)
	}
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(static))))

	mux.Handle("/api/", s.apiRouter())

	return mux
}

func (s *Server) handleResult(w http.ResponseWriter, r *http.Request) {
	s.renderComponent(r.Context(), w, views.Result(views.Page{
		Title:   "Ваши показатели",
		Tab:     "bmi",
		Scripts: []string{"/static/js/result.js"},
	}))
}

// Онбординг идёт без таб-бара: пустой Tab убирает его в Layout.
func (s *Server) handleOnboarding(w http.ResponseWriter, r *http.Request) {
	s.renderComponent(r.Context(), w, views.Onboarding(views.Page{
		Title:   "Знакомимся",
		Scripts: []string{"/static/js/onboarding.js"},
	}))
}

func (s *Server) handleDiary(w http.ResponseWriter, r *http.Request) {
	s.render(w, "diary.html", PageData{Title: "Дневник калорий", Tab: "calories"})
}

func (s *Server) handleHistory(w http.ResponseWriter, r *http.Request) {
	s.render(w, "history.html", PageData{Title: "История", Tab: "history"})
}

func (s *Server) handleGoal(w http.ResponseWriter, r *http.Request) {
	s.render(w, "goal.html", PageData{Title: "Цель по весу", Tab: "settings"})
}

func (s *Server) handleSettings(w http.ResponseWriter, r *http.Request) {
	s.render(w, "settings.html", PageData{Title: "Ещё", Tab: "settings", Name: "Alex Kireev"})
}

func (s *Server) handleLoading(w http.ResponseWriter, r *http.Request) {
	s.render(w, "loading.html", PageData{Title: "Загрузка", Tab: "bmi"})
}

func (s *Server) handleOffline(w http.ResponseWriter, r *http.Request) {
	s.render(w, "offline.html", PageData{Title: "Нет соединения", Centered: true})
}

func (s *Server) handleFailure(w http.ResponseWriter, r *http.Request) {
	s.renderStatus(w, "failure.html", PageData{
		Title:    "Сервис недоступен",
		Centered: true,
		Detail:   "503, запрос 9f2a1c, 14:32:07",
	}, http.StatusServiceUnavailable)
}

func (s *Server) handleOutside(w http.ResponseWriter, r *http.Request) {
	s.render(w, "outside.html", PageData{Title: "Откройте через Telegram", Centered: true})
}
