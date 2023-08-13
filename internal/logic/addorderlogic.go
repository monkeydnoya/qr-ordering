package logic

import (
	"context"
	"fmt"

	"qr-ordering-service/internal/svc"
	"qr-ordering-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddOrderLogic {
	return &AddOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddOrderLogic) AddOrder(req *types.AddToOrder) error {
	var notExistingItems []types.ItemCreate
	order, err := l.svcCtx.Db.GetOrder(l.ctx, req.Id)
	fmt.Println(order)

	if err != nil {
		l.Logger.Errorw("order: could not get order",
			logx.LogField{Key: "id", Value: req.Id},
			logx.LogField{Key: "err", Value: err})
		return err
	}

	if order.Status == "done" || order.Status == "deliviring" {
		l.Logger.Errorw("order: order is already deliviring or done",
			logx.LogField{Key: "id", Value: req.Id},
			logx.LogField{Key: "err", Value: err})
		return fmt.Errorf("order is already deliviring or done")
	}

	similar := false
	for _, newItem := range req.Items {
		similar = false
		for _, orderItem := range order.Items {
			if newItem.Name == orderItem.Name {
				err := l.svcCtx.Db.UpdateItemCount(l.ctx, orderItem.Id, newItem.Count)
				if err != nil {
					l.Logger.Errorw("order: failed to update existing item count",
						logx.LogField{Key: "id", Value: req.Id},
						logx.LogField{Key: "err", Value: err})
					return err
				}
				err = l.svcCtx.Db.UpdateOrderSummaryPrice(l.ctx, req.Id, newItem.UnitPrice*float64(newItem.Count))
				if err != nil {
					l.Logger.Errorw("order: failed to update order summary price",
						logx.LogField{Key: "id", Value: req.Id},
						logx.LogField{Key: "err", Value: err})
					return err
				}
				similar = true
				break
			}
		}

		if !similar {
			l.Logger.Infow("order: unique new items:",
				logx.LogField{Key: "id", Value: req.Id},
				logx.LogField{Key: "not existing items", Value: err})
			notExistingItems = append(notExistingItems, newItem)
		}
	}

	if len(notExistingItems) > 0 {
		err = l.svcCtx.Db.AddToOrder(l.ctx, *req)
		if err != nil {
			l.Logger.Errorw("order: could not add new items to order",
				logx.LogField{Key: "id", Value: req.Id},
				logx.LogField{Key: "err", Value: err})
			return err
		}
	}

	return nil
}
