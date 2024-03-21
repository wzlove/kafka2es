# 第一层是构建一个builder镜像,目的是在其中编译出可执行文件 kafka2es
FROM golang:alpine AS builder

LABEL stage=gobuilder

# 禁用cgo
ENV CGO_ENABLED 0
# 添加代理
ENV GOPROXY https://goproxy.cn,direct
RUN apk update --no-cache && apk add --no-cache tzdata


WORKDIR /app

ADD go.mod .
ADD go.sum .
RUN go mod tidy
COPY .. .
COPY ../src/etc/ /app/etc
RUN go build -ldflags="-s -w" -o /app/kafka2es src/main.go


# 从第一个镜像里copy出来可执行文件,并且使用尽可能小的基础镜像 alphine 以保障最终镜像尽可能小
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/kafka2es /app/kafka2es
COPY --from=builder /app/etc /app/etc

EXPOSE "8080"

CMD ["./kafka2es","--config","etc/config.yaml"]
