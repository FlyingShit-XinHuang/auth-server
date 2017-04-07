PKG=whispir/auth-server
BUILD=docker run -ti --rm -v `pwd`:/go/src/$(PKG) -w /go/src/$(PKG) iron/go:1.7-dev go build

auth-cmd:
	$(BUILD) -o auth $(PKG)/examples/cmd/auth

manage-cmd:
	$(BUILD) -o manage $(PKG)/examples/cmd/resources

docker-demo-auth-cmd: auth-cmd
	docker build -f Dockerfile.demo-auth-cmd -t demo-auth .

docker-demo-manage-cmd: manage-cmd
	docker build -f Dockerfile.demo-manage-cmd -t demo-manage .

docker-build-demo: docker-demo-auth-cmd docker-demo-manage-cmd

server:
	$(BUILD) -o auth-server

docker-build: server
	docker build -t auth-server .

.PHONY: server manage-cmd auth-cmd docker-demo-auth-cmd docker-demo-manage-cmd
