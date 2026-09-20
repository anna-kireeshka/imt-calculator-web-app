# ── Build stage ────────────────────────────────────────────────────────────────
FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux \
    go build -mod=mod -ldflags="-s -w" -trimpath \
    -o /bin/server ./cmd/server

# ── Final stage ───────────────────────────────────────────────────────────────
FROM alpine:latest
WORKDIR /app

COPY --from=builder /bin/server /app/server


RUN apk --no-cache add ca-certificates

EXPOSE 8080

CMD ["/app/server"]

