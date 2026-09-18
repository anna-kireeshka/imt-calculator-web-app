# ── Build stage ────────────────────────────────────────────────────────────────
FROM golang:1.27-alpine AS builder

WORKDIR /app

# Копируем только go.mod; go.sum генерируется заново через go mod tidy,
# чтобы избежать конфликтов контрольных сумм между окружениями.
COPY go.mod go.sum ./
RUN go mod download
# Копируем исходники и собираем.
# -mod=mod разрешает обновить go.sum если он неполный (безопасно для dev).
COPY . .
RUN CGO_ENABLED=0 GOOS=linux \
    go build -mod=mod -ldflags="-s -w" -trimpath \
    -o /bin/server ./cmd/server

# ── Final stage ───────────────────────────────────────────────────────────────
FROM alpine:latest
WORKDIR /app

# Копируем собранный бинарник из builder stage.
COPY --from=builder /bin/server /app/server
#migrations
COPY --from=builder /app/database/migrations /app/database/migrations


# Устанавливаем необходимые зависимости для работы приложения.
RUN apk --no-cache add ca-certificates

# Указываем порт, который будет слушать приложение.
EXPOSE 8080
# Запускаем приложение при старте контейнера.
CMD ["/app/server"]

