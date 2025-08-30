package handler

import (
	"net/http"

	"PaiPai/apps/eino_chat/api/internal/logic"
	"PaiPai/apps/eino_chat/api/internal/svc"
	"PaiPai/apps/eino_chat/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 调用智能体
func AgentInvokeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AgentInvokeReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewAgentInvokeLogic(r.Context(), svcCtx)
		resp, err := l.AgentInvoke(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
