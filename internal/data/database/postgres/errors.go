package postgres

import "errors"

var (
	ErrTableIsServicing = errors.New("table is already servicing")
)
