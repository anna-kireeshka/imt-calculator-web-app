package database

import (
	"app/imt-calculator-web-app/internal/domain"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	conn        *pgx.Conn
	User        domain.UserRepository
	Measurement domain.MeasurementsRepository
}

func NewRepository(ctx context.Context, dsn string) (r *Repository, err error) {
	var config *pgx.ConnConfig
	if config, err = pgx.ParseConfig(dsn); err != nil {
		return nil, fmt.Errorf("parse dsn: %w", err)
	}

	var conn *pgx.Conn

	if conn, err = pgx.ConnectConfig(ctx, config); err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err = conn.Ping(pingCtx); err != nil {
		_ = conn.Close(ctx)
		return nil, fmt.Errorf("ping: %w", err)
	}

	r = &Repository{conn: conn}
	r.User = NewUserRepository(conn)
	r.Measurement = NewMeasurementsRepository(conn)

	return r, nil
}

func (r *Repository) Close(ctx context.Context) error {
	return r.conn.Close(ctx)
}
