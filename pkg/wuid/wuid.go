package wuid

import (
	"database/sql"
	"fmt"
	"github.com/edwingeng/wuid/mysql/wuid"
	"sort"
	"strconv"
)

// 定义一个 WUID 实例变量
var w *wuid.WUID

// Init 函数用于初始化 WUID 实例
func Init(dsn string) {
	// 定义一个函数，用于创建和返回一个新的数据库连接
	newDB := func() (*sql.DB, bool, error) {
		// 使用提供的数据源名称 (DSN) 打开一个新的数据库连接
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			// 如果打开连接出错，返回错误
			return nil, false, err
		}
		// 返回数据库连接对象和成功状态
		return db, true, nil
	}

	// 创建一个新的 WUID 实例
	w = wuid.NewWUID("default", nil)
	// 从 MySQL 数据库中加载 H28 组件
	_ = w.LoadH28FromMysql(newDB, "wuid")
}

// GenUid 函数用于生成一个唯一 ID
func GenUid(dsn string) string {
	// 如果 WUID 实例尚未初始化，则进行初始化
	if w == nil {
		Init(dsn)
	}

	// 返回格式化的唯一 ID
	return fmt.Sprintf("%#016x", w.Next())
}

// CombineId 函数用于组合两个 ID
func CombineId(aid, bid string) string {
	// 将两个 ID 放入切片
	ids := []string{aid, bid}

	// 对 ID 切片进行排序
	sort.Slice(ids, func(i, j int) bool {
		// 将 ID 字符串解析为无符号整数
		a, _ := strconv.ParseUint(ids[i], 0, 64)
		b, _ := strconv.ParseUint(ids[j], 0, 64)
		// 返回比较结果，用于排序
		return a < b
	})

	// 返回格式化后的组合 ID
	return fmt.Sprintf("%s_%s", ids[0], ids[1])
}
