/**
 * @author: dn-jinmin/dn-jinmin
 * @doc:
 */

package ctxdata

import "context"

// GetUId 从 context.Context 中提取用户 ID (Identify)
func GetUId(ctx context.Context) string {
	// 尝试从上下文中获取 Identify 对应的值，并断言为字符串类型
	if u, ok := ctx.Value(Identify).(string); ok {
		// 如果断言成功，返回获取到的字符串值
		return u
	}
	// 如果断言失败，返回空字符串
	return ""
}
