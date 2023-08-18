# 第二阶段-构建可执行文件
FROM gogo_dy_base AS base

WORKDIR /app/gogo_dy_api

RUN go env -w GOPROXY=https://goproxy.cn

RUN go build main.go

EXPOSE 8080
# 运行Go应用程序
CMD ["./main"]


# 第三阶段-构建最小的可执行镜像
FROM debian:buster-slim
WORKDIR /app/gogo_dy_api
COPY --from=base /app/gogo_dy_api .
CMD ["./main"]