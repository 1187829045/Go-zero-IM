package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

// Conn 表示 WebSocket 连接的结构体
type Conn struct {
	idleMu sync.Mutex // 用于保护 idle 变量的互斥锁

	Uid string // 用户 ID，标识 WebSocket 连接所属的用户

	*websocket.Conn // 嵌入的 Gorilla WebSocket 连接，提供 WebSocket 通信功能

	s *Server // 关联的服务器实例，指向 WebSocket 服务器对象

	idle              time.Time     // 记录连接的最后活动时间，用于判断连接是否空闲
	maxConnectionIdle time.Duration // 连接的最大空闲时间，超过该时间后连接会被关闭

	messageMu      sync.Mutex          // 用于保护消息队列的互斥锁，防止并发访问冲突
	readMessage    []*Message          // 存储未确认的消息，确保消息可靠传输
	readMessageSeq map[string]*Message // 存储消息 ID 和对应的消息，用于处理消息确认
	message        chan *Message       // 消息通道，用于异步发送和接收消息
	done           chan struct{}       // 用于关闭连接的信号通道，当通道关闭时，表示连接应被关闭
}

// 创建一个新的 WebSocket 连接
func NewConn(s *Server, w http.ResponseWriter, r *http.Request) *Conn {
	// 使用 Gorilla WebSocket 提供的 upgrader 将 HTTP 连接升级为 WebSocket 连接
	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		// 如果升级失败，记录错误并返回 nil，表示连接创建失败
		s.Errorf("upgrade err %v", err)
		return nil
	}

	// 创建并返回 Conn 实例，初始化各种字段
	conn := &Conn{
		Conn:              c,                            // WebSocket 连接对象
		s:                 s,                            // 关联的服务器实例
		idle:              time.Now(),                   // 初始化连接的最后活动时间为当前时间
		maxConnectionIdle: s.opt.maxConnectionIdle,      // 从服务器选项获取最大空闲时间
		readMessage:       make([]*Message, 0, 2),       // 初始化消息队列
		readMessageSeq:    make(map[string]*Message, 2), // 初始化消息 ID 和消息映射
		message:           make(chan *Message, 1),       // 初始化消息通道
		done:              make(chan struct{}),          // 初始化关闭信号通道
	}

	// 启动连接的心跳检测 goroutine，定期检查连接的活动状态
	go conn.keepalive()
	return conn
}

// appendMsgMq 将消息追加到消息队列中
func (c *Conn) appendMsgMq(msg *Message) {
	// 锁定消息队列，防止并发修改
	c.messageMu.Lock()
	defer c.messageMu.Unlock()

	// 检查消息是否已存在于 readMessageSeq 中（通过消息 ID）
	if m, ok := c.readMessageSeq[msg.Id]; ok {
		// 消息已存在，检查消息的确认序号
		if len(c.readMessage) == 0 {
			// 队列中没有该消息，则直接返回
			return
		}

		// 如果消息的确认序号大于或等于已有记录的确认序号，则不处理
		if m.AckSeq >= msg.AckSeq {
			return
		}

		// 更新消息记录为新的消息
		c.readMessageSeq[msg.Id] = msg
		return
	}

	// 如果消息是 ACK 类型（确认消息），则不处理直接返回
	if msg.FrameType ==
		FrameAck {
		return
	}

	// 将消息追加到消息队列中，并更新消息序列映射
	c.readMessage = append(c.readMessage, msg)
	c.readMessageSeq[msg.Id] = msg
}

// ReadMessage 读取 WebSocket 消息
func (c *Conn) ReadMessage() (messageType int, p []byte, err error) {
	// 调用 Gorilla WebSocket 的 ReadMessage 方法读取消息
	messageType, p, err = c.Conn.ReadMessage()

	// 更新连接的最后活动时间为初始时间（表示不再空闲）
	c.idleMu.Lock()
	defer c.idleMu.Unlock()
	c.idle = time.Time{} // 将 idle 设置为零值，表示连接是非空闲状态
	return
}

// WriteMessage 发送 WebSocket 消息
func (c *Conn) WriteMessage(messageType int, data []byte) error {
	// 锁定空闲时间，防止并发修改
	c.idleMu.Lock()
	defer c.idleMu.Unlock()

	// 调用 Gorilla WebSocket 的 WriteMessage 方法发送消息
	err := c.Conn.WriteMessage(messageType, data)
	// 记录发送消息的时间为最后活动时间
	c.idle = time.Now()
	return err
}

// Close 关闭 WebSocket 连接
func (c *Conn) Close() error {
	// 选择关闭 done 通道，标识连接即将关闭
	select {
	case <-c.done: // 如果 done 已经关闭，不执行任何操作
	default:
		close(c.done) // 关闭 done 通道，发出关闭信号
	}

	// 调用 Gorilla WebSocket 的 Close 方法关闭 WebSocket 连接
	return c.Conn.Close()
}

// keepalive 处理连接的心跳检测
func (c *Conn) keepalive() {
	// 初始化一个定时器，用于检测连接的最大空闲时间
	idleTimer := time.NewTimer(c.maxConnectionIdle)
	defer func() {
		idleTimer.Stop() // 当心跳检测结束时，停止定时器
	}()

	for {
		select {
		// 当定时器到期（即连接的空闲时间超过最大空闲时间）时，执行以下逻辑
		case <-idleTimer.C:
			// 锁定空闲时间，防止并发修改
			c.idleMu.Lock()
			idle := c.idle     // 读取最后活动时间
			if idle.IsZero() { // 如果连接非空闲状态（idle 为零值）
				c.idleMu.Unlock()
				idleTimer.Reset(c.maxConnectionIdle) // 重置定时器，重新开始检测
				continue
			}

			// 计算剩余的空闲时间
			val := c.maxConnectionIdle - time.Since(idle)
			c.idleMu.Unlock()

			if val <= 0 {
				// 如果剩余时间小于等于零，表示空闲时间超过最大空闲时间
				// 通过服务器实例优雅地关闭连接
				c.s.Close(c)
				return
			}

			// 如果还有剩余时间，重置定时器继续检测
			idleTimer.Reset(val)
		// 当连接关闭时（done 通道关闭），结束心跳检测
		case <-c.done:
			return
		}
	}
}
