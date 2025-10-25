package main

import (
	"PaiPai/pkg/configserver"
	"flag"
	"fmt"
	"os"

	"PaiPai/apps/social/rpc/internal/config"
	"PaiPai/apps/social/rpc/internal/server"
	"PaiPai/apps/social/rpc/internal/svc"
	"PaiPai/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/dev/social.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	//conf.MustLoad(*configFile, &c)
	var configs = "social.yaml"
	// 使用etcd容器名称而不是etcd，以便在Docker网络中直接解析
	hostIP := "etcd"
	if envHostIP := os.Getenv("HOST_IP"); envHostIP != "" {
		hostIP = envHostIP
	}
	err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
		// 修改为字符串数组格式，以匹配Config结构体中的类型定义
		ETCDEndpoints:  []string{fmt.Sprintf("%s:2379", hostIP)},
		ProjectKey:     "paipai",
		Namespace:      "social",
		Configs:        configs,
		// ConfigFilePath应该是目录路径而非文件路径
		ConfigFilePath: "/social/conf",
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

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		social.RegisterSocialServer(grpcServer, server.NewSocialServer(svc.NewServiceContext(c)))

		// err := grpcServer.AddUnaryInterceptors(
		// 拦截器
		// )
		// if err != nil {
		//  panic(err)
		// }

		// 添加反射支持，便于调试
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
