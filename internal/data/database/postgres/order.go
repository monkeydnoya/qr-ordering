package postgres

import (
	"context"
	"fmt"
	"qr-ordering-service/internal/types"
	"time"

	uuid "github.com/satori/go.uuid"
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
		CreatedDate:  time.Now().Unix(),
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

func (pg *pgConnection) AddToOrder(ctx context.Context, addition types.AddToOrder) error {
	var newSummaryPrice float64
	order := types.OrderEntity{
		Id: uuid.FromStringOrNil(addition.Id),
	}
	newItemsEntity := make([]types.ItemEntity, 0)
	for _, v := range addition.Items {
		newItemsEntity = append(newItemsEntity, v.ToEntity())
		newSummaryPrice = newSummaryPrice + v.ToEntity().SummaryPrice
	}
	tx := pg.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			pg.log.Errorw("postgres: recover transaction: failed to add new items to order",
				logx.LogField{Key: "id", Value: addition.Id},
				logx.LogField{Key: "err", Value: r})
		}
	}()
	for _, item := range newItemsEntity {
		if err := tx.Model(&order).Association("Items").Append(&item); err != nil {
			tx.Rollback()
			pg.log.Errorw("postgres: can't find order with that id",
				logx.LogField{Key: "id", Value: addition.Id},
				logx.LogField{Key: "err", Value: err})
			return err
		}
	}
	if err := tx.Select("summary_price").Find(&order).Error; err != nil {
		tx.Rollback()
		pg.log.Errorw("postgres: can't find order with that id",
			logx.LogField{Key: "id", Value: addition.Id},
			logx.LogField{Key: "err", Value: err})
		return err
	}
	order.SummaryPrice = order.SummaryPrice + newSummaryPrice
	if err := tx.Model(&order).Update("summary_price", order.SummaryPrice).Error; err != nil {
		tx.Rollback()
		pg.log.Errorw("postgres: failed to update summary price",
			logx.LogField{Key: "id", Value: addition.Id},
			logx.LogField{Key: "err", Value: err.Error()})
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		pg.log.Errorw("postgres: failed to add new items to order - commit error",
			logx.LogField{Key: "id", Value: addition.Id},
			logx.LogField{Key: "err", Value: err.Error()})
		return err
	}
	return nil
}
