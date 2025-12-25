FROM registry.cn-zhangjiakou.aliyuncs.com/xiaoduoai/alpine:3.15.0
WORKDIR /app
COPY dist /app
CMD ["./goods_center", "-c", "/etc/config/goods_center.toml"]