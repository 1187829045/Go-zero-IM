package jwt_exercise

import (
	"github.com/golang-jwt/jwt/v4"
	"io/ioutil"
	"log"
	"testing"
)

// 对称加密
// jwt.RegisteredClaims包含一下重要字段
type RegisteredClaims struct {
	Issuer    string           `json:"iss,omitempty"` // 签发者，表示这个 JWT 的发布者
	Subject   string           `json:"sub,omitempty"` // 主题，通常是 JWT 面向的用户或实体
	Audience  jwt.ClaimStrings `json:"aud,omitempty"` // 受众，表示 JWT 的接收者
	ExpiresAt *jwt.NumericDate `json:"exp,omitempty"` // 过期时间，JWT 的有效期
	NotBefore *jwt.NumericDate `json:"nbf,omitempty"` // 在此时间之前，JWT 无效
	IssuedAt  *jwt.NumericDate `json:"iat,omitempty"` // 签发时间，表示 JWT 的发布时间
	ID        string           `json:"jti,omitempty"` // JWT ID，唯一标识符，防止重放攻击
}

func TestHs256(t *testing.T) {
	// 定义一个 User 结构体，用于存储用户信息
	type User struct {
		Id   int64
		Name string
	}

	// 定义一个 UserClaims 结构体，包含用户信息和注册声明
	type UserClaims struct {
		User                 User
		jwt.RegisteredClaims // 注册声明，来自 jwt 包
	}

	// 1. 使用 jwt.NewWithClaims 生成一个新的 Token
	user := User{
		Id:   101,   // 设置用户ID
		Name: "llb", // 设置用户名
	}
	// 创建 UserClaims 实例，包含用户信息和空的注册声明
	userClaims := UserClaims{
		User:             user,                   // 设置用户信息
		RegisteredClaims: jwt.RegisteredClaims{}, // 设置注册声明
	}
	// 使用 HS256 签名方法生成一个新的 Token，传入用户声明
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)

	// 2. 对 Token 进行签名，生成加密的 Token 字符串
	mySigningKey := []byte("ushjlwmwnwht")      // 定义用于签名的密钥
	ss, err := token.SignedString(mySigningKey) // 使用密钥对 Token 进行签名
	// 打印签名后的 Token 字符串和可能的错误信息
	t.Log(ss, err)
}

// 解密
func TestHs256Parse(t *testing.T) {
	// 定义一个 JWT token 字符串，这里是一个示例 token
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VyIjp7IklkIjoxMDEsIk5hbWUiOiJoaXNoZW5nIn19.ij1kWID03f_CiELe0fPLZJ-Y64dkf2nDE-f6nGERBSE"

	// 定义一个 User 结构体，用于存储用户信息
	type User struct {
		Id   int64  // 用户ID
		Name string // 用户名
	}

	// 定义一个 UserClaims 结构体，包含用户信息和注册声明
	type UserClaims struct {
		User                 User // 用户信息
		jwt.RegisteredClaims      // 注册声明，来自 jwt 包
	}

	// 使用 jwt.ParseWithClaims 解析 token
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 这个回调函数提供了用于验证 token 签名的密钥
		// 这里返回预设的密钥 "ushjlwmwnwht"
		return []byte("ushlgende"), nil
	})

	// 检查 token 是否成功解析并且有效
	if userClaims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		// 如果解析成功且 token 有效，打印 userClaims 和注册声明中的 Issuer
		t.Log(userClaims, userClaims.RegisteredClaims.Issuer)
	} else {
		// 如果解析失败或者 token 无效，打印错误信息
		t.Log(err)
	}
}
func TestHs256Parse2(t *testing.T) {
	// 定义一个 JWT token 字符串，这里是一个示例 token
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VyIjp7IklkIjoxMDEsIk5hbWUiOiJoaXNoZW5nIn19.ij1kWID03f_CiELe0fPLZJ-Y64dkf2nDE-f6nGERBSE"

	// 定义一个 User 结构体，用于存储用户信息
	type User struct {
		Id   int64
		Name string // 用户名
	}

	// 定义一个 UserClaims 结构体，包含用户信息和注册声明
	type UserClaims struct {
		User                 User
		jwt.RegisteredClaims // 注册声明，来自 jwt 包
	}

	// 创建一个新的 JWT 解析器
	parser := jwt.NewParser()

	// 使用 parser.ParseWithClaims 解析 token
	token, err := parser.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 这个回调函数提供了用于验证 token 签名的密钥
		// 这里返回预设的密钥 "ushjlwmwnwht"
		return []byte("ushjlwmwnwht"), nil
	})

	// 检查 token 是否成功解析并且有效
	if userClaims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		// 如果解析成功且 token 有效，打印 userClaims 和注册声明中的 Issuer
		t.Log(userClaims, userClaims.RegisteredClaims.Issuer)
	} else {
		// 如果解析失败或者 token 无效，打印错误信息
		t.Log(err)
	}
}

