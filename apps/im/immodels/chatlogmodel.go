package immodels

import "github.com/zeromicro/go-zero/core/stores/mon"

//这段代码的作用是定义并实现一个与 MongoDB 交互的 ChatLogModel 数据模型，用于操作特定集合中的数据（如 chat_log 集合）。
//通过使用 go-zero 的 mon 包，这个模型可以方便地与 MongoDB 进行数据存储和检索操作。
// 确保 customChatLogModel 实现了 ChatLogModel 接口

var _ ChatLogModel = (*customChatLogModel)(nil)

type (
	ChatLogModel interface {
		chatLogModel
	}

	customChatLogModel struct {
		*defaultChatLogModel
	}
)

// NewChatLogModel 返回一个用于 MongoDB 的模型实例。
func NewChatLogModel(url, db, collection string) ChatLogModel {
	// 连接到 MongoDB 并创建一个模型实例
	conn := mon.MustNewModel(url, db, collection)
	return &customChatLogModel{
		// 初始化 defaultChatLogModel
		defaultChatLogModel: newDefaultChatLogModel(conn),
	}
}

// MustChatLogModel 返回一个用于 chat_log 集合的模型实例。
func MustChatLogModel(url, db string) ChatLogModel {
	// 调用 NewChatLogModel，指定集合名称为 chat_log
	return NewChatLogModel(url, db, "chat_log")
}
