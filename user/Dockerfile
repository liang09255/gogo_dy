FROM golang:1.20-alpine
WORKDIR /app
COPY ../../user .
RUN go env -w GOPROXY=https://goproxy.cn
RUN go mod tidy
RUN go build .
EXPOSE 8081 8881
# 运行Go应用程序
CMD ["./main"]
