package server

import (
	"net/http"
)

const apiPrefix = "/api/v1/"

func (h *HandlerMeasurements) HandleRouter(mux *http.ServeMux) {
	mux.HandleFunc("GET"+" "+apiPrefix+"measurements", h.get)
	mux.HandleFunc("POST"+" "+apiPrefix+"measurements", h.save)

	mux.HandleFunc("GET"+" "+apiPrefix+"measurements/bmr", h.getBMR)
	mux.HandleFunc("GET"+" "+apiPrefix+"measurements/history", h.getHistory)
}

func (h *HandlerUser) HandleUserRouter(max *http.ServeMux) {
	max.HandleFunc("GET"+" "+apiPrefix+"user", h.get)
	max.HandleFunc("POST"+" "+apiPrefix+"user", h.save)
}

func (s *Server) apiRouter() http.Handler {
	var mux = http.NewServeMux()

	newMeasurementsHandler(s.repo.Measurement, s.config.TgBotToken).HandleRouter(mux)
	newUserHandler(s.repo.User, s.config.TgBotToken).HandleUserRouter(mux)
	return authMiddleware(s.config.TgBotToken, mux)
}
