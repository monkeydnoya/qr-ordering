package logic

import (
	"context"

	"qr-ordering-service/internal/svc"
	"qr-ordering-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateOrderLogic) CreateOrder(req *types.OrderCreate) (resp *types.Order, err error) {
	createdOrder, err := l.svcCtx.Db.CreateOrder(l.ctx, *req)
	if err != nil {
		l.Logger.Errorw("order: could not create order",
			logx.LogField{Key: "table", Value: req.Table},
			logx.LogField{Key: "err", Value: err})
		return &createdOrder, err
	}
	return &createdOrder, nil
}
