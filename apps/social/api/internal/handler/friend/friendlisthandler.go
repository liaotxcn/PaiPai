package friend

import (
	"net/http"

	"PaiPai/apps/social/api/internal/logic/friend"
	"PaiPai/apps/social/api/internal/svc"
	"PaiPai/apps/social/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 好友列表
func FriendListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FriendListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := friend.NewFriendListLogic(r.Context(), svcCtx)
		resp, err := l.FriendList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
