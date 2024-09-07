package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/threading"
	"time"

	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

// 定义一个AckType类型的枚举，用于表示ACK的类型

type AckType int

const (
	NoAck    AckType = iota // 不需要ACK
	OnlyAck                 // 只需要一次ACK
	RigorAck                // 严格ACK，需要多次确认
)

// 将AckType类型转换为字符串表示

func (t AckType) ToString() string {
	switch t {
	case OnlyAck:
		return "OnlyAck"
	case RigorAck:
		return "RigorAck"
	}
	return "NoAck"
}

// 结构体表示WebSocket服务器

type Server struct {
	sync.RWMutex
	*threading.TaskRunner                // 任务执行器
	opt                   *serverOption  // 服务器选项
	authentication        Authentication // 认证接口

	routes map[string]HandlerFunc // 路由表
	addr   string                 // 服务器地址
	patten string                 // URL模式

	connToUser map[*Conn]string // 连接到用户ID的映射
	userToConn map[string]*Conn // 用户ID到连接的映射

	upgrader websocket.Upgrader // 升级HTTP连接到WebSocket的工具
	logx.Logger
}

// 创建一个新的Server实例

func NewServer(addr string, opts ...ServerOptions) *Server {
	opt := newServerOptions(opts...) // 解析服务器选项

	return &Server{
		routes:   make(map[string]HandlerFunc), //方法名和处理函数
		addr:     addr,
		patten:   opt.patten,
		opt:      &opt,
		upgrader: websocket.Upgrader{},

		authentication: opt.Authentication,

		connToUser: make(map[*Conn]string),
		userToConn: make(map[string]*Conn),

		Logger:     logx.WithContext(context.Background()),
		TaskRunner: threading.NewTaskRunner(opt.concurrency),
	}
}

// 处理WebSocket连接
//接受请求并处理请求

func (s *Server) ServerWs(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			s.Errorf("server handler ws recover err %v", r)
		}
	}()

	conn := NewConn(s, w, r) //创建一个新的WebSocket连接
	if conn == nil {
		return
	}

	if !s.authentication.Auth(w, r) { // 认证请求
		// conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprint("不具备访问权限")))
		s.Send(&Message{FrameType: FrameData, Data: fmt.Sprint("不具备访问权限")}, conn)
		conn.Close()
		return
	}

	// 记录连接
	s.addConn(conn, r)

	// 处理连接
	go s.handlerConn(conn)
}

// 根据连接对象执行任务处理
func (s *Server) handlerConn(conn *Conn) {
	uids := s.GetUsers(conn) // 获取用户ID
	conn.Uid = uids[0]       // 设置连接的用户ID

	// 处理任务
	go s.handlerWrite(conn)

	if s.isAck(nil) { // 检查是否需要ACK
		go s.readAck(conn)
	}

	for {
		// 获取请求消息
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.Errorf("websocket conn read message err %v", err)
			s.Close(conn)
			return
		}
		// 解析消息
		var message Message
		if err = json.Unmarshal(msg, &message); err != nil {
			s.Errorf("json unmarshal err %v, msg %v", err, string(msg))
			s.Close(conn)
			return
		}
		// 给客户端回复一个ACK
		// 依据消息进行处理
		if s.isAck(&message) {
			s.Infof("conn message read ack msg %v", message)
			conn.appendMsgMq(&message)
		} else {
			conn.message <- &message
		}
	}
}

// 判断消息是否需要ACK确认
func (s *Server) isAck(message *Message) bool {
	if message == nil {
		return s.opt.ack != NoAck // 如果消息为nil，检查服务器选项中是否设置了ACK
	}
	return s.opt.ack != NoAck && message.FrameType != FrameNoAck // 检查服务器选项和消息帧类型，确定是否需要ACK
}

