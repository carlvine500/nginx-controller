all: build
build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build && mv nginx-controller nginx-controller.amd64; \
    GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build && mv nginx-controller nginx-controller.arm64
clean:
	rm -rf ./clickpaas-login.amd64 ./clickpaas-login.arm64;

