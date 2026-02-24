#!/bin/bash

# 简单的链接检查脚本
echo "检查内部链接..."

# 检查是否存在指向/tag/的链接
find _site -name "*.html" -exec grep -l "/tag/" {} \;

if [ $? -eq 0 ]; then
    echo "发现指向/tag/的链接"
    exit 1
else
    echo "未发现指向/tag/的链接 - 修复成功!"
fi

# 检查图片链接
echo "检查图片链接..."
find _site -name "*.html" -exec grep -l "image-20260214174940676.png" {} \;

if [ $? -eq 0 ]; then
    echo "发现损坏的图片链接"
    exit 1
else
    echo "未发现损坏的图片链接"
fi

echo "所有检查通过!"