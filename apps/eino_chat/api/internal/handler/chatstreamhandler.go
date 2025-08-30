package handler

import (
	"net/http"

	"PaiPai/apps/eino_chat/api/internal/logic"
	"PaiPai/apps/eino_chat/api/internal/svc"
	"PaiPai/apps/eino_chat/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 聊天流式输出
func ChatStreamHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChatStreamReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewChatStreamLogic(r.Context(), svcCtx)
		resp, err := l.ChatStream(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
