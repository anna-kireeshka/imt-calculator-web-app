package server

import (
	"log"
	"net/http"
)

const (
	messageInternal   = "Что-то пошло не так, попробуйте позже"
	messageBadRequest = "Некорректные данные запроса"
)

func internalError(w http.ResponseWriter, op string, err error) {
	log.Printf("%s: %v", op, err)
	http.Error(w, messageInternal, http.StatusInternalServerError)
}

func badRequest(w http.ResponseWriter, op string, err error) {
	log.Printf("%s: %v", op, err)
	http.Error(w, messageBadRequest, http.StatusBadRequest)
}
