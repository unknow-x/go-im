.PHONY:all build run test cover docker help

BINARY_IMD=imd
all: build run

build:
	@go build -o $(BINARY_IMD) cmd/im/main.go

run:
	./"${BINARY_IMD}"

clean:
	@go clean
	rm --force "cover.out"

test:
	@GO_IM_ENV_PATH=".." go test -v test/*

cover:
	@go test -coverprofile cover.out
	@go tool cover -html=cover.out

docker:
	@docker build -t go-im:latest ./build

help:
	@echo "make 格式化go代码 并编译生成二进制文件"
	@echo "make build 编译go代码生成二进制文件"
	@echo "make clean 清理中间目标文件"
	@echo "make test 执行测试case"
	@echo "make cover 检查测试覆盖率"
	@echo "make run 直接运行程序"
	@echo "make lint 执行代码检查"
	@echo "make docker 构建docker镜像"