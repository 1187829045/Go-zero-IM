package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"llb-chat/apps/im/api/internal/logic"
	"llb-chat/apps/im/api/internal/svc"
	"llb-chat/apps/im/api/internal/types"
)

//这段代码定义了一个 HTTP 处理函数 putConversationsHandler，用于处理 HTTP 请求并执行相应的业务逻辑

func putConversationsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PutConversationsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewPutConversationsLogic(r.Context(), svcCtx)
		resp, err := l.PutConversations(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
