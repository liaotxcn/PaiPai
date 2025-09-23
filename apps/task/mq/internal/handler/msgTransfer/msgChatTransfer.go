package msgTransfer

import (
	model "PaiPai/apps/im/immodels"
	"PaiPai/apps/im/models"
	"PaiPai/apps/im/ws/ws"
	"PaiPai/apps/task/mq/internal/svc"
	"PaiPai/apps/task/mq/mq"
	"PaiPai/pkg/bitmap"
	"context"
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MsgChatTransfer struct {
	*baseMsgTransfer
}

func NewMsgChatTransfer(svc *svc.ServiceContext) *MsgChatTransfer {
	return &MsgChatTransfer{
		baseMsgTransfer: NewBaseMsgTransfer(svc),
	}
}

func (m *MsgChatTransfer) Consume(ctx context.Context, key, value string) error {

	fmt.Println("key : ", key, " value : ", value)

	var (
		data  mq.MsgChatTransfer
		msgId = primitive.NewObjectID()
	)

	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}

	// 记录数据
	if err := m.addChatLog(ctx, msgId, &data); err != nil {
		return err
	}

	return m.Transfer(ctx, &ws.Push{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		RecvIds:        data.RecvIds,
		SendTime:       data.SendTime,
		MType:          data.MType,
		MsgId:          msgId.Hex(),
		Content:        data.Content,
	})
}

func (m *MsgChatTransfer) addChatLog(ctx context.Context, msgId primitive.ObjectID, data *mq.MsgChatTransfer) error {
	// 记录消息
	chatLog := models.ChatLog{
		ID:             msgId,
		ConversationId: data.ConversationId,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		ChatType:       data.ChatType,
		MsgFrom:        0,
		MsgType:        data.MType,
		MsgContent:     data.Content,
		SendTime:       data.SendTime,
	}

	// 设置发送者本人已读
	readRecords, err := bitmap.NewBitmap(0)
	if err != nil {
		return err
	}
	if err := readRecords.Set(chatLog.SendId); err != nil {
		return err
	}
	chatLog.ReadRecords = readRecords.Export()

	// 转换为model.ChatLog类型
	modelChatLog := &model.ChatLog{
		ID:             chatLog.ID,
		ConversationId: chatLog.ConversationId,
		SendId:         chatLog.SendId,
		RecvId:         chatLog.RecvId,
		ChatType:       chatLog.ChatType,
		MsgFrom:        chatLog.MsgFrom,
		MsgType:        chatLog.MsgType,
		MsgContent:     chatLog.MsgContent,
		SendTime:       chatLog.SendTime,
		ReadRecords:    chatLog.ReadRecords,
	}

	err = m.svcCtx.ChatLogModel.Insert(ctx, modelChatLog)
	if err != nil {
		return err
	}

	// 更新会话最新消息
	return m.svcCtx.ConversationModel.UpdateMsg(ctx, modelChatLog)
}
