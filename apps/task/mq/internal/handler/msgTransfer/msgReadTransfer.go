package msgTransfer

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"llb-chat/apps/im/ws/ws"
	"llb-chat/apps/task/mq/internal/svc"
	"llb-chat/apps/task/mq/mq"
	"llb-chat/pkg/bitmap"
	"llb-chat/pkg/constants"
	"sync"
	"time"
)

// 定义全局变量，控制消息读取记录的延迟时间和记录数量
var (
	GroupMsgReadRecordDelayTime  = time.Second
	GroupMsgReadRecordDelayCount = 10
)

// 定义常量，用于选择消息处理的策略
const (
	GroupMsgReadHandlerAtTransfer = iota
	GroupMsgReadHandlerDelayTransfer
)

// 包含消息读取处理相关字段
type MsgReadTransfer struct {
	*baseMsgTransfer //包含基础的消息传输功能

	cache.Cache //包含缓存功能

	mu sync.Mutex // 保护并发访问的锁

	groupMsgs map[string]*groupMsgRead // 存储群聊消息读取记录
	push      chan *ws.Push            // 推送消息的通道
}

func NewMsgReadTransfer(svc *svc.ServiceContext) kq.ConsumeHandler {
	// 创建 MsgReadTransfer 实例并初始化
	m := &MsgReadTransfer{
		baseMsgTransfer: NewBaseMsgTransfer(svc),           // 初始化基础消息传输，传入服务上下文
		groupMsgs:       make(map[string]*groupMsgRead, 1), // 初始化群聊消息记录映射，容量设置为1
		push:            make(chan *ws.Push, 1),            // 初始化推送通道，容量设置为1
	}

	// 配置消息读取处理策略，根据服务配置修改延迟计数和时间
	if svc.Config.MsgReadHandler.GroupMsgReadHandler != GroupMsgReadHandlerAtTransfer {
		// 如果配置中定义的处理策略不是在传输时处理
		if svc.Config.MsgReadHandler.GroupMsgReadRecordDelayCount > 0 {
			// 设置消息读取记录延迟计数
			GroupMsgReadRecordDelayCount = svc.Config.MsgReadHandler.GroupMsgReadRecordDelayCount
		}

		if svc.Config.MsgReadHandler.GroupMsgReadRecordDelayTime > 0 {
			// 设置消息读取记录延迟时间
			GroupMsgReadRecordDelayTime = time.Duration(svc.Config.MsgReadHandler.GroupMsgReadRecordDelayTime) * time.Second
		}
	}

	// 启动一个 goroutine 来处理消息传输
	go m.transfer()

	// 返回 MsgReadTransfer 实例，实现 kq.ConsumeHandler 接口
	return m
}

//处理消费的消息
//主要处理消息读取状态的更新。这个逻辑通常发生在用户读取消息时，需要更新相应的聊天日志状态，并生成推送记录。

func (m *MsgReadTransfer) Consume(key, value string) error {
	m.Info("MsgReadTransfer ", value)

	var (
		data mq.MsgMarkRead // 消息读取数据
		ctx  = context.Background()
	)
	// 解析消息数据
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}

	// 更新聊天日志读取状态
	readRecords, err := m.UpdateChatLogRead(ctx, &data)
	if err != nil {
		return err
	}

	// 创建推送记录
	push := &ws.Push{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		ContentType:    constants.ContentMakeRead,
		ReadRecords:    readRecords,
	}

	switch data.ChatType {
	case constants.SingleChatType:
		// 单聊直接推送
		m.push <- push
	case constants.GroupChatType:
		// 群聊处理
		if m.svcCtx.Config.MsgReadHandler.GroupMsgReadHandler == GroupMsgReadHandlerAtTransfer {
			m.push <- push
		}

		m.mu.Lock()
		defer m.mu.Unlock()

		push.SendId = ""

		if _, ok := m.groupMsgs[push.ConversationId]; ok {
			m.Infof("merge push %v", push.ConversationId)
			// 合并请求
			m.groupMsgs[push.ConversationId].mergePush(push)
		} else {
			m.Infof("newGroupMsgRead push %v", push.ConversationId)
			// 新建群聊消息读取记录
			m.groupMsgs[push.ConversationId] = newGroupMsgRead(push, m.push)
		}
	}

	return nil
}

// 更新聊天日志的已读状态
func (m *MsgReadTransfer) UpdateChatLogRead(ctx context.Context, data *mq.MsgMarkRead) (map[string]string, error) {

	res := make(map[string]string) // 记录更新后的结果

	// 根据消息 ID 获取聊天日志
	chatLogs, err := m.svcCtx.ChatLogModel.ListByMsgIds(ctx, data.MsgIds)
	if err != nil {
		return nil, err
	}

	// 处理已读状态
	for _, chatLog := range chatLogs {
		switch chatLog.ChatType {
		case constants.SingleChatType:
			// 单聊消息设置为已读
			chatLog.ReadRecords = []byte{1}
		case constants.GroupChatType:
			// 群聊消息更新已读记录
			readRecords := bitmap.Load(chatLog.ReadRecords)
			readRecords.Set(data.SendId)
			chatLog.ReadRecords = readRecords.Export()
		}

		// 更新记录到结果映射
		res[chatLog.ID.Hex()] = base64.StdEncoding.EncodeToString(chatLog.ReadRecords)

		// 更新聊天日志
		err = m.svcCtx.ChatLogModel.UpdateMakeRead(ctx, chatLog.ID, chatLog.ReadRecords)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

// 处理推送消息的合并和发送
func (m *MsgReadTransfer) transfer() {
	for push := range m.push {
		if push.RecvId != "" || len(push.RecvIds) > 0 {
			// 处理单聊和群聊的消息传输
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

		// 清理空闲状态的群聊消息记录
		m.mu.Lock()
		if _, ok := m.groupMsgs[push.ConversationId]; ok && m.groupMsgs[push.ConversationId].IsIdle() {
			m.groupMsgs[push.ConversationId].clear()
			delete(m.groupMsgs, push.ConversationId)
		}

		m.mu.Unlock()
	}
}
