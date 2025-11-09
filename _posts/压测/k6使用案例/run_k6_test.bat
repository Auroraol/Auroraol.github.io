@echo off
chcp 65001 > nul
REM K6 压测脚本快速启动脚本 (Windows)

REM 默认配置
set BASE_URL=%BASE_URL%
if "%BASE_URL%"=="" set BASE_URL=https://jiqiren.xiaoduoai.com/sdk/spi/test

set APP_KEY=%APP_KEY%
if "%APP_KEY%"=="" set APP_KEY=3409409348479354011

set TIMESTAMP=%TIMESTAMP%
if "%TIMESTAMP%"=="" set TIMESTAMP=2021-06-06 13:39:42

set SIGN=%SIGN%
if "%SIGN%"=="" set SIGN=8abb21bcfc4cc7ba4a501e2dc73a5e0c

set SCRIPT_FILE=%SCRIPT_FILE%
if "%SCRIPT_FILE%"=="" set SCRIPT_FILE=k6_benchmark.js

REM 检查 K6 是否安装
where k6 >nul 2>nul
if %errorlevel% neq 0 (
    echo 错误: 未找到 k6 命令
    echo 请先安装 K6: https://k6.io/docs/getting-started/installation/
    exit /b 1
)

REM 检查脚本文件是否存在
if not exist "%SCRIPT_FILE%" (
    echo 错误: 找不到脚本文件 %SCRIPT_FILE%
    exit /b 1
)

REM 检查数据文件是否存在
if not exist "new2new.jsonl" (
    echo 警告: 找不到数据文件 new2new.jsonl
    echo 请确保 new2new.jsonl 文件在脚本同目录下
    exit /b 1
)

REM 显示配置信息
echo ========== K6 压测配置 ==========
echo 目标URL: %BASE_URL%
echo APP_KEY: %APP_KEY%
echo TIMESTAMP: %TIMESTAMP%
echo SIGN: %SIGN%
echo 脚本文件: %SCRIPT_FILE%
echo ================================
echo.

REM 运行 K6 压测
echo 开始运行 K6 压测...
echo.

REM 启用 K6 Web Dashboard
set K6_WEB_DASHBOARD=true

set BASE_URL=%BASE_URL%
set APP_KEY=%APP_KEY%
set TIMESTAMP=%TIMESTAMP%
set SIGN=%SIGN%
k6 run %SCRIPT_FILE% %*

REM 检查退出码
if %errorlevel% equ 0 (
    echo.
    echo 压测完成！
) else (
    echo.
    echo 压测失败或未通过阈值检查
    exit /b 1
)

echo.
pause

