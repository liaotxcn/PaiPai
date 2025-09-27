package main

import (
	"PaiPai/pkg/interceptor/rpcserver"
	"flag"
	"fmt"
	"os"
	"strings"

	"PaiPai/apps/user/rpc/internal/config"
	"PaiPai/apps/user/rpc/internal/server"
	"PaiPai/apps/user/rpc/internal/svc"
	"PaiPai/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config

	// 从环境变量获取HOST_IP，如果没有则使用默认值
	hostIP := "host.docker.internal"
	if envHostIP := os.Getenv("HOST_IP"); envHostIP != "" {
		hostIP = envHostIP
	}

	// 创建一个临时配置文件，替换其中的环境变量
	content, err := os.ReadFile(*configFile)
	if err != nil {
		// 如果无法读取配置文件，尝试从容器内路径读取
		content, err = os.ReadFile("/user/conf/user.yaml")
		if err != nil {
			panic(err)
		}
	}

	// 替换配置文件中的${HOST_IP}占位符
	content = []byte(strings.Replace(string(content), "${HOST_IP}", hostIP, -1))

	// 写入临时文件
	tmpFile := "/tmp/user.yaml"
	err = os.WriteFile(tmpFile, content, 0644)
	if err != nil {
		panic(err)
	}

	// 使用go-zero的默认方式加载处理过的配置文件
	fmt.Println("使用处理后的配置文件:", tmpFile)
	conf.MustLoad(tmpFile, &c)

	ctx := svc.NewServiceContext(c)

	if err := ctx.SetRootToken(); err != nil {
		panic(err)
	}

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		user.RegisterUserServer(grpcServer, server.NewUserServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	s.AddUnaryInterceptors(rpcserver.LogInterceptor)

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
