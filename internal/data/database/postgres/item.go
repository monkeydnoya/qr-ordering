package postgres

import (
	"context"
	"fmt"
	"qr-ordering-service/internal/types"

	uuid "github.com/satori/go.uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

func (pg *pgConnection) UpdateItemCount(ctx context.Context, itemId string, count int) error {
	var itemEntity types.ItemEntity
	id, _ := uuid.FromString(itemId)

	tx := pg.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			pg.log.Errorw("postgres: recover transaction: failed to update item count",
				logx.LogField{Key: "item", Value: itemId},
				logx.LogField{Key: "err", Value: r})
		}
	}()
	if err := tx.Model(types.ItemEntity{}).Where("item_id = ?", id).Find(&itemEntity).Error; err != nil {
		tx.Rollback()
		pg.log.Errorw("postgres: failed to find item",
			logx.LogField{Key: "item", Value: itemId},
			logx.LogField{Key: "err", Value: err})
		return err
	}

	updatedCount := itemEntity.Count + count
	updatedSummaryPrice := itemEntity.SummaryPrice + itemEntity.UnitPrice*float64(count)
	fmt.Println(updatedSummaryPrice)

	if err := tx.Model(itemEntity).Update("count", updatedCount).Error; err != nil {
		tx.Rollback()
		pg.log.Errorw("postgres: failed to update item count",
			logx.LogField{Key: "item", Value: itemId},
			logx.LogField{Key: "err", Value: err})
		return err
	}

	if err := tx.Model(itemEntity).Update("summary_price", updatedSummaryPrice).Error; err != nil {
		tx.Rollback()
		pg.log.Errorw("postgres: failed to update item summary price",
			logx.LogField{Key: "item", Value: itemId},
			logx.LogField{Key: "err", Value: err})
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		pg.log.Errorw("postgres: failed to update item",
			logx.LogField{Key: "item", Value: itemId},
			logx.LogField{Key: "err", Value: err})
		return err
	}

	return nil
}
