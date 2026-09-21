package server

import (
	"app/imt-calculator-web-app/internal/domain"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type HandlerUser struct {
	User     domain.UserRepository
	BotToken string
}

func newUserHandler(repo domain.UserRepository, botToken string) *HandlerUser {
	return &HandlerUser{User: repo, BotToken: botToken}
}

func (h *HandlerUser) save(w http.ResponseWriter, r *http.Request) {
	var userID, ok = userIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var body struct {
		DoB    time.Time `json:"dob"`
		Gender bool      `json:"gender"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		badRequest(w, "user save: decode body", err)
		return
	}

	var result, err = h.User.Save(r.Context(), userID, body.DoB, body.Gender)

	if err != nil {
		internalError(w, "user save", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)

}

func (h *HandlerUser) get(w http.ResponseWriter, r *http.Request) {
	var userID, ok = userIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var result, err = h.User.GetByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			http.Error(w, "пользователь не заведён", http.StatusNotFound)
			return
		}
		internalError(w, "user get", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(result)
}
