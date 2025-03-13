@echo off
chcp 65001
set "REDIS_PATH=D:\project\Redis-x64-3.2.100"

:: 检查路径是否存在
if not exist "%REDIS_PATH%" (
    echo Redis路径不存在: %REDIS_PATH%
    pause
    exit /b 1
)

:: 检查文件是否存在
if not exist "%REDIS_PATH%\redis-server.exe" (
    echo Redis服务器程序不存在
    pause
    exit /b 1
)

if not exist "%REDIS_PATH%\redis-cli.exe" (
    echo Redis客户端程序不存在
    pause
    exit /b 1
)

echo 正在启动Redis服务器...
cd /d "%REDIS_PATH%"
start "Redis Server" cmd /k redis-server.exe
timeout /t 2
echo Redis服务器已启动

echo 正在启动Redis客户端...
start "Redis CLI" cmd /k redis-cli.exe
echo Redis客户端已启动

echo Redis环境启动完成!
pause 