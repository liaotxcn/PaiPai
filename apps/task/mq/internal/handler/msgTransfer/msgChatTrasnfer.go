package msgTransfer

import (
	model "PaiPai/apps/im/immodels"
	"PaiPai/apps/im/models"
	"PaiPai/apps/im/ws/websocket"
	"PaiPai/apps/social/rpc/socialclient"
	"PaiPai/apps/task/mq/internal/svc"
	"PaiPai/apps/task/mq/mq"
	constants "PaiPai/pkg/constant"
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
)

type MsgChatTransfer struct {
	logx.Logger
	svc *svc.ServiceContext
}

func NewMsgChatTransfer(svc *svc.ServiceContext) *MsgChatTransfer {
	return &MsgChatTransfer{
		Logger: logx.WithContext(context.Background()),
		svc:    svc,
	}
}

func (m *MsgChatTransfer) Consume(key, value string) error {
	fmt.Println("key : ", key, " value : ", value)

	var (
		data mq.MsgChatTransfer
		ctx  = context.Background()
	)
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}

	// 记录数据
	if err := m.addChatLog(ctx, &data); err != nil {
		switch data.ChatType {
		case constants.SingleChatType:
			return m.single(&data)
		case constants.GroupChatType:
			return m.group(ctx, &data)
		}
	}
	return nil
}

func (m *MsgChatTransfer) single(data *mq.MsgChatTransfer) error {
	// 私聊推送消息
	return m.svc.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FormId:    constants.SYSTEM_ROOT_UID,
		Data:      data,
	})
}

func (m *MsgChatTransfer) group(ctx context.Context, data *mq.MsgChatTransfer) error {
	// 群聊推送消息
	users, err := m.svc.Social.GroupUsers(ctx, &socialclient.GroupUsersReq{
		GroupId: data.RecvId,
	})
	if err != nil {
		return err
	}
	data.RecvIds = make([]string, 0, len(users.List))
	for _, menbers := range users.List {
		if menbers.UserId == data.SendId {
			continue
		}
		data.RecvIds = append(data.RecvIds, menbers.UserId)
	}
	return m.svc.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FormId:    constants.SYSTEM_ROOT_UID,
		Data:      data,
	})
}

func (m *MsgChatTransfer) addChatLog(ctx context.Context, data *mq.MsgChatTransfer) error {
	// 记录消息
	chatLog := models.ChatLog{
		ConversationId: data.ConversationId,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		ChatType:       data.ChatType,
		MsgFrom:        0,
		MsgType:        data.MType,
		MsgContent:     data.Content,
		SendTime:       data.SendTime,
	}
	err := m.svc.ChatLogModel.Insert(ctx, (*model.ChatLog)(&chatLog))
	if err != nil {
		return err
	}

	return m.svc.ConversationModel.UpdateMsg(ctx, (*model.ChatLog)(&chatLog))
}
