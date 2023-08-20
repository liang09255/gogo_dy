ARG GOGO_VERSION
# 第二阶段-构建可执行文件
FROM gogo_base:$GOGO_VERSION AS base
WORKDIR /app/user
RUN go env -w GOPROXY=https://goproxy.cn
RUN go build main.go

# 第三阶段-构建最小的可执行镜像
FROM debian:buster-slim
COPY --from=base /app /app
WORKDIR /app/user
EXPOSE 8081 8081
CMD ["./main"]
