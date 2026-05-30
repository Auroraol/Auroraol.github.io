---
title: traefik ingress
date: 2026-05-25 11:00:00 +0800
categories: [k8s]
tags: [traefik ]
---

# requirements.txt

```
# 删除旧的虚拟环境
rm -rf .venv

# 创建新的虚拟环境
python -m venv .venv

# 激活虚拟环境
.venv\Scripts\activate

# 安装所有依赖项
pip install -r requirements.txt

```

进入虚拟环境

![image-20250116192848360](D:\Github\python-note\Python2.0.assets\image-20250116192848360.png)





# 安装 pipreqs
pip install pipreqs

# 扫描项目并生成 requirements.txt
pipreqs . --force

# 或者只更新现有文件
pipreqs . --savepath requirements.txt
