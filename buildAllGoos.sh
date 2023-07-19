#!/bin/bash

appName="test";

# 创建文件夹
mkdir -p ./allGoosBuild


# 构建 Windows 可执行文件
echo "正在打包Windows";
env GOOS=windows GOARCH=amd64 go build -o ./allGoosBuild/${appName}-windows.exe

# 构建 Mac 可执行文件
echo "正在打包MAC";
env GOOS=darwin GOARCH=amd64 go build -o ./allGoosBuild/${appName}-mac

# 构建 Linux 可执行文件
echo "正在打包Linux";
env GOOS=linux GOARCH=amd64 go build -o ./allGoosBuild/${appName}-linux