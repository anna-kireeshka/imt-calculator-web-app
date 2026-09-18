package server

import (
	"net/http"
)

const apiPrefix = "/api/v1/"

type HTTPError struct {
	Code       int
	Message    string
	InnerError error
}

func NewHTTPError(code int, message string, inner error) *HTTPError {
	return &HTTPError{
		Code:       code,
		Message:    message,
		InnerError: inner,
	}
}

type HadlerFuncWithError = func(w http.ResponseWriter, r http.Request)

func (e *HTTPError) Error() string {
	return e.Message
}

func wrapError(endpoint HadlerFuncWithError) {

}

func (h *HandlerMeasurements) HandleRouter(mux *http.ServeMux) {
	mux.HandleFunc("GET"+" "+apiPrefix+"measurements", h.get)
	mux.HandleFunc("POST"+" "+apiPrefix+"measurements", h.save)
	mux.HandleFunc("PATCH /api/v1/measurements", h.update)

	mux.HandleFunc("GET"+" "+apiPrefix+"measurements/bmr", h.getBMR)
}

func (h *HandlerUser) HandleUserRouter(max *http.ServeMux) {
	max.HandleFunc("POST"+" "+apiPrefix+"user", h.save)
}

func (s *Server) apiRouter() http.Handler {
	var mux = http.NewServeMux()

	newMeasurementsHandler(s.repo.Measurement, s.config.TgBotToken).HandleRouter(mux)
	newUserHandler(s.repo.User, s.config.TgBotToken).HandleUserRouter(mux)
	return authMiddleware(s.config.TgBotToken, mux)
}