// 读取消息的ACK llb
func (s *Server) readAck(conn *Conn) {
	for {
		select {
		case <-conn.done: // 如果连接关闭，退出循环
			s.Infof("close message ack uid %v ", conn.Uid)
			return
		default:
		}

		// 从队列中读取新的消息
		conn.messageMu.Lock()           // 加锁以保证消息队列的线程安全
		if len(conn.readMessage) == 0 { // 如果消息队列为空
			conn.messageMu.Unlock()            // 解锁
			time.Sleep(100 * time.Microsecond) // 增加短暂睡眠
			continue
		}

		// 读取第一条消息
		message := conn.readMessage[0]

		// 判断ACK的方式
		switch s.opt.ack {
		case OnlyAck:
			// 直接给客户端回复ACK
			s.Send(&Message{
				FrameType: FrameAck,
				Id:        message.Id,
				AckSeq:    message.AckSeq + 1,
			}, conn)
			// 进行业务处理
			// 把消息从队列中移除
			conn.readMessage = conn.readMessage[1:]
			conn.messageMu.Unlock() // 解锁

			conn.message <- message // 将消息发送到消息通道
		case RigorAck:
			// 先回复ACK
			if message.AckSeq == 0 {
				// 如果ACK序列号为0，表示未确认
				conn.readMessage[0].AckSeq++             // 增加ACK序列号
				conn.readMessage[0].ackTime = time.Now() // 记录ACK时间
				s.Send(&Message{
					FrameType: FrameAck,
					Id:        message.Id,
					AckSeq:    message.AckSeq,
				}, conn)
				s.Infof("message ack RigorAck send mid %v, seq %v , time%v", message.Id, message.AckSeq,
					message.ackTime)
				conn.messageMu.Unlock() // 解锁
				continue
			}

			// 再次验证ACK

			// 1. 客户端返回结果，再一次确认
			// 获取客户端的ACK序列号
			msgSeq := conn.readMessageSeq[message.Id]
			if msgSeq.AckSeq > message.AckSeq {
				// 如果客户端的ACK序列号大于消息的ACK序列号，表示确认成功
				conn.readMessage = conn.readMessage[1:]
				conn.messageMu.Unlock() // 解锁
				conn.message <- message // 将消息发送到消息通道
				s.Infof("message ack RigorAck success mid %v", message.Id)
				continue
			}

			// 2. 客户端没有确认，考虑是否超过了ACK的确认时间
			val := s.opt.ackTimeout - time.Since(message.ackTime)
			if !message.ackTime.IsZero() && val <= 0 {
				// 如果ACK确认时间超过设定的超时时间
				delete(conn.readMessageSeq, message.Id) // 删除消息序列号记录
				conn.readMessage = conn.readMessage[1:] // 从消息队列中移除消息
				conn.messageMu.Unlock()                 // 解锁
				continue
			}
			// 2.1 未超过，重新发送ACK
			conn.messageMu.Unlock() // 解锁
			s.Send(&Message{
				FrameType: FrameAck,
				Id:        message.Id,
				AckSeq:    message.AckSeq,
			}, conn)
			// 睡眠一定的时间后重试
			time.Sleep(3 * time.Second)
		}
	}
}

// 处理消息的任务
// 这段代码负责在 WebSocket 连接的生命周期内持续监听并处理来自客户端的消息，同时处理心跳和确认逻辑，确保连接稳定且消息得到正确处理。
func (s *Server) handlerWrite(conn *Conn) {
	// 进入无限循环，持续处理该连接的消息
	for {
		select {
		case <-conn.done:
			// 如果接收到 conn.done 的关闭信号，表示连接已关闭，退出循环
			return
		case message := <-conn.message:
			// 从 conn.message 通道中读取消息，开始处理消息

			switch message.FrameType {
			case FramePing:
				// 如果消息的类型是心跳消息（FramePing）
				// 处理心跳消息，回复一个心跳消息以维持连接
				s.Send(&Message{FrameType: FramePing}, conn)

			case FrameData:
				// 如果消息的类型是数据消息（FrameData）

				// 根据请求的方法名，在路由表中查找对应的处理函数
				if handler, ok := s.routes[message.Method]; ok {
					// 如果找到处理函数，调用该函数处理消息
					handler(s, conn, message)
				} else {
					// 如果没有找到对应的方法，表示该方法不存在
					// 发送一条错误消息，告知客户端请求的方法不存在
					s.Send(&Message{FrameType: FrameData, Data: fmt.Sprintf("不存在执行的方法 %v 请检查", message.Method)}, conn)
				}
			}

			// 检查消息是否需要进行确认（ACK）
			if s.isAck(message) {
				conn.messageMu.Lock()                   // 加锁，确保线程安全
				delete(conn.readMessageSeq, message.Id) // 删除该消息的 ACK 记录
				conn.messageMu.Unlock()                 // 解锁
			}
		}
	}
}

