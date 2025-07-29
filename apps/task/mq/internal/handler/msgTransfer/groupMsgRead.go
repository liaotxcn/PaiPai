package msgTransfer

import (
	"PaiPai/apps/im/ws/ws"
	constants "PaiPai/pkg/constant"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
	"time"
)

type groupMsgRead struct {
	mu             sync.Mutex
	conversationId string
	push           *ws.Push
	pushCh         chan *ws.Push
	count          int
	pushTime       time.Time // 上一次推送时间
	done           chan struct{}
}

func newGroupMsgRead(push *ws.Push, pushCh chan *ws.Push) *groupMsgRead {
	m := &groupMsgRead{
		conversationId: push.ConversationId,
		push:           push,
		pushCh:         pushCh,
		count:          0,
		pushTime:       time.Now(),
		done:           make(chan struct{}),
	}
	go m.transfer()
	return m
}

// 合并消息
func (m *groupMsgRead) mergePush(push *ws.Push) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.count++
	for msgId, read := range push.ReadRecords {
		m.push.ReadRecords[msgId] = read
	}
}

func (m *groupMsgRead) transfer() {
	// 超时发送
	// 超量发送

	timer := time.NewTimer(GroupMsgReadRecordDelayTime / 2)
	defer timer.Stop()
	for {
		select {
		case <-m.done:
			return
		case <-timer.C:
			m.mu.Lock()
			pushTime := m.pushTime
			val := GroupMsgReadRecordDelayTime - time.Since(pushTime)
			push := m.push

			if val > 0 && m.count < GroupMsgReadRecordDelayCount || push == nil {
				//未达标
				m.mu.Unlock()
				continue
			}

			m.pushTime = time.Now()
			m.push = nil
			m.count = 0
			m.mu.Unlock()

			// 推送
			logx.Infof("超过合并的条件推送 %v", push)
			m.pushCh <- push

		// 其他情况
		default:
			m.mu.Lock()
			if m.count >= GroupMsgReadRecordDelayCount {
				push := m.push
				m.push = nil
				m.count = 0
				m.mu.Unlock()
				logx.Infof("ddefault 超过合并的条件推送 %v", push)
				m.pushCh <- push
				continue
			}
			// 为空时
			if m.isIdle() {
				m.mu.Unlock()
				// 使msgReadTransfer释放
				m.pushCh <- &ws.Push{
					ChatType:       constants.GroupChatType,
					ConversationId: m.conversationId,
				}
				continue
			}
			// 不为空时
			m.mu.Unlock()
			tempDelay := GroupMsgReadRecordDelayTime / 4
			if tempDelay > time.Second {
				tempDelay = time.Second
			}
			time.Sleep(tempDelay)
		}
	}
}

func (m *groupMsgRead) IsIdle() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.IsIdle()
}

func (m *groupMsgRead) isIdle() bool {
	pushTime := m.pushTime
	val := GroupMsgReadRecordDelayTime*2 - time.Since(pushTime)

	if val <= 0 && m.push == nil && m.count == 0 {
		return true
	}
	return false
}

func (m *groupMsgRead) clear() {
	select {
	case <-m.done:
	default:
		close(m.done)
	}
	m.push = nil
}
