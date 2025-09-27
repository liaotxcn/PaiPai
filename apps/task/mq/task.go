package mq

import (
	"PaiPai/apps/task/mq/internal/config"
	"PaiPai/apps/task/mq/internal/handler"
	"PaiPai/apps/task/mq/internal/svc"
	"PaiPai/pkg/configserver"
	"flag"
	"fmt"
	"os"
	"github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "etc/dev/task.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	var configs = "task.yaml"
	// 从环境变量获取HOST_IP，如果没有则使用默认值
	hostIP := "host.docker.internal"
	if envHostIP := os.Getenv("HOST_IP"); envHostIP != "" {
		hostIP = envHostIP
	}
	err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
		ETCDEndpoints:  fmt.Sprintf("%s:3379", hostIP),
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

	if err := c.SetUp(); err != nil {
		panic(err)
	}
	ctx := svc.NewServiceContext(c)
	listen := handler.NewListen(ctx)

	serviceGroup := service.NewServiceGroup()
	for _, s := range listen.Services() {
		serviceGroup.Add(s)
	}
	fmt.Println("starting service at ...", c.ListenOn)
	serviceGroup.Start()
}
