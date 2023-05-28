package auth

import "qr-ordering-service/internal/types"

type QrAuth interface {
	ValidateTable(table types.Table) error
}
