# eino_chat service Makefile

# 生成rpc代码
eino-rpc-gen:
	goctl rpc protoc apps/eino_chat/rpc/eino.proto --go_out=./apps/eino_chat/rpc --go-grpc_out=./apps/eino_chat/rpc --zrpc_out=./apps/eino_chat/rpc

# 生成api代码
eino-api-gen:
	goctl api go -api apps/eino_chat/api/eino.api -dir apps/eino_chat/api

# 启动eino_chat api服务
eino-api-dev:
	source apps/eino_chat/.env && go run apps/eino_chat/api/eino.go -f apps/eino_chat/api/etc/eino-api.yaml

# 启动知识库服务
eino-knowledge-dev:
	source apps/eino_chat/.env && go run apps/eino_chat/cmd/knowledge_base/main.go

# 构建eino_chat docker镜像
eino-docker-build:
	docker build -t paipai/eino_chat -f deploy/dockerfile/dockerfile_eino_chat .

# 运行eino_chat docker容器
eino-docker-run:
	docker run -d --name eino_chat -p 8888:8888 --env-file apps/eino_chat/.env paipai/eino_chat