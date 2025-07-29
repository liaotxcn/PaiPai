package handler

import (
	"net/http"

	"PaiPai/apps/im/api/internal/logic"
	"PaiPai/apps/im/api/internal/svc"
	"PaiPai/apps/im/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 更新会话
func putConversationsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PutConversationsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewPutConversationsLogic(r.Context(), svcCtx)
		resp, err := l.PutConversations(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
