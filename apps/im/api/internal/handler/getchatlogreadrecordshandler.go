package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"llb-chat/apps/im/api/internal/logic"
	"llb-chat/apps/im/api/internal/svc"
	"llb-chat/apps/im/api/internal/types"
)

// 返回一个处理获取聊天记录已读未读信息的 HTTP 处理函数
func getChatLogReadRecordsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetChatLogReadRecordsReq

		if err := httpx.Parse(r, &req); err != nil {
			// 如果解析请求时发生错误，返回错误响应
			httpx.Error(w, err)
			return
		}

		l := logic.NewGetChatLogReadRecordsLogic(r.Context(), svcCtx)

		// 调用 GetChatLogReadRecords 方法获取聊天记录已读未读信息
		resp, err := l.GetChatLogReadRecords(&req)

		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
