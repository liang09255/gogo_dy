FROM golang:1.20-alpine
WORKDIR /app
COPY . .
RUN go env -w GOPROXY=https://goproxy.cn
RUN go mod tidy
RUN go build .
# 运行Go应用程序
CMD ["./main"]
