package database

import (
	"context"
	"qr-ordering-service/internal/types"
)

type Database interface {
	Start() error
	Stop() error

	GetOrder(ctx context.Context, id string) (types.CreatedOrder, error)
	CreateOrder(ctx context.Context, order types.Order) (types.CreatedOrder, error)
	AddToOrder(ctx context.Context, order types.AddToOrder) error
}
