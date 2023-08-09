package database

import (
	"context"
	"qr-ordering-service/internal/types"
)

type Database interface {
	Start() error
	Stop() error

	GetOrder(ctx context.Context, id string) (types.Order, error)
	CreateOrder(ctx context.Context, order types.OrderCreate) (types.Order, error)

	AddToOrder(ctx context.Context, order types.AddToOrder) error

	UpdateOrderSummaryPrice(ctx context.Context, orderId string, summaryPrice float64) error

	UpdateItemCount(ctx context.Context, itemId string, count int) error
}
