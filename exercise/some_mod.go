package exercise

//用于记录一些常用的结构体
import (
	"github.com/zeromicro/go-zero/core/service"
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

)
