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
)

type MsgChatTransfer struct {
	*baseMsgTransfer
}

func NewMsgChatTransfer(svc *svc.ServiceContext) *MsgChatTransfer {
	return &MsgChatTransfer{
		baseMsgTransfer: NewBaseMsgTransfer(svc),
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

	return m.Transfer(ctx, &ws.Push{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		RecvIds:        data.RecvIds,
		SendTime:       data.SendTime,
		MType:          data.MType,
		Content:        data.Content,
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

	readRecords := bitmap.NewBitmap(0)
	readRecords.Set(chatLog.SendId)
	chatLog.ReadRecords = readRecords.Export()

	err := m.svcCtx.ChatLogModel.Insert(ctx, &model.ChatLog{})
	if err != nil {
		return err
	}

	return m.svcCtx.ConversationModel.UpdateMsg(ctx, &model.ChatLog{})
}
