package wuid

import (
	"database/sql"
	"fmt"
	"github.com/bwmarrin/snowflake" // 雪花算法库
	"github.com/edwingeng/wuid/mysql/wuid"
	"log"
	"sort"
	"strconv"
)

// Node 是雪花节点
var Node *snowflake.Node

func init() {
	//雪花算法init
	err := InitIDGen(1)
	if err != nil {
		panic(err)
	}
}

// InitIDGen 初始化雪花 ID 生成器
func InitIDGen(nodeID int64) error {
	var err error
	Node, err = snowflake.NewNode(nodeID)
	if err != nil {
		return err
	}
	return nil
}

// id生成工具
// id可以使用直接自增的方式处理，但需考虑项目后期数据量增长(拆分数据库)
// 使用wuid库：可以生成全局唯一id，并记录id存储方式(mysql、redis均可)

var w *wuid.WUID

func Init(dsn string) {

	newDB := func() (*sql.DB, bool, error) {
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			return nil, false, err
		}
		return db, true, nil
	}

	w = wuid.NewWUID("default", nil)
	_ = w.LoadH28FromMysql(newDB, "wuid")
}

func GenUid(dsn string) string {
	if w == nil {
		Init(dsn)
	}

	return fmt.Sprintf("%#016x", w.Next())
}

func CombineId(aid, bid string) string {
	ids := []string{aid, bid}

	sort.Slice(ids, func(i, j int) bool {
		a, _ := strconv.ParseUint(ids[i], 0, 64)
		b, _ := strconv.ParseUint(ids[j], 0, 64)
		return a < b
	})

	return fmt.Sprintf("%s_%s", ids[0], ids[1])
}

// GenerateID 生成一个新的唯一 ID
func GenerateID() string {
	if Node == nil {
		log.Fatal("Snowflake Node is not initialized")
	}
	return Node.Generate().String()
}
