# user service
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
