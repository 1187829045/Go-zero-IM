// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"user/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/v1/user",
				Handler: userHandler(serverCtx),
			},
		},
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.LoginVerification},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/userinfo",
					Handler: userinfoHandler(serverCtx),
				},
			}...,
		),
	)
}
