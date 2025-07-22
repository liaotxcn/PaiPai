package wuid

import (
	"database/sql"
	"fmt"
	"github.com/edwingeng/wuid/mysql/wuid"
)

// id生成工具
// id可以使用直接自增的方式处理，但需考虑项目后期数据量增长(拆分数据库)
// 使用wuid库：可以生成全局唯一id，并记录id存储方式(mysql、redis均可)

var w *wuid.WUID

func Init(dns string) {
	newDB := func() (*sql.DB, bool, error) {
		db, err := sql.Open("mysql", dns)
		if err != nil {
			return nil, false, err
		}
		return db, true, nil
	}

	w = wuid.NewWUID("default", nil)
	_ = w.LoadH28FromMysql(newDB, "wuid")
}

func GenUid(dns string) string {
	if w == nil {
		Init(dns)
	}
	return fmt.Sprintf("%#016x", w.Next())
}
