package handler

import (
	"net/http"

	"PaiPai/apps/im/api/internal/logic"
	"PaiPai/apps/im/api/internal/svc"
	"PaiPai/apps/im/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取会话
func getConversationsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetConversationsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGetConversationsLogic(r.Context(), svcCtx)
		resp, err := l.GetConversations(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
