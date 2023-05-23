package database

import (
	"context"
	"qr-ordering-service/internal/types"
)

type Database interface {
	Start() error
	Stop() error

	CreateOrder(ctx context.Context, order types.Order) (types.CreatedOrder, error)
}
