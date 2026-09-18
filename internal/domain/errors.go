package domain

import "errors"

// ErrNotFound — общий признак «данных нет». Хендлеры отличают по нему
// пустую выборку (404) от поломки хранилища (500), не зная про pgx.
var ErrNotFound = errors.New("не найдено")
