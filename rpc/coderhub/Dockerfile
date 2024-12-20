# 第 1 阶段：构建阶段
FROM golang:1.22-alpine3.19 AS builder

# 设置 Go 模块代理
ENV GOPROXY=https://goproxy.cn,direct

# 设置工作目录
WORKDIR /app

# 将代码复制到工作目录中
COPY . .

# 安装依赖
RUN go mod download

# 编译 Go 程序
RUN go build -o /app/bin/cmd ./rpc/coderhub/coderhub.go

# 第 2 阶段：运行阶段（轻量级）
FROM alpine:3.19

# 设置工作目录
WORKDIR /app

# 复制编译后的二进制文件和配置文件到最终镜像中
COPY --from=builder /app/bin/cmd /app/bin/cmd
COPY --from=builder /app/rpc/coderhub/etc /app/bin/etc

# 暴露服务的端口
EXPOSE 8888

# 定义容器启动时运行的命令
CMD ["/app/bin/cmd", "-f", "/app/bin/etc/coderhub.yaml"]