//非对称加密

func TestRs256(t *testing.T) {
	// 定义一个 User 结构体，用于存储用户信息
	type User struct {
		Id   int64
		Name string // 用户名
	}

	// 定义一个 UserClaims 结构体，包含用户信息和注册声明
	type UserClaims struct {
		User                 User
		jwt.RegisteredClaims // 注册声明，来自 jwt 包
	}

	// 1. 使用 jwt.NewWithClaims 生成一个新的 Token
	user := User{
		Id:   101,
		Name: "hisheng", // 设置用户名
	}
	// 创建 UserClaims 实例，包含用户信息和空的注册声明
	userClaims := UserClaims{
		User:             user,
		RegisteredClaims: jwt.RegisteredClaims{},
	}

	// 读取 PEM 格式的私钥文件
	privateKeyDataPem, err := ioutil.ReadFile("private-key.pem")
	if err != nil {
		log.Fatal(err) // 如果读取文件出错，打印错误并退出程序
	}

	// 从 PEM 数据解析 RSA 私钥
	privateKeyData, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyDataPem)
	if err != nil {
		log.Fatal(err) // 如果解析私钥出错，打印错误并退出程序
	}

	// 2. 使用 RS256 签名方法生成 Token，并用 RSA 私钥对其进行签名
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, userClaims).SignedString(privateKeyData)
	// 打印生成的 Token 和可能的错误信息
	t.Log(token, err)
}
func TestRs256Parse(t *testing.T) {
	// 定义一个 JWT token 字符串，这里是一个示例 token
	tokenString := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VyIjp7IklkIjoxMDEsIk5hbWUiOiJoaXNoZW5nIn19.GOS-d7iwaLDCSaSsBpArbtH-3JlD1KFNoJeyZjQ6Xv4XySo599WY784mVj-rR9kxtLYmkdAG3dPG9am6NZubHBLYWDi9b1w5gXcC7o3nAZirFGgS7bwf-7DmEetwUWzJZe-QXr1CIiVyHPU2ZCketYkIqgMGixVNfxfS5pJ5LhiUM_7J9Zlf5Kq15P9Y7U35Nu25-UXGgs-BHvH-Db6PJ4vHZq-dla_yuqRN276JBxdc24SEnML_iOHmgXVOXEWjtC_p6LmsP0VGMqwXAHN4FH0XbMzpQtTQKJgskggI41T1Ruzb9zpzJdmiL2DkeMgbvc0xVzV3CjM_AA5cUqcZaQ"

	// 定义一个 User 结构体，用于存储用户信息
	type User struct {
		Id   int64  // 用户ID
		Name string // 用户名
	}

	// 定义一个 UserClaims 结构体，包含用户信息和注册声明
	type UserClaims struct {
		User                 User // 用户信息
		jwt.RegisteredClaims      // 注册声明，来自 jwt 包
	}

	// 使用 jwt.ParseWithClaims 解析 token
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 从文件系统读取 PEM 格式的公钥文件
		publicKeyDataPem, err := ioutil.ReadFile("public-key.pem")
		if err != nil {
			log.Fatal(err) // 如果读取文件出错，打印错误并退出程序
		}
		// 从 PEM 数据解析 RSA 公钥
		return jwt.ParseRSAPublicKeyFromPEM(publicKeyDataPem)
	})

	// 检查 token 是否成功解析为 UserClaims 类型，并且 token 是有效的
	if userClaims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		// 如果解析成功且 token 有效，打印 userClaims 和注册声明中的 Issuer
		t.Log(userClaims, userClaims.RegisteredClaims.Issuer)
	} else {
		// 如果解析失败或者 token 无效，打印错误信息
		t.Log(err)
	}
}
