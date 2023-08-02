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
	order, err := l.svcCtx.Db.GetOrder(l.ctx, req.Id)
	if err != nil {
		l.Logger.Errorw("order: could not get order",
			logx.LogField{Key: "id", Value: req.Id},
			logx.LogField{Key: "err", Value: err})
		return err
	}

	fmt.Println(order.Status)
	if order.Status == "done" || order.Status == "deliviring" {
		l.Logger.Errorw("order: order is already deliviring or done",
			logx.LogField{Key: "id", Value: req.Id},
			logx.LogField{Key: "err", Value: err})
		return fmt.Errorf("order is already deliviring or done")
	} else {
		err = l.svcCtx.Db.AddToOrder(l.ctx, *req)
		if err != nil {
			l.Logger.Errorw("order: could not add new items to order",
				logx.LogField{Key: "id", Value: req.Id},
				logx.LogField{Key: "err", Value: err})
			return err
		}
		return nil
	}

}
