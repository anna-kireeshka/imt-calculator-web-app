package database

import (
	"errors"
	"fmt"

	"app/imt-calculator-web-app/internal/domain"

	"github.com/jackc/pgx/v5"
)

func wrapQueryError(err error, op string, notFoundMessage string) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("%w: %s", domain.ErrNotFound, notFoundMessage)
	}
	return fmt.Errorf("%s: %w", op, err)
}
