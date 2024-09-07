package websocket

//鉴权
import (
	"fmt"
	"net/http"
	"time"
)

type Authentication interface {
	Auth(w http.ResponseWriter, r *http.Request) bool
	UserId(r *http.Request) string
}

type authentication struct{}

func (*authentication) Auth(w http.ResponseWriter, r *http.Request) bool {
	return true
}
func (*authentication) UserId(r *http.Request) string {
	// 这个方法是 `authentication` 类型的一个方法，接受一个 `*http.Request` 类型的参数 `r`，返回一个字符串。
	query := r.URL.Query()
	// 从请求 `r` 中获取 URL 的查询参数，返回一个 `url.Values` 类型的 map，包含所有的查询参数。

	if query != nil && query["userId"] != nil {
		// 检查查询参数 `query` 是否不为空，且其中是否包含 `userId` 这个键。
		return fmt.Sprintf("%v", query["userId"])
		// 如果 `userId` 存在，则返回其对应的值。
		// 由于 `query["userId"]` 是一个字符串切片，这里用 `fmt.Sprintf` 将其格式化为字符串。
	}

	return fmt.Sprintf("%v", time.Now().UnixMilli())
	// 如果 `userId` 不存在，返回当前时间的 Unix 毫秒时间戳，转换为字符串形式。
}
