user-rpc-dev:
	@make -f deploy/mk/user-rpc.mk release-test

user-api-dev:
	@make -f deploy/mk/user-api.mk release-test

# social service
social-rpc-dev:
	@make -f deploy/mk/social-rpc.mk release-test

social-api-dev:
	@make -f deploy/mk/social-api.mk release-test

# im service
im-rpc-dev:
	@make -f deploy/mk/im-rpc.mk release-test

im-api-dev:
	@make -f deploy/mk/im-api.mk release-test

# im_websocket
im-ws-dev:
	@make -f deploy/mk/im-ws.mk release-test

# task_mq
task-mq-dev:
	@make -f deploy/mk/task-mq.mk release-test

release-test: user-rpc-dev user-api-dev social-rpc-dev social-api-dev im-rpc-dev im-api-dev im-ws-dev task-mq-dev

# eino_chat service
eino-rpc-gen:
	make -f deploy/mk/eino_chat.mk eino-rpc-gen
eino-api-gen:
	make -f deploy/mk/eino_chat.mk eino-api-gen
eino-api-dev:
	make -f deploy/mk/eino_chat.mk eino-api-dev
eino-knowledge-dev:
	make -f deploy/mk/eino_chat.mk eino-knowledge-dev
eino-docker-build:
	make -f deploy/mk/eino_chat.mk eino-docker-build
eino-docker-run:
	make -f deploy/mk/eino_chat.mk eino-docker-run

install-server:
	cd ./deploy/script && chmod +x release-test.sh && ./release-test.sh

# 本地构建并运行所有服务，不依赖阿里云仓库
# 使用方法: make local-run [LOCAL_MIRROR=true]
local-run:
	cd ./deploy/script && chmod +x local-run.sh && ./local-run.sh $(if $(LOCAL_MIRROR),--local-mirror)

# 停止所有服务
stop-all:
	@if command -v docker-compose &> /dev/null; then \
		docker-compose down; \
	else \
		docker compose down; \
	fi
	docker ps -aq -f "name=pai-pai-" | xargs -r docker stop
	docker ps -aq -f "name=pai-pai-" | xargs -r docker rm

# 查看所有镜像和服务状态
status-all:
	@cd ./deploy/script && chmod +x status-all.sh && bash ./status-all.sh
