package user

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"llb-chat/apps/user/api/internal/logic/user"
	"llb-chat/apps/user/api/internal/svc"
	"llb-chat/apps/user/api/internal/types"
	"net/http"
)

// 将一个普通函数转换为 HTTP 处理器，从而使它能够处理 HTTP 请求

func DetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	// 返回一个符合 http.HandlerFunc 类型的函数，用于处理 HTTP 请求
	return func(w http.ResponseWriter, r *http.Request) {
		// 定义一个 UserInfoReq 类型的变量 req 用于存储请求的数据
		var req types.UserInfoReq
		// 使用 httpx.Parse 解析 HTTP 请求 r，并将解析后的数据存入 req
		if err := httpx.Parse(r, &req); err != nil {
			// 如果解析出错，返回错误信息给客户端
			httpx.Error(w, err)
			return
		}

		// 创建一个 DetailLogic 实例 l，用于处理业务逻辑
		l := user.NewDetailLogic(r.Context(), svcCtx)
		// 调用 l.Detail 方法处理请求，返回响应 resp 和错误信息 err
		resp, err := l.Detail(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
