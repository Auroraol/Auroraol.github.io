# Redis MONITOR 日志使用

## 一、`mon.tar.gz` 是什么

`mon` = `redis-cli monitor`。用 MONITOR 命令把 Redis 一段时间内收到的**所有命令流**实时抓下来，打包成 tar.gz。本目录的 `mon.tar.gz` 解压后是一个约 20MB 的纯文本文件（名叫 `23`），每行一条命令。

## 二、怎么抓取

```bash
redis-cli -h <host> -p <port> -a <password> monitor > mon.log
```

按条数截断（避免文件无限增长）：

```bash
redis-cli monitor | head -n 1000000 > mon.log
```

打包：

```bash
tar -czf mon.tar.gz mon.log
```

![image-20260819175312717](C:/Users/16658/AppData/Roaming/Typora/typora-user-images/image-20260819175312717.png)

## 三、日志格式

```
1786704776.711720 [0 10.31.164.232:60350] "set" "key" "value" "ex" "3" "nx"
```

| 字段 | 含义 |
|------|------|
| `1786704776.711720` | Unix 时间戳（秒.微秒） |
| `0` | 数据库编号（db0） |
| `10.31.164.232:60350` | 客户端 IP:端口 |
| `"set"` | 执行的命令 |
| 后续字符串 | key / value / 参数（`evalsha` 是 Lua 脚本执行） |

## 四、怎么用这份 dump

解压：

```bash
tar -xzf mon.tar.gz        # 得到文件 23
```

### 1. 看某类命令

```bash
grep '"set"' 23 | head
grep '"evalsha"' 23 | wc -l     # Lua 脚本执行次数
```

### 2. 统计各命令调用次数

```bash
awk -F'"' '{print $4}' 23 | sort | uniq -c | sort -rn | head -20
```

### 3. 找出访问最频繁的 key（定位热点 key）

```bash
awk -F'"' '{print $6}' 23 | sort | uniq -c | sort -rn | head -20
```

### 4. 找出某个 key 的全部操作

```bash
grep 'concurrent_controller_report_client' 23
```

### 5. 找出大 value 的写入

```bash
awk -F'"' '{print length($6), $4, $6}' 23 | sort -rn | head -20
```

## 五、注意事项

1. **MONITOR 会拖慢 Redis 性能**：所有命令都要经过监控缓冲区，高并发线上慎用，抓完立刻 `Ctrl+C`。
2. **日志含明文 key/value**，属敏感数据，别外传。
3. **命令量巨大**：生产环境一秒几万条，务必用 `head -n` 截断或定时轮转，否则磁盘会被打满。