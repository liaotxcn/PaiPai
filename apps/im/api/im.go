package main

import (
	"PaiPai/pkg/configserver"
	"PaiPai/pkg/resultx"
	"flag"
	"fmt"
	"os"
	"github.com/zeromicro/go-zero/rest/httpx"

	"PaiPai/apps/im/api/internal/config"
	"PaiPai/apps/im/api/internal/handler"
	"PaiPai/apps/im/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/dev/im.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	//conf.MustLoad(*configFile, &c)
	var configs = "im.yaml"
	// 从环境变量获取HOST_IP，如果没有则使用默认值
	hostIP := "host.docker.internal"
	if envHostIP := os.Getenv("HOST_IP"); envHostIP != "" {
		hostIP = envHostIP
	}
	err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
		ETCDEndpoints:  fmt.Sprintf("%s:3379", hostIP),
		ProjectKey:     "paipai",
		Namespace:      "im",
		Configs:        configs,
		// ConfigFilePath应该是目录路径而非文件路径
		ConfigFilePath: "/im/conf",
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

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	httpx.SetErrorHandlerCtx(resultx.ErrHandler(c.Name))
	httpx.SetOkHandler(resultx.OkHandler)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
