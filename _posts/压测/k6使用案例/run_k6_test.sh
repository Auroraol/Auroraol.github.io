#!/bin/bash

# K6 压测脚本快速启动脚本

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 默认配置
BASE_URL="${BASE_URL:-http://localhost:8092/sdk/spi/test}"
APP_KEY="${APP_KEY:-3409409348479354011}"
TIMESTAMP="${TIMESTAMP:-2021-06-06 13:39:42}"
SIGN="${SIGN:-8abb21bcfc4cc7ba4a501e2dc73a5e0c}"
SCRIPT_FILE="${SCRIPT_FILE:-k6_benchmark.js}"

# 检查 K6 是否安装
if ! command -v k6 &> /dev/null; then
    echo -e "${RED}错误: 未找到 k6 命令${NC}"
    echo "请先安装 K6: https://k6.io/docs/getting-started/installation/"
    exit 1
fi

# 检查脚本文件是否存在
if [ ! -f "$SCRIPT_FILE" ]; then
    echo -e "${RED}错误: 找不到脚本文件 $SCRIPT_FILE${NC}"
    exit 1
fi

# 检查数据文件是否存在
if [ ! -f "new2new.jsonl" ]; then
    echo -e "${YELLOW}警告: 找不到数据文件 new2new.jsonl${NC}"
    echo "请确保 new2new.jsonl 文件在脚本同目录下"
    exit 1
fi

# 显示配置信息
echo -e "${GREEN}========== K6 压测配置 ==========${NC}"
echo "目标URL: $BASE_URL"
echo "APP_KEY: $APP_KEY"
echo "TIMESTAMP: $TIMESTAMP"
echo "SIGN: $SIGN"
echo "脚本文件: $SCRIPT_FILE"
echo -e "${GREEN}================================${NC}"
echo ""

# 运行 K6 压测
# 启用 K6 Web Dashboard
export K6_WEB_DASHBOARD=true

echo -e "${GREEN}开始运行 K6 压测...${NC}"
echo ""

BASE_URL="$BASE_URL" \
APP_KEY="$APP_KEY" \
TIMESTAMP="$TIMESTAMP" \
SIGN="$SIGN" \
k6 run "$SCRIPT_FILE" "$@"

# 检查退出码
if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}压测完成！${NC}"
else
    echo ""
    echo -e "${RED}压测失败或未通过阈值检查${NC}"
    exit 1
fi

