package msgTransfer

import (
	"PaiPai/apps/im/ws/ws"
	"PaiPai/apps/task/mq/internal/svc"
	"PaiPai/apps/task/mq/mq"
	"PaiPai/pkg/bitmap"
	constants "PaiPai/pkg/constant"
	"context"
	"encoding/base64"
	"encoding/json"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
)

var (
	GroupMsgReadRecordDelayTime  = time.Second
	GroupMsgReadRecordDelayCount = 10
)

const (
	GroupMsgReadHandlerAtTransfer = iota
	GroupMsgReadHandlerDelayTransfer
)

// 消费者-处理已读未读
type MsgReadTransfer struct {
	*baseMsgTransfer

	cache.Cache

	mu sync.Mutex

	groupMsgs map[string]*groupMsgRead // 群消息存储(用map因为群数量多)
	push      chan *ws.Push
}

func NewMsgReadTransfer(svc *svc.ServiceContext) *MsgReadTransfer {
	m := &MsgReadTransfer{
		baseMsgTransfer: NewBaseMsgTransfer(svc),
		groupMsgs:       make(map[string]*groupMsgRead, 1),
		push:            make(chan *ws.Push, 1),
	}

	if svc.Config.MsgReadHandler.GroupMsgReadHandler != GroupMsgReadHandlerAtTransfer {
		if svc.Config.MsgReadHandler.GroupMsgReadRecordDelayCount > 0 {
			GroupMsgReadRecordDelayCount = svc.Config.MsgReadHandler.GroupMsgReadRecordDelayCount
		}
		if svc.Config.MsgReadHandler.GroupMsgReadRecordDelayTime > 0 {
			GroupMsgReadRecordDelayTime = time.Duration(svc.Config.MsgReadHandler.GroupMsgReadRecordDelayTime)
		}
	}

	go m.transfer() // 协程调用transfer

	return m
}

func (m *MsgReadTransfer) Consume(ctx context.Context, key, value string) error {
	m.Info("MsgChatTransfer.Consume", value)

	var (
		data mq.MsgMarkRead
	)
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}

	// 业务处理 更新
	readRecords, err := m.UpdateChatLogRead(ctx, &data)
	if err != nil {
		return err
	}

	// map[string]string 已读记录

	push := &ws.Push{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		ContentType:    constants.ContentMakeRead,
		ReadRecords:    readRecords,
	}
	// 判断发送类型
	switch data.ChatType {
	case constants.SingleChatType:
		// 直接推送
		m.push <- push
	case constants.GroupChatType:
		// 判断是否启用合并消息的处理
		m.mu.Lock()
		defer m.mu.Unlock()

		push.SendId = ""
		if _, ok := m.groupMsgs[push.ConversationId]; ok {
			m.Infof("merge push %v", push.ConversationId)
			// 合并请求
			m.groupMsgs[push.ConversationId].mergePush(push)
		} else {
			m.Infof("newGroupMsgRead push %v", push.ConversationId)
			m.groupMsgs[push.ConversationId] = newGroupMsgRead(push, m.push)
		}
	}

	return nil
}

func (m *MsgReadTransfer) transfer() {
	for push := range m.push {
		if push.RecvId != "" || len(push.ReadRecords) > 0 {
			if err := m.Transfer(context.Background(), push); err != nil {
				m.Errorf("m transfer err %v push %v", err, push)
			}
		}
		if push.ChatType == constants.SingleChatType {
			continue
		}
		if m.svcCtx.Config.MsgReadHandler.GroupMsgReadHandler == GroupMsgReadHandlerAtTransfer {
			continue
		}
		// 清空数据 - 重构以避免死锁
		m.mu.Lock()
		msgRead, exists := m.groupMsgs[push.ConversationId]
		m.mu.Unlock()

		if exists {
			// 在锁外检查闲置状态
			if msgRead.IsIdle() {
				m.mu.Lock()
				// 再次检查以避免竞态条件
				if msg, stillExists := m.groupMsgs[push.ConversationId]; stillExists {
					msg.clear()
					delete(m.groupMsgs, push.ConversationId)
				}
				m.mu.Unlock()
			}
		}
	}
}

func (m *MsgReadTransfer) UpdateChatLogRead(ctx context.Context, data *mq.MsgMarkRead) (map[string]string, error) {
	res := make(map[string]string)

	chatLogs, err := m.svcCtx.ChatLogModel.ListByMsgIds(ctx, data.MsgIds)
	if err != nil {
		return nil, err
	}

	// 处理已读
	for _, chatLog := range chatLogs {
		switch chatLog.ChatType {
		case constants.SingleChatType:
			chatLog.ReadRecords = []byte{1}
		case constants.GroupChatType:
			readRecords, err := bitmap.Load(chatLog.ReadRecords)
			if err != nil {
				m.Errorf("加载已读记录失败, chatLogId: %s, err: %v", chatLog.ID.Hex(), err)
				continue
			}
			if err := readRecords.Set(data.SendId); err != nil {
				m.Errorf("设置已读失败, chatLogId: %s, sendId: %s, err: %v", chatLog.ID.Hex(), data.SendId, err)
				continue
			}
			chatLog.ReadRecords = readRecords.Export()
		}

		res[chatLog.ID.Hex()] = base64.StdEncoding.EncodeToString(chatLog.ReadRecords)

		err = m.svcCtx.ChatLogModel.UpdateMakeRead(ctx, chatLog.ID, chatLog.ReadRecords)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
