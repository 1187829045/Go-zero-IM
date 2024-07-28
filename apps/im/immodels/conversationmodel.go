package immodels

import "github.com/zeromicro/go-zero/core/stores/mon"

// 确保 customConversationModel 实现了 ConversationModel 接口
var _ ConversationModel = (*customConversationModel)(nil)

type (
	ConversationModel interface {
		conversationModel
	}

	customConversationModel struct {
		*defaultConversationModel
	}
)

// NewConversationModel 返回一个新的 Mongo 数据库模型
func NewConversationModel(url, db, collection string) ConversationModel {
	// 创建 Mongo 数据库连接
	conn := mon.MustNewModel(url, db, collection)
	// 返回自定义的会话模型实例
	return &customConversationModel{
		defaultConversationModel: newDefaultConversationModel(conn),
	}
}

// MustConversationModel 返回一个指定集合的 Mongo 数据库模型
func MustConversationModel(url, db string) ConversationModel {
	// 使用 "conversation" 集合作为默认集合名
	return NewConversationModel(url, db, "conversation")
}
