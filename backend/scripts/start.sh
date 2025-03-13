#!/bin/bash

# 设置环境
export APP_ENV=${APP_ENV:-dev}  # 默认为开发环境

# 启动应用
go run main.go 