package msgTransfer

import (
	"github.com/zeromicro/go-zero/core/logx"
	"llb-chat/apps/im/ws/ws"
	"llb-chat/pkg/constants"
	"sync"
	"time"
)

type groupMsgRead struct {
	mu             sync.Mutex
	conversationId string
	push           *ws.Push
	pushCh         chan *ws.Push
	count          int
	// 上次推送时间
	pushTime time.Time
	done     chan struct{}
}

func newGroupMsgRead(push *ws.Push, pushCh chan *ws.Push) *groupMsgRead {
	m := &groupMsgRead{
		conversationId: push.ConversationId,
		push:           push,
		pushCh:         pushCh,
		count:          1,
		pushTime:       time.Now(),
		done:           make(chan struct{}),
	}

	go m.transfer()
	return m
}

// 合并消息
// 将新的推送记录合并到当前的推送记录中。
func (m *groupMsgRead) mergePush(push *ws.Push) {
	// 加锁以保护共享数据
	m.mu.Lock()
	defer m.mu.Unlock()

	// 更新记录计数
	m.count++

	// 将新记录合并到现有记录中
	for msgId, read := range push.ReadRecords {
		m.push.ReadRecords[msgId] = read
	}
}

// transfer 方法处理推送记录的合并和发送。
func (m *groupMsgRead) transfer() {
	// 1. 超时发送
	// 2. 超量发送

	// 创建一个定时器，用于触发超时发送
	timer := time.NewTimer(GroupMsgReadRecordDelayTime / 2)
	defer timer.Stop() // 确保定时器在方法结束时停止

	for {
		select {
		case <-m.done:
			// 如果接收到停止信号，退出循环
			return
		case <-timer.C:
			// 定时器触发，检查推送条件
			m.mu.Lock()

			// 计算剩余的延迟时间
			pushTime := m.pushTime
			val := GroupMsgReadRecordDelayTime - time.Since(pushTime)
			push := m.push
			logx.Infof("timer.C %v val %v", time.Now(), val)
			if val > 0 && m.count < GroupMsgReadRecordDelayCount || push == nil {
				if val > 0 {
					// 如果还没有达到条件，重置定时器
					timer.Reset(val)
				}

				// 未达标，释放锁并继续循环
				m.mu.Unlock()
				continue
			}

			// 达到条件，重置推送记录和计数
			m.pushTime = time.Now()
			m.push = nil
			m.count = 0
			timer.Reset(GroupMsgReadRecordDelayTime / 2)
			m.mu.Unlock()

			// 推送记录到频道
			logx.Infof("超过 合并的条件推送 %v ", push)
			m.pushCh <- push
		default:
			// 如果定时器未触发
			m.mu.Lock()

			// 检查是否超过了记录数量
			if m.count >= GroupMsgReadRecordDelayCount {
				push := m.push
				m.push = nil
				m.count = 0
				m.mu.Unlock()

				// 推送记录到频道
				logx.Infof("default 超过 合并的条件推送 %v ", push)
				m.pushCh <- push
				continue
			}

			// 检查是否处于空闲状态
			if m.isIdle() {
				m.mu.Unlock()
				// 使得msgReadTransfer释放
				m.pushCh <- &ws.Push{
					ChatType:       constants.GroupChatType,
					ConversationId: m.conversationId,
				}
				continue
			}
			m.mu.Unlock()

			// 等待一段时间后再次检查
			tempDelay := GroupMsgReadRecordDelayTime / 4
			if tempDelay > time.Second {
				tempDelay = time.Second
			}
			time.Sleep(tempDelay)
		}
	}
}

// IsIdle 方法检查当前状态是否为空闲。
func (m *groupMsgRead) IsIdle() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.isIdle()
}

// isIdle 方法检查当前状态是否为空闲。
func (m *groupMsgRead) isIdle() bool {
	// 计算距离推送时间的时间差
	pushTime := m.pushTime
	val := GroupMsgReadRecordDelayTime*2 - time.Since(pushTime)

	// 如果推送记录为空且计数器为零，认为是空闲状态
	if val <= 0 && m.push == nil && m.count == 0 {
		return true
	}

	return false
}

// clear 方法清理当前的推送记录和状态。
func (m *groupMsgRead) clear() {
	select {
	case <-m.done:
	default:
		// 关闭 done 通道以停止循环
		close(m.done)
	}

	// 清除推送记录
	m.push = nil
}
