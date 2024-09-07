package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"llb-chat/apps/user/api/internal/logic/user"
	"llb-chat/apps/user/api/internal/svc"
	"llb-chat/apps/user/api/internal/types"
)

// LoginHandler 是一个处理登录请求的 HTTP 处理函数。
// 它接受一个 *svc.ServiceContext 参数，用于在处理请求时访问服务的上下文资源（如数据库连接、缓存等）。
func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	// 返回一个匿名的 http.HandlerFunc 函数，这个函数将在接收到 HTTP 请求时被调用。
	return func(w http.ResponseWriter, r *http.Request) {
		// 定义一个 LoginReq 类型的变量 req，用于存储解析后的请求数据。
		var req types.LoginReq

		// 使用 httpx.Parse 函数从 HTTP 请求中解析数据到 req 中。
		// 如果解析过程中发生错误（如请求数据格式不正确），
		// 则调用 httpx.Error 函数返回错误信息，并停止进一步处理。
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		// 创建一个 LoginLogic 实例 l，用于处理具体的登录逻辑。
		// 传入 HTTP 请求的上下文 r.Context() 和服务上下文 svcCtx 以初始化 LoginLogic。
		l := user.NewLoginLogic(r.Context(), svcCtx)

		// 调用 LoginLogic 的 Login 方法，传入解析后的请求数据 req，执行登录操作。
		// 返回登录结果 resp 和可能的错误 err。
		resp, err := l.Login(&req)
		if err != nil {
			// 如果登录过程中发生错误，使用 httpx.Error 函数将错误信息返回给客户端。
			httpx.Error(w, err)
		} else {
			// 如果登录成功，使用 httpx.OkJson 函数将登录结果以 JSON 格式返回给客户端。
			httpx.OkJson(w, resp)
		}
	}
}
