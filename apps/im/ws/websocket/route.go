package websocket

type Route struct {
	Method  string      //方法
	Handler HandlerFunc //方法调用的函数
}

//服务对象，给请求方用户返回消息，也可以一个用户给另一个用户发送消息

type HandlerFunc func(srv *Server, conn *Conn, msg *Message)
