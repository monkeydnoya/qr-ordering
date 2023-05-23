package logic

import (
	"context"

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

func (l *AddOrderLogic) AddOrder(req *types.Order) (resp *types.CreatedOrder, err error) {
	// todo: add your logic here and delete this line

	return
}
