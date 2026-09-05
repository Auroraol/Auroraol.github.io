# Redis 大 key 排查

## 一、快速定位：redis-cli 在线扫描

```bash
redis-cli --bigkeys
```

内部用 `SCAN` 游标遍历（不阻塞），按 key 类型统计，最后输出各类型最大的 key。

### 局限

1. **按元素数量判断，不是按内存**。string 按长度，但 list/set/hash/zset 报的是成员数，成员数少但单个成员巨大的情况会漏掉。
2. **抽样**，非全量精确扫描。

### 变体

```bash
redis-cli --memkeys          # 按实际内存占用找，比 --bigkeys 更接近真实
redis-cli --bigkeys -i 0.1   # 每个 SCAN 之间 sleep 0.1s，降低线上压力
```

按单个 key 精确查：

```bash
redis-cli MEMORY USAGE <key>
```

---

## 二、生产级方案：RDB 离线分析

量级大（上亿 key）或怕影响线上时，直接解析 dump 文件，不碰线上实例。

### 0. 拿到 RDB 文件

```bash
redis-cli CONFIG GET dir          # 目录，如 /var/lib/redis
redis-cli CONFIG GET dbfilename   # 文件名，如 dump.rdb

redis-cli BGSAVE                 # 手动触发一次持久化，确保 dump 最新
```

拿到完整路径，如 `/var/lib/redis/dump.rdb`。也可从备份机/从库拿，或复制到另一台机器解析。

### 方案 A：redis-rdb-cli（推荐，Go 单二进制，快）

```bash
wget https://github.com/leonchen83/redis-rdb-cli/releases/download/v0.9.1/redis-rdb-cli-release.zip
unzip redis-rdb-cli-release.zip
cd redis-rdb-cli-release        # 里面是 rct 可执行文件
```

**找出内存最大的前 10 个 key：**

```bash
./rct -f biggest -n 10 -s /var/lib/redis/dump.rdb
```

输出示例：

```
database,key,type,encoding,bytes,num_elements,len_largest_element,expiry
0,user:session:99123,hash,hashtable,52428800,1200000,96,
0,feed:list:888,list,quicklist,33554432,800000,64,
```

- `bytes` = 内存字节数
- `num_elements` = 成员数

**全量内存报表再排序：**

```bash
./rct -f memory -o mem.csv -s /var/lib/redis/dump.rdb
sort -t',' -k4 -nr mem.csv | head -20
```

### 方案 B：rdb-tools（Python 版）

```bash
pip install rdbtools python-lzf   # python-lzf 解析压缩数据，必装
```

```bash
# -l 100 只看最大的 100 个，-b 1024 过滤掉小于 1KB 的
rdb -c memory -l 100 -b 1024 /var/lib/redis/dump.rdb -f memory.csv

sort -t',' -k4 -nr memory.csv | head -20
```

CSV 列：`database,type,key,size_in_bytes,encoding,num_elements,len_largest_element,expiry`

---

## 三、选型对比

| | redis-rdb-cli | rdb-tools |
|---|---|---|
| 语言 | Go（单二进制） | Python（要装依赖） |
| 速度 | 快 | 慢（大文件差很多倍） |
| 找大 key | `-f biggest` 一条命令 | 生成 CSV 自己 sort |
| 适合 | 生产大文件 | 小文件 / 已装 Python |

**结论**：先用 `redis-cli --bigkeys` 跑一遍，九成场景够；要精确内存用 `--memkeys`；量级大或怕影响线上，用 `rct -f biggest -n 10 -s dump.rdb`。
