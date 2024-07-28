package immodels

import "github.com/zeromicro/go-zero/core/stores/mon"

// 确保 customConversationsModel 实现了 ConversationsModel 接口
var _ ConversationsModel = (*customConversationsModel)(nil)

type (
	ConversationsModel interface {
		conversationsModel
	}

	customConversationsModel struct {
		*defaultConversationsModel
	}
)

// NewConversationsModel 返回一个新的 Mongo 数据库模型
func NewConversationsModel(url, db, collection string) ConversationsModel {
	// 创建 Mongo 数据库连接
	conn := mon.MustNewModel(url, db, collection)
	// 返回自定义的会话模型实例
	return &customConversationsModel{
		defaultConversationsModel: newDefaultConversationsModel(conn),
	}
}

// MustConversationsModel 返回一个指定集合的 Mongo 数据库模型
func MustConversationsModel(url, db string) ConversationsModel {
	// 使用 "conversations" 集合作为默认集合名
	return NewConversationsModel(url, db, "conversations")
}
