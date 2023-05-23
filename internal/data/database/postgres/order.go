package postgres

import (
	"context"
	"fmt"
	"qr-ordering-service/internal/types"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

func (pg *pgConnection) CreateOrder(ctx context.Context, order types.Order) (types.CreatedOrder, error) {
	var summaryPrice float64
	fmt.Println(order.Items)
	itemsEntity := make([]types.ItemEntity, 0)
	for i, item := range order.Items {
		itemsEntity = append(itemsEntity, item.ToEntity())
		summaryPrice = summaryPrice + float64(itemsEntity[i].SummaryPrice)
	}
	fmt.Println(itemsEntity)

	orderEntity := types.OrderEntity{
		CreatedDate:  time.Now(),
		Table:        order.Table,
		Paid:         false,
		Status:       "pending",
		SummaryPrice: summaryPrice,
		Items:        itemsEntity,
	}
	tx := pg.db.WithContext(ctx).Begin()
	// TODO: Read about Transaction Recover
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			pg.log.Errorw("postgres: recover transaction: failed to create order",
				logx.LogField{Key: "table", Value: order.Table},
				logx.LogField{Key: "err", Value: r})
		}
	}()
	fmt.Println(orderEntity)
	if err := tx.Create(&orderEntity).Error; err != nil {
		tx.Rollback()
		pg.log.Errorw("postgres: failed to create order",
			logx.LogField{Key: "table", Value: order.Table},
			logx.LogField{Key: "err", Value: err})
		return types.CreatedOrder{}, err
	}
	if err := tx.First(&orderEntity).Error; err != nil {
		tx.Rollback()
		pg.log.Errorw("postgres: failed to select order",
			logx.LogField{Key: "table", Value: order.Table},
			logx.LogField{Key: "err", Value: err})
		return types.CreatedOrder{}, err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		pg.log.Errorw("postgres: failed to create order",
			logx.LogField{Key: "table", Value: order.Table},
			logx.LogField{Key: "err", Value: err})
		return types.CreatedOrder{}, err
	}

	return orderEntity.ToModel(), nil
}
