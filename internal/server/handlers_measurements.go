package server

import (
	"app/imt-calculator-web-app/internal/domain"
	"app/imt-calculator-web-app/internal/fitness"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"time"
)

type HandlerMeasurements struct {
	Measurement domain.MeasurementsRepository
	BotToken    string
}

func newMeasurementsHandler(repo domain.MeasurementsRepository, botToken string) *HandlerMeasurements {
	return &HandlerMeasurements{Measurement: repo, BotToken: botToken}
}

// GET /api/measurement
func (h *HandlerMeasurements) get(w http.ResponseWriter, r *http.Request) {
	var userID, ok = userIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var result, err = h.Measurement.GetByID(r.Context(), userID)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

// POST measurement

func (h *HandlerMeasurements) save(w http.ResponseWriter, r *http.Request) {
	var userID, ok = userIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var body struct {
		Height       int     `json:"height"`
		Weight       float64 `json:"weight"`
		ActivityType string  `json:"activity_type"`
		Goal         string  `json:"goal"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if body.Height <= 0 {
		http.Error(w, "рост должен быть больше нуля", http.StatusBadRequest)
		return
	}

	if body.Weight <= 0 {
		http.Error(w, "вес должен быть больше нуля", http.StatusBadRequest)
		return
	}

	var result, err = h.Measurement.Save(r.Context(), userID, body.Height, body.Weight, body.ActivityType, body.Goal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

func (h *HandlerMeasurements) getBMR(w http.ResponseWriter, r *http.Request) {
	var userID, ok = userIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var resp, err = h.Measurement.GetDataForCalculation(r.Context(), userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			http.Error(w, "нет данных для расчёта", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var age = fitness.AgeAt(resp.DoB, time.Local)

	var ccal, errCcal = fitness.Callories(resp.Weight, resp.Height, age, resp.ActivityType, resp.Goal, resp.Gender)
	if errCcal != nil {
		http.Error(w, errCcal.Error(), http.StatusInternalServerError)
		return
	}

	var macros, errMacros = fitness.Macros(resp.Weight, resp.Goal)
	if errMacros != nil {
		http.Error(w, errMacros.Error(), http.StatusInternalServerError)
		return
	}

	var result = domain.BMR{
		IMT:    fitness.IMT(resp.Height, resp.Weight),
		Ccal:   int(math.Round(ccal)),
		Macros: macros,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(result)
}

func (h *HandlerMeasurements) getHistory(w http.ResponseWriter, r *http.Request) {
	var userID, ok = userIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var resp, err = h.Measurement.GetHistory(r.Context(), userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i := range resp {
		resp[i].IMT = fitness.IMT(resp[i].Height, resp[i].Weight)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(resp)
}
