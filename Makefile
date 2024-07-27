GO = /usr/local/go/bin/go

user-rpc-dev:
	@make -f deploy/mk/user-rpc.mk GO=$(GO) release-test

user-api-dev:
	@make -f deploy/mk/user-api.mk GO=$(GO) release-test


social-rpc-dev:
	@make -f deploy/mk/social-rpc.mk GO=$(GO) release-test

social-api-dev:
	@make -f deploy/mk/social-api.mk GO=$(GO) release-test

im-ws-dev:
	@make -f deploy/mk/im-ws.mk GO=$(GO) release-test

im-rpc-dev:
	@make -f deploy/mk/im-rpc.mk GO=$(GO) release-test

im-api-dev:
	@make -f deploy/mk/im-api.mk GO=$(GO) release-test

task-mq-dev:
	@make -f deploy/mk/task-mq.mk GO=$(GO) release-test

release-test: user-rpc-dev user-api-dev social-rpc-dev social-api-dev im-ws-dev im-rpc-dev im-api-dev task-mq-dev

install-server:
	cd ./deploy/script && chmod +x release-test.sh && ./release-test.sh

install-server-user-rpc:
	cd ./deploy/script && chmod +x user-rpc-test.sh && ./release-test.sh

