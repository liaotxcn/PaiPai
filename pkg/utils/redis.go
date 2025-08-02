package utils

import (
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func InitRDB(dsn string, pass string) *redis.Client {
	//fmt.Printf("init redis on dsn %v, password %v \n", dsn, pass)
	rdb = redis.NewClient(&redis.Options{
		Addr:     dsn,
		Password: pass,
		DB:       0,
	})
	return rdb
}
