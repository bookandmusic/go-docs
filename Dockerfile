# 第一阶段：编译Go项目和压缩静态文件和模板
FROM golang:1.21 AS builder

WORKDIR /app

# 设置Go代理
ENV GOPROXY=https://goproxy.cn,direct

# 将项目代码复制到容器中
COPY . .

# 安装minify工具
RUN go install github.com/tdewolff/minify/v2/cmd/minify@latest

# 使用 minify 压缩静态文件
RUN minify -r -o . ./static

# 整理并删除不再使用的依赖项
RUN go mod tidy

# 编译Go项目为二进制文件
RUN GOOS=linux CGO_ENABLED=1 GOARCH=amd64 go build -ldflags="-s -w" -installsuffix cgo -o godocs 

# 第二阶段：运行时环境
FROM ubuntu:latest

WORKDIR /app

# 从第一阶段复制编译好的二进制文件
COPY --from=builder /app/godocs .

# 复制静态文件和模板
COPY --from=builder /app/static ./static
COPY --from=builder /app/templates ./templates
# 复制配置文件
COPY --from=builder /app/config/config.ini ./config/config.ini
# 复制gojieba字典文件
COPY --from=builder /app/lib/gojieba/dict ./lib/gojieba/dict

CMD ["-c", "/app/config/config.ini"]
ENTRYPOINT ["/app/godocs", "server"]
