package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"llb-chat/apps/im/api/internal/logic"
	"llb-chat/apps/im/api/internal/svc"
	"llb-chat/apps/im/api/internal/types"
)

// getChatLogHandler 返回一个处理获取聊天记录的 HTTP 处理函数
func getChatLogHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	// 返回一个 HTTP 处理函数
	return func(w http.ResponseWriter, r *http.Request) {
		// 定义一个 ChatLogReq 结构体变量来存储请求数据
		var req types.ChatLogReq
		// 解析 HTTP 请求，将请求数据绑定到 req 结构体
		if err := httpx.Parse(r, &req); err != nil {
			// 如果解析请求时发生错误，返回错误响应
			httpx.Error(w, err)
			return
		}
		// 创建一个新的 GetChatLogLogic 实例
		l := logic.NewGetChatLogLogic(r.Context(), svcCtx)
		// 调用 GetChatLog 方法获取聊天记录
		resp, err := l.GetChatLog(&req)
		// 如果获取聊天记录时发生错误，返回错误响应
		if err != nil {
			httpx.Error(w, err)
		} else {
			// 否则，返回成功的 JSON 响应
			httpx.OkJson(w, resp)
		}
	}
}
