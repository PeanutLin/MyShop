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
mkdir -p productshop

# 移动编译后的应用程序到目标文件夹
mv consumerapp productshop/
mv fontedapp productshop/
mv backendapp productshop/
mv getproductapp productshop/
mv validateapp productshop/

# 输出提示信息
echo "build finish"

