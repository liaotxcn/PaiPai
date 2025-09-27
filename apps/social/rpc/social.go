package main

import (
	"PaiPai/pkg/configserver"
	"PaiPai/pkg/interceptor"
	"PaiPai/pkg/interceptor/rpcserver"
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

var configFile = flag.String("f", "etc/social.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	var configs = "social.yaml"
	// 从环境变量获取HOST_IP，如果没有则使用默认值
	hostIP := "host.docker.internal"
	if envHostIP := os.Getenv("HOST_IP"); envHostIP != "" {
		hostIP = envHostIP
	}
	err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
		ETCDEndpoints:  fmt.Sprintf("%s:3379", hostIP),
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
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		social.RegisterSocialServer(grpcServer, server.NewSocialServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	s.AddUnaryInterceptors(rpcserver.LogInterceptor, rpcserver.SyncLimiterInterceptor(10))
	s.AddUnaryInterceptors(interceptor.NewIdempotenceServer(interceptor.NewDefaultIdempotent(c.Cache[0].RedisConf)))

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
