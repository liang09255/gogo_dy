ARG GOGO_VERSION
# 第二阶段-构建可执行文件
FROM gogo_base:$GOGO_VERSION AS base
WORKDIR /app/video
RUN go env -w GOPROXY=https://goproxy.cn
RUN go build main.go

# 第三阶段-构建最小的可执行镜像
FROM debian:buster-slim
COPY --from=base /app /app
WORKDIR /app/video
EXPOSE 8085 8085
CMD ["./main"]
