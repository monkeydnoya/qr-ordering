// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"qr-ordering-service/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/orders/create",
				Handler: createOrderHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/api/orders/:id/update",
				Handler: updateOrderHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/api/orders/:id/add",
				Handler: addOrderHandler(serverCtx),
			},
		},
	)
}