// 添加连接llb
func (s *Server) addConn(conn *Conn, req *http.Request) {
	uid := s.authentication.UserId(req) // 从请求中获取用户ID

	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	// 验证用户是否之前已登录
	if c := s.userToConn[uid]; c != nil {
		// 关闭之前的连接
		c.Close()
	}

	s.connToUser[conn] = uid // 记录连接与用户的映射
	s.userToConn[uid] = conn // 记录用户与连接的映射
}

// 根据用户ID获取连接
func (s *Server) GetConn(uid string) *Conn {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	return s.userToConn[uid] // 返回用户ID对应的连接
}

// 根据用户ID列表获取连接列表
func (s *Server) GetConns(uids ...string) []*Conn {
	if len(uids) == 0 {
		return nil
	}

	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	res := make([]*Conn, 0, len(uids))
	for _, uid := range uids {
		res = append(res, s.userToConn[uid]) // 添加对应的连接到结果列表
	}
	return res
}

// 根据连接列表获取用户ID列表llb

func (s *Server) GetUsers(conns ...*Conn) []string {

	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	var res []string
	if len(conns) == 0 {
		// 获取全部用户ID
		res = make([]string, 0, len(s.connToUser))
		for _, uid := range s.connToUser {
			res = append(res, uid) // 添加用户ID到结果列表
		}
	} else {
		// 获取部分用户ID
		res = make([]string, 0, len(conns))
		for _, conn := range conns {
			res = append(res, s.connToUser[conn]) // 添加连接对应的用户ID到结果列表
		}
	}

	return res
}

// 关闭连接 ok
func (s *Server) Close(conn *Conn) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	uid := s.connToUser[conn]
	if uid == "" {
		// 连接已关闭
		return
	}

	delete(s.connToUser, conn) // 删除连接与用户的映射
	delete(s.userToConn, uid)  // 删除用户与连接的映射

	conn.Close() // 关闭连接
}

// 根据用户ID发送消息
func (s *Server) SendByUserId(msg interface{}, sendIds ...string) error {
	if len(sendIds) == 0 {
		return nil // 如果没有用户ID，返回nil
	}

	return s.Send(msg, s.GetConns(sendIds...)...) // 获取连接并发送消息
}

// 发送消息到多个连接

func (s *Server) Send(msg interface{}, conns ...*Conn) error {
	if len(conns) == 0 {
		return nil // 如果没有连接，返回nil
	}

	data, err := json.Marshal(msg) // 将消息序列化为JSON格式
	if err != nil {
		return err // 序列化失败，返回错误
	}

	for _, conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			return err // 发送消息失败，返回错误
		}
	}

	return nil
}

// 添加路由

func (s *Server) AddRoutes(rs []Route) {
	for _, r := range rs {
		s.routes[r.Method] = r.Handler // 将路由方法和处理器添加到路由映射中
	}
}

// 启动服务器
func (s *Server) Start() {
	http.HandleFunc(s.patten, s.ServerWs)    // 注册WebSocket处理函数
	s.Info(http.ListenAndServe(s.addr, nil)) // 启动HTTP服务器
}

// 停止服务器
func (s *Server) Stop() {
	fmt.Println("停止服务") // 输出停止服务的消息
}
