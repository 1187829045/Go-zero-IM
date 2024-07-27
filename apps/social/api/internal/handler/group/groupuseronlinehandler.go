package group

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"llb-chat/apps/social/api/internal/logic/group"
	"llb-chat/apps/social/api/internal/svc"
	"llb-chat/apps/social/api/internal/types"
)

func GroupUserOnlineHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GroupUserOnlineReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := group.NewGroupUserOnlineLogic(r.Context(), svcCtx)
		resp, err := l.GroupUserOnline(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
