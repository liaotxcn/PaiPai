package main

import (
	"PaiPai/pkg/configserver"
	"flag"
	"fmt"
	"os"

	"PaiPai/apps/task/mq/internal/config"
	"PaiPai/apps/task/mq/internal/handler"
	"PaiPai/apps/task/mq/internal/svc"

	"github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "etc/dev/task.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	//conf.MustLoad(*configFile, &c)
	var configs = "task.yaml"
	// 使用etcd容器名称而不是etcd，以便在Docker网络中直接解析
	hostIP := "etcd"
	if envHostIP := os.Getenv("HOST_IP"); envHostIP != "" {
		hostIP = envHostIP
	}
	err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
		// 修改为字符串数组格式，以匹配Config结构体中的类型定义
		ETCDEndpoints:  []string{fmt.Sprintf("%s:2379", hostIP)},
		ProjectKey:     "paipai",
		Namespace:      "task",
		Configs:        configs,
		// ConfigFilePath应该是目录路径而非文件路径
		ConfigFilePath: "/task/conf",
		LogLevel: "DEBUG",
	})).MustLoad(&c, func(bytes []byte) error {
		var c config.Config
		err := configserver.LoadFromJsonBytes(bytes, &c)
		if err != nil {
			fmt.Println("config read err :", err)
			return nil
		}
		fmt.Printf("%s config has changed :%+v \n", configs, c)
		return nil
	})
	if err != nil {
		panic(err)
	}

	svcCtx := svc.NewServiceContext(c)
	listen := handler.NewListen(svcCtx)
	services := listen.Services()

	// 启动所有服务
	serviceGroup := service.NewServiceGroup()
	for _, s := range services {
		serviceGroup.Add(s)
	}
	defer serviceGroup.Stop()

	fmt.Printf("Starting server...\n")
	serviceGroup.Start()
}
