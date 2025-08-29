package handler

import (
	"net/http"

	"PaiPai/apps/eino_chat/api/internal/logic"
	"PaiPai/apps/eino_chat/api/internal/svc"
	"PaiPai/apps/eino_chat/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 聊天补全
func ChatCompletionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChatCompletionReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewChatCompletionLogic(r.Context(), svcCtx)
		resp, err := l.ChatCompletion(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
