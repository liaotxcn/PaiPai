package push

import (
	"PaiPai/apps/im/ws/internal/svc"
	"PaiPai/apps/im/ws/websocket"
	"PaiPai/apps/im/ws/ws"
	constants "PaiPai/pkg/constant"
	"github.com/mitchellh/mapstructure"
)

func Push(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		var data ws.Push
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err))
			return
		}

		// 发送的目标
		switch data.ChatType {
		case constants.SingleChatType:
			err := single(srv, &data, data.RecvId)
			if err != nil {
				srv.Error(err)
			}
		case constants.GroupChatType:
			err := group(srv, &data)
			if err != nil {
				srv.Error(err)
			}
		default:
		}
	}
}

// 私聊推送处理
func single(srv *websocket.Server, data *ws.Push, recvId string) error {
	rconn := srv.GetConn(recvId)
	if rconn == nil {
		// todo: 目标离线
		return nil
	}

	srv.Infof("push msg %v", data)

	return srv.Send(websocket.NewMessage(data.SendId, &ws.Chat{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendTime:       data.SendTime,
		Msg: ws.Msg{
			MsgId:       data.MsgId,
			MType:       data.MType,
			Content:     data.Content,
			ReadRecords: data.ReadRecords,
		},
	}), rconn[0])

}

// 群聊推送处理
func group(srv *websocket.Server, data *ws.Push) error {
	for _, id := range data.RecvIds {
		func(id string) {
			srv.Schedule(func() {
				single(srv, data, id)
			})
		}(id)
	}
	return nil
}
