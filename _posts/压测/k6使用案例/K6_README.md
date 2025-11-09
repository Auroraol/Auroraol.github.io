# K6 压测脚本使用说明

## 文件说明

- `k6_benchmark.js`: 主要的 K6 压测脚本
- `new2new.jsonl`: 测试数据文件（JSONL 格式，每行一个 JSON 对象）

## 安装 K6

### Windows
```powershell
# 使用 Chocolatey
choco install k6

# 或使用 Scoop
scoop install k6

# 或下载安装包
# https://github.com/grafana/k6/releases
```

### macOS
```bash
brew install k6
```

### Linux
```bash
# Debian/Ubuntu
sudo gpg -k
sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
sudo apt-get update
sudo apt-get install k6
```

## 使用方法

### 1. 基本使用

```bash
# 使用默认配置运行
k6 run k6_benchmark.js
```

### 2. 自定义配置

#### 通过环境变量配置

```bash
# Windows (PowerShell)
$env:BASE_URL="http://localhost:8092/sdk/spi/test"
$env:APP_KEY="3409409348479354011"
$env:TIMESTAMP="2021-06-06 13:39:42"
$env:SIGN="8abb21bcfc4cc7ba4a501e2dc73a5e0c"
k6 run k6_benchmark.js

# Linux/macOS
BASE_URL="http://localhost:8092/sdk/spi/test" \
APP_KEY="3409409348479354011" \
TIMESTAMP="2021-06-06 13:39:42" \
SIGN="8abb21bcfc4cc7ba4a501e2dc73a5e0c" \
k6 run k6_benchmark.js
```

#### 修改压测配置

编辑 `k6_benchmark.js` 中的 `options` 对象：

```javascript
export const options = {
    stages: [
        { duration: '10s', target: 5 },    // 10秒内增加到5个虚拟用户
        { duration: '30s', target: 10 },   // 30秒内增加到10个虚拟用户
        { duration: '1m', target: 20 },    // 1分钟内增加到20个虚拟用户
        { duration: '30s', target: 10 },   // 30秒内减少到10个虚拟用户
        { duration: '10s', target: 0 },    // 10秒内减少到0个虚拟用户
    ],
    // ...
};
```

### 3. 高级用法

#### 固定虚拟用户数运行指定时间

```javascript
export const options = {
    vus: 20,           // 虚拟用户数
    duration: '5m',    // 运行时间
};
```

#### 固定 QPS 压测

```javascript
export const options = {
    scenarios: {
        constant_request_rate: {
            executor: 'constant-arrival-rate',
            rate: 20,              // 每秒20个请求
            timeUnit: '1s',
            duration: '5m',
            preAllocatedVUs: 10,
            maxVUs: 50,
        },
    },
};
```

#### 输出结果到文件

```bash
# 输出到 JSON 文件
k6 run --out json=results.json k6_benchmark.js

# 输出到 CSV 文件
k6 run --out csv=results.csv k6_benchmark.js

# 输出到 InfluxDB
k6 run --out influxdb=http://localhost:8086/k6 k6_benchmark.js
```

## 配置说明

### 环境变量

- `BASE_URL`: 接口基础 URL（默认: `http://localhost:8092/sdk/spi/test`）
- `APP_KEY`: app_key 参数（默认: `3409409348479354011`）
- `TIMESTAMP`: timestamp 参数（默认: `2021-06-06 13:39:42`）
- `SIGN`: sign 参数（默认: `8abb21bcfc4cc7ba4a501e2dc73a5e0c`）

### 压测阶段说明

K6 使用 stages 来定义压测的不同阶段：

```javascript
stages: [
    { duration: '10s', target: 5 },   // 阶段1: 10秒内从0增加到5个虚拟用户
    { duration: '30s', target: 10 },  // 阶段2: 30秒内从5增加到10个虚拟用户
    { duration: '1m', target: 20 },   // 阶段3: 1分钟内从10增加到20个虚拟用户
    { duration: '30s', target: 10 },  // 阶段4: 30秒内从20减少到10个虚拟用户
    { duration: '10s', target: 0 },   // 阶段5: 10秒内从10减少到0个虚拟用户
]
```

### 阈值设置

```javascript
thresholds: {
    'http_req_duration': ['p(95)<3000', 'p(99)<5000'], // 95%的请求应该在3秒内完成，99%在5秒内
    'errors': ['rate<0.1'],                            // 错误率应该小于10%
    'http_req_failed': ['rate<0.1'],                   // HTTP 请求失败率应该小于10%
}
```

如果阈值不满足，K6 会以非零退出码退出。

## 输出说明

K6 会输出以下指标：

- `http_reqs`: HTTP 请求总数
- `http_req_duration`: 请求响应时间（平均值、最小值、最大值、P95、P99等）
- `http_req_failed`: 请求失败率
- `errors`: 错误率
- `vus`: 虚拟用户数
- `iteration_duration`: 迭代持续时间

## 示例输出

```
     ✓ 状态码是 200
     ✓ 响应时间 < 5000ms
     ✓ 响应体不为空

     checks.........................: 100.00% ✓ 1350  ✗ 0
     data_received..................: 2.1 MB  35 kB/s
     data_sent......................: 15 MB   250 kB/s
     http_req_duration..............: avg=234.5ms min=45.2ms med=198.3ms max=1234.5ms p(90)=456.7ms p(95)=678.9ms
     http_req_failed................: 0.00%   ✓ 0    ✗ 450
     http_reqs......................: 450     7.5/s
     iteration_duration.............: avg=334.5ms min=145.2ms med=298.3ms max=1334.5ms p(90)=556.7ms p(95)=778.9ms
     iterations.....................: 450     7.5/s
     vus............................: 20      min=1  max=20
     vus_max........................: 20      min=1  max=20
```

## 注意事项

1. 确保 `new2new.jsonl` 文件在脚本同目录下
2. 确保 JSONL 文件的每行都是有效的 JSON 格式
3. 如果数据量很大，考虑使用 SharedArray 来优化内存使用
4. 根据实际服务器性能调整虚拟用户数和请求速率
5. 建议先在测试环境进行压测，避免影响生产环境

## 故障排查

### 问题1: 无法读取文件

**错误**: `无法读取文件 new2new.jsonl: ...`

**解决**: 
- 确保 `new2new.jsonl` 文件在脚本同目录
- 检查文件路径是否正确
- 检查文件权限

### 问题2: JSON 解析失败

**错误**: `解析第 X 行失败: ...`

**解决**:
- 检查 JSONL 文件格式是否正确
- 确保每行都是有效的 JSON
- 检查是否有特殊字符需要转义

### 问题3: 连接失败

**错误**: `请求失败: status=0`

**解决**:
- 检查目标服务器是否运行
- 检查 URL 是否正确
- 检查网络连接
- 检查防火墙设置

## 参考资源

- [K6 官方文档](https://k6.io/docs/)
- [K6 JavaScript API](https://k6.io/docs/javascript-api/)
- [K6 压测场景](https://k6.io/docs/using-k6/scenarios/)

