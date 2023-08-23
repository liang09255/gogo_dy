# 第一阶段-把所有文件放进基础镜像
FROM golang:1.20-alpine

WORKDIR /app

COPY . .