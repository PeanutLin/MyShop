#!/bin/bash

# 设置环境变量
export GOOS=linux
export GOARCH=amd64

# 编译应用程序
go build -o consumerapp consumer/consumer.go
go build -o fontedapp fonted/web/main.go
go build -o backendapp backend/web/main.go 
go build -o getproductapp getProduct/getProduct.go
go build -o validateapp validate/validate.go

# 创建目标文件夹
mkdir -p productshopapp

# 移动编译后的应用程序到目标文件夹
mv consumerapp productshopapp/
mv fontedapp productshopapp/
mv backendapp productshopapp/
mv getproductapp productshopapp/
mv validateapp productshopapp/

# 输出提示信息
echo "build finish"

