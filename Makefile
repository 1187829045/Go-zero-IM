GO = /usr/local/go/bin/go

user-rpc-dev:
	@make -f deploy/mk/user-rpc.mk GO=$(GO) release-test

user-api-dev:
	@make -f deploy/mk/user-api.mk GO=$(GO) release-test

release-test: user-rpc-dev user-api-dev

install-server:
	cd ./deploy/script && chmod +x release-test.sh && ./release-test.sh
