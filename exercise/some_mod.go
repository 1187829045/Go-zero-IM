package exercise

//用于记录一些常用的结构体
import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
	"llb-chat/pkg/constants"
	"strings"
	"time"
)

type (
	// RedisConf 是一个 Redis 配置结构体。
	RedisConf struct {
		Host     string // Redis 服务的主机地址。
		Type     string `json:",default=node,options=node|cluster"` // Redis 类型，默认为节点模式，可以选择节点模式或集群模式。
		Pass     string `json:",optional"`                          // Redis 服务的密码，非必须项。
		Tls      bool   `json:",optional"`                          // 是否启用 TLS 加密连接，非必须项。
		NonBlock bool   `json:",default=true"`                      // 是否启用非阻塞模式，默认为启用。
		// PingTimeout 是 ping Redis 的超时时间。
		PingTimeout time.Duration `json:",default=1s"` // ping 操作的超时时间，默认为 1 秒。
	}

	// RedisKeyConf 是一个带有键的 Redis 配置结构体。
	RedisKeyConf struct {
		RedisConf        // 内嵌 RedisConf 结构体。
		Key       string // 用于存储 Redis 键的字符串。
	}

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	PrivateKeyConf struct {
		Fingerprint string
		KeyFile     string
	}
	SignatureConf struct {
		Strict      bool          `json:",default=false"`
		Expiry      time.Duration `json:",default=1h"`
		PrivateKeys []PrivateKeyConf
	}
	MiddlewaresConf struct {
		Trace      bool `json:",default=true"`
		Log        bool `json:",default=true"`
		Prometheus bool `json:",default=true"`
		MaxConns   bool `json:",default=true"`
		Breaker    bool `json:",default=true"`
		Shedding   bool `json:",default=true"`
		Timeout    bool `json:",default=true"`
		Recover    bool `json:",default=true"`
		Metrics    bool `json:",default=true"`
		MaxBytes   bool `json:",default=true"`
		Gunzip     bool `json:",default=true"`
	}

	RestConf struct {
		service.ServiceConf        // 嵌入 ServiceConf 结构体，包含服务的通用配置
		Host                string `json:",default=0.0.0.0"` // 服务监听的主机地址，默认为 0.0.0.0
		Port                int    // 服务监听的端口
		CertFile            string `json:",optional"`        // SSL 证书文件路径，可选
		KeyFile             string `json:",optional"`        // SSL 密钥文件路径，可选
		Verbose             bool   `json:",optional"`        // 是否开启详细日志，可选
		MaxConns            int    `json:",default=10000"`   // 最大连接数，默认值为 10000
		MaxBytes            int64  `json:",default=1048576"` // 请求体的最大字节数，默认值为 1048576（1MB）
		// 毫秒为单位的超时时间，默认值为 3000 毫秒（3 秒）
		Timeout int64 `json:",default=3000"`
		// CPU 使用率的阈值，范围为 0 到 1000，默认值为 900
		CpuThreshold int64 `json:",default=900,range=[0:1000)"`
		// 签名配置，可选
		Signature SignatureConf `json:",optional"`
		// 中间件配置，所有项目都有默认值
		Middlewares MiddlewaresConf
		// 跟踪中间件的路径黑名单，可选
		TraceIgnorePaths []string `json:",optional"`
	}
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	MsgChatTransfer struct {
		MsgId string `mapstructure:"msgId"`

		ConversationId     string `json:"conversationId"`
		constants.ChatType `json:"chatType"`
		SendId             string   `json:"sendId"`
		RecvId             string   `json:"recvId"`
		RecvIds            []string `json:"recvIds"` // 接收者ID列表，表示消息接收方的用户ID列表（群聊）
		SendTime           int64    `json:"sendTime"`

		constants.MType `json:"mType"` // 消息类型，使用常量定义，可能是文本、图片、视频等
		Content         string         `json:"content"` // 消息内容，表示消息的实际内容
	}

	// MsgMarkRead 结构体表示已读消息的标记信息
	MsgMarkRead struct {
		constants.ChatType `json:"chatType"` // 聊天类型，使用常量定义，可能是单聊或群聊
		ConversationId     string            `json:"conversationId"` // 会话ID，表示消息所属的会话
		SendId             string            `json:"sendId"`         // 发送者ID，表示消息发送方的用户ID
		RecvId             string            `json:"recvId"`         // 接收者ID，表示消息接收方的用户ID
		MsgIds             []string          `json:"msgIds"`         // 消息ID列表，表示已读的消息ID列表
	}
)

var (
	groupMembersFieldNames          = builder.RawFieldNames(&GroupMembers{})
	groupMembersRows                = strings.Join(groupMembersFieldNames, ",")
	groupMembersRowsExpectAutoSet   = strings.Join(stringx.Remove(groupMembersFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	groupMembersRowsWithPlaceHolder = strings.Join(stringx.Remove(groupMembersFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheGroupMembersIdPrefix = "cache:groupMembers:id:"
)

type (
	groupMembersModel interface {
		Insert(ctx context.Context, session sqlx.Session, data *GroupMembers) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*GroupMembers, error)
		FindByGroudIdAndUserId(ctx context.Context, userId, groupId string) (*GroupMembers, error)
		ListByUserId(ctx context.Context, userId string) ([]*GroupMembers, error)
		ListByGroupId(ctx context.Context, groupId string) ([]*GroupMembers, error)
		Update(ctx context.Context, data *GroupMembers) error
		Delete(ctx context.Context, id int64) error
	}

	defaultGroupMembersModel struct {
		sqlc.CachedConn
		table string
	}

	GroupMembers struct {
		Id          int64          `db:"id"`
		GroupId     string         `db:"group_id"`
		UserId      string         `db:"user_id"`
		RoleLevel   int            `db:"role_level"`
		JoinTime    sql.NullTime   `db:"join_time"`
		JoinSource  sql.NullInt64  `db:"join_source"`
		InviterUid  sql.NullString `db:"inviter_uid"`
		OperatorUid string         `db:"operator_uid"`
	}
)

func (m *defaultGroupMembersModel) ListByUserId(ctx context.Context, userId string) ([]*GroupMembers, error) {
	// 构建查询语句
	query := fmt.Sprintf("select %s from %s where `user_id` = ?", groupMembersRows, m.table)
	var resp []*GroupMembers
	// 执行查询操作，不使用缓存
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}
