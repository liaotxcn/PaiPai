package main

import (
	"fmt"
	"github.com/HYY-yu/sail-client"
	"time"
) // sail客户端库

// sail应用实例

type Config struct {
	Name     string
	Host     string
	Port     string
	Mode     string
	Datebase string

	UserRpc struct {
		Etcd struct {
			Hosts []string
			Key   string
		}
	}
	Redisx struct {
		Hosts []string
		Pass  string
	}
	JwtAuth struct {
		AccessSecret string
	}
}

func main() {
	var cfg Config
	s := sail.New(&sail.MetaConfig{
		ETCDEndpoints:  "x.x.x.x:3379",
		ProjectKey:     "xxxxxx",
		Namespace:      "user",
		Configs:        "user-api.yaml",
		ConfigFilePath: "./conf",
		LogLevel:       "DEBUG",
	}, sail.WithOnConfigChange(func(configFileKey string, s *sail.Sail) {
		if s.Err() != nil {
			fmt.Println(s.Err())
			return
		}
		fmt.Println(s.Pull())
		v, err := s.MergeVipers()
		if err != nil {
			fmt.Println(s.Err())
			return
		}
		v.Unmarshal(&cfg)
		fmt.Println(cfg, "\n", cfg.Datebase)
	}))
	if s.Err() != nil {
		fmt.Println(s.Err())
		return
	}
	fmt.Println(s.Pull())
	v, err := s.MergeVipers()
	if err != nil {
		fmt.Println(s.Err())
		return
	}
	v.Unmarshal(&cfg)
	fmt.Println(cfg, "\n", cfg.Datebase)

	for {
		time.Sleep(time.Second)
	}
}
