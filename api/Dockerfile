ARG GOGO_VERSION
# 第二阶段-构建可执行文件
FROM gogo_base:$GOGO_VERSION AS base
WORKDIR /app/api
RUN go env -w GOPROXY=https://goproxy.cn
RUN go build main.go

# 第三阶段-构建最小的可执行镜像
FROM debian:buster-slim
COPY --from=base /app /app
WORKDIR /app/api
EXPOSE 8080 8080
CMD ["./main"]