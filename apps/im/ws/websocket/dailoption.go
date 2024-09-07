package websocket

import "net/http"

// 定义了一个函数类型，用于配置拨号选项

type DailOptions func(option *dailOption)

// 结构体用于存储拨号选项
type dailOption struct {
	pattern string      // WebSocket 连接的路径模式
	header  http.Header // 连接请求的 HTTP 头
}

// 使用给定的 DailOptions 函数创建并返回一个 dailOption 实例
func newDailOptions(opts ...DailOptions) dailOption {
	// 初始化 dailOption 的默认值
	o := dailOption{
		pattern: "/ws", // 默认的 WebSocket 路径模式
		header:  nil,   // 默认不设置 HTTP 头
	}

	// 应用所有的 DailOptions 函数来定制 dailOption
	for _, opt := range opts {
		opt(&o)
	}

	return o
}

// 返回一个 DailOptions 函数，用于设置 WebSocket 连接的路径模式
func WithClientPatten(pattern string) DailOptions {
	return func(opt *dailOption) {
		opt.pattern = pattern
	}
}

// 返回一个 DailOptions 函数，用于设置连接请求的 HTTP 头
func WithClientHeader(header http.Header) DailOptions {
	return func(opt *dailOption) {
		opt.header = header
	}
}
