package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"qr-ordering-service/internal/logic"
	"qr-ordering-service/internal/svc"
	"qr-ordering-service/internal/types"
)

func updateOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Order
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUpdateOrderLogic(r.Context(), svcCtx)
		resp, err := l.UpdateOrder(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
