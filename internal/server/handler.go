package server

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-telegram/bot"
)

type ctxKey int

// userIDKey — ключ, под которым authMiddleware кладёт id пользователя.
// Собственный неэкспортируемый тип: чужой пакет не сможет затереть значение.
const userIDKey ctxKey = iota

// authMiddleware проверяет подпись Telegram initData и кладёт id пользователя
// в контекст. Хендлеры за ним не знают ни про токен, ни про формат заголовка,
// и не могут случайно продолжить работу с неаутентифицированным запросом.
func authMiddleware(botToken string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userID, ok = userIDFromInitData(r, botToken)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userIDKey, userID)))
	})
}

// userIDFromInitData достаёт initData из заголовка "Authorization: tma <initData>"
// и сверяет его подпись токеном бота.
func userIDFromInitData(r *http.Request, botToken string) (int64, bool) {
	var raw, ok = strings.CutPrefix(r.Header.Get("Authorization"), "tma ")
	if !ok {
		return 0, false
	}

	var values, err = url.ParseQuery(raw)
	if err != nil {
		return 0, false
	}

	user, ok := bot.ValidateWebappRequest(values, botToken)
	if !ok {
		return 0, false
	}

	return user.ID, true
}

// userIDFromContext возвращает id, положенный authMiddleware.
func userIDFromContext(ctx context.Context) (int64, bool) {
	var userID, ok = ctx.Value(userIDKey).(int64)
	return userID, ok
}
