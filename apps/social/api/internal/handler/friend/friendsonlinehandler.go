package friend

import (
	"PaiPai/apps/social/api/internal/logic/friend"
	"PaiPai/apps/social/api/internal/svc"
	"PaiPai/apps/social/api/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// FriendsOnlineHandler 好友在线情况
func FriendsOnlineHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FriendsOnlineReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := friend.NewFriendsOnlineLogic(r.Context(), svcCtx)
		resp, err := l.FriendsOnline(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
