package postgres

import (
	"context"
	"qr-ordering-service/internal/types"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

func (pg *pgConnection) GetOrder(ctx context.Context, id string) (types.Order, error) {
	var orderEntity types.OrderEntity
	if err := pg.db.Model(types.OrderEntity{}).Where("id = ?", id).Find(&orderEntity).Error; err != nil {
		pg.log.Errorw("postgres: failed to find order",
			logx.LogField{Key: "order", Value: id},
			logx.LogField{Key: "err", Value: err})
		return types.Order{}, err
	}

	var items []types.ItemEntity
	if err := pg.db.Model(orderEntity).Association("Items").Find(&items); err != nil {
		pg.log.Errorw("postgres: failed to find order",
			logx.LogField{Key: "order", Value: id},
			logx.LogField{Key: "err", Value: err})
		return types.Order{}, err
	}

	orderEntity.Items = items
	return orderEntity.ToModel(), nil
}

func (pg *pgConnection) CreateOrder(ctx context.Context, order types.OrderCreate) (types.Order, error) {
	var summaryPrice float64
	var existingOrders []types.OrderEntity
	itemsEntity := make([]types.ItemEntity, 0)
	for i, item := range order.Items {
		itemEntity := types.ItemEntity{
			Name:         item.Name,
			Count:        item.Count,
			PiecePrice:   item.PiecePrice,
			SummaryPrice: float64(item.Count) * item.PiecePrice,
			IsDilivered:  false,
		}
		itemsEntity = append(itemsEntity, itemEntity)
		summaryPrice = summaryPrice + float64(itemsEntity[i].SummaryPrice)
	}

	orderEntity := types.OrderEntity{
		CreatedDate:  time.Now().Unix(),
		Table:        order.Table,
		Paid:         false,
		Status:       "pending",
		SummaryPrice: summaryPrice,
		Items:        itemsEntity,
	}

	tx := pg.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			pg.log.Errorw("postgres: recover transaction: failed to create order",
				logx.LogField{Key: "table", Value: order.Table},
				logx.LogField{Key: "err", Value: r})
		}
	}()
	// Search order with this table with status !done.
	if err := tx.Model(types.OrderEntity{Table: order.Table}).Where("public.orders.table = ? AND status != 'done'", order.Table).Find(&existingOrders).Error; err != nil {
		tx.Rollback()
		pg.log.Errorw("postgres: failed to search order of existing table",
			logx.LogField{Key: "table", Value: order.Table},
			logx.LogField{Key: "err", Value: err})
		return types.Order{}, err
	}

	if len(existingOrders) > 0 {
		tx.Rollback()
		pg.log.Errorw("postgres: failed to create order",
			logx.LogField{Key: "table", Value: order.Table},
			logx.LogField{Key: "err", Value: ErrTableIsServicing})
		return types.Order{}, ErrTableIsServicing
	}

	if err := tx.Create(&orderEntity).Error; err != nil {
		tx.Rollback()
		pg.log.Errorw("postgres: failed to create order",
			logx.LogField{Key: "table", Value: order.Table},
			logx.LogField{Key: "err", Value: err})
		return types.Order{}, err
	}
	if err := tx.First(&orderEntity).Error; err != nil {
		tx.Rollback()
		pg.log.Errorw("postgres: failed to select order",
			logx.LogField{Key: "table", Value: order.Table},
			logx.LogField{Key: "err", Value: err})
		return types.Order{}, err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		pg.log.Errorw("postgres: failed to create order",
			logx.LogField{Key: "table", Value: order.Table},
			logx.LogField{Key: "err", Value: err})
		return types.Order{}, err
	}

	return orderEntity.ToModel(), nil
}

func (pg *pgConnection) AddToOrder(ctx context.Context, addition types.AddToOrder) error {
	var newSummaryPrice float64
	order := types.OrderEntity{
		Id: uuid.FromStringOrNil(addition.Id),
	}
	newItemsEntities := make([]types.ItemEntity, 0)
	for _, item := range addition.Items {
		newItemsEntity := types.ItemEntity{
			Name:         item.Name,
			Count:        item.Count,
			PiecePrice:   item.PiecePrice,
			SummaryPrice: float64(item.Count) * item.PiecePrice,
			IsDilivered:  false,
		}
		newItemsEntities = append(newItemsEntities, newItemsEntity)
		newSummaryPrice = newSummaryPrice + newItemsEntity.SummaryPrice
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
	for _, item := range newItemsEntities {
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

func (pg *pgConnection) UpdateOrderSummaryPrice(ctx context.Context, orderId string, summaryPrice float64) error {
	var order types.OrderEntity
	id, _ := uuid.FromString(orderId)

	tx := pg.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			pg.log.Errorw("postgres: recover transaction: failed to add new items to order",
				logx.LogField{Key: "id", Value: orderId},
				logx.LogField{Key: "err", Value: r})
		}
	}()
	if err := tx.Model(types.OrderEntity{}).Where("id = ?", id).Find(&order).Error; err != nil {
		tx.Rollback()
		pg.log.Errorw("postgres: can't find order with that id",
			logx.LogField{Key: "id", Value: orderId},
			logx.LogField{Key: "err", Value: err})
		return err
	}

	newSummaryPrice := order.SummaryPrice + summaryPrice
	if err := tx.Model(&order).Update("summary_price", newSummaryPrice).Error; err != nil {
		tx.Rollback()
		pg.log.Errorw("postgres: failed to update summary price",
			logx.LogField{Key: "id", Value: orderId},
			logx.LogField{Key: "err", Value: err.Error()})
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		pg.log.Errorw("postgres: failed to add new items to order - commit error",
			logx.LogField{Key: "id", Value: orderId},
			logx.LogField{Key: "err", Value: err.Error()})
		return err
	}
	return nil
}
