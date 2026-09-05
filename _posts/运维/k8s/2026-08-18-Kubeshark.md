---
title: NodePort类型的Service
date: 2025-11-20 11:00:00 +0800
categories: [分类1, 分类1.1]
tags: [标签1, 标签2]
---

 ![image-20251120191713394](C:/Users/16658/AppData/Roaming/Typora/typora-user-images/image-20251120191713394.png)

集群内部（能执行 kubectl）    用节点内网 IP + NodePort：10.248.33.111:31072 

本地电脑（集群外）    用 kubectl port-forward 或 VPN

```
# 用节点 10.248.33.111 的 NodePort 31072 测试

psql -h 10.248.33.111 -p 31072 -U pguser -d postgres 快速实操（假设你在集群外部） bash 复制

# 在能执行 kubectl 的机器上开转发

kubectl -n dev-lane port-forward svc/postgresql-service 5432:5432

# 然后在本机连

psql -h 127.0.0.1 -p 5432 -U pguser -d postgres 一句话总结 ClusterIP 只能集群内用； NodePort 用节点 IP + 端口，但节点必须能被你访问到； 所有节点都没公网 IP，所以集群外必须做转发或 VPN。
```



| 名称                 | 类型     | 集群 IP        | 映射端口       | 年龄   |
| -------------------- | -------- | -------------- | -------------- | ------ |
| `postgresql`         | NodePort | 10.100.155.154 | 5432:32605/TCP | 2y358d |
| `postgresql-service` | NodePort | 10.107.191.220 | 5432:31072/TCP | 2y98d  |

![image-20251120192825576](C:/Users/16658/AppData/Roaming/Typora/typora-user-images/image-20251120192825576.png)

| 场景     | 连接方式示例                                                 |
| :------- | :----------------------------------------------------------- |
| 集群外部 | `psql -h <节点IP> -p 32605 -U postgres`                      |
| 集群内部 | `psql -h postgresql.default.svc.cluster.local -p 5432 -U postgres` |





## 驱逐

![image-20260326192818683](C:/Users/16658/AppData/Roaming/Typora/typora-user-images/image-20260326192818683.png)

[使用Kubeshark分析Kubernetes集群API流量-开发者社区-阿里云](https://developer.aliyun.com/article/1273795)





[终极指南：5分钟快速上手Kubeshark，轻松实现Kubernetes流量监控_gitblog_00081-开源鸿蒙跨平台开发者社区](https://openharmonycrossplatform.csdn.net/69dda88754b52172bc696b66.html)







```
# 设置版本号（最新稳定版，您也可以去 GitHub Releases 查看最新）
KUBESHARK_VERSION=$(curl -s https://api.github.com/repos/kubeshark/kubeshark/releases/latest | grep tag_name | cut -d '"' -f 4)

# 下载对应 Linux amd64 版本
curl -L -o kubeshark "https://github.com/kubeshark/kubeshark/releases/download/${KUBESHARK_VERSION}/kubeshark_linux_amd64"

# 赋予执行权限
chmod +x kubeshark

kubeshark version
```

### 步骤四：开始抓取业务请求 Header

由于您的服务器没有浏览器，您有两种方式获取 Header 数据：

#### 方式 A：将流量输出为 JSON 并查看 Header（纯终端）

bash

```
kubeshark tap -n default deployment/chatbot-client-common-deploy -o json | jq '.request.headers'
```



如果您没有 `jq`，可以直接查看原始输出并用 `grep` 过滤：

bash

```
kubeshark tap -n default deployment/chatbot-client-common-deploy -o json | grep -i '"headers"'
```



> **注意**：该命令会持续输出，您可以通过 `Ctrl+C` 终止。

#### 方式 B：录制流量到文件，然后离线查看

bash

```
# 录制 60 秒的数据到 dump.pcap
kubeshark tap -n default deployment/chatbot-client-common-deploy --dump -d 60 -o dump.pcap
```



录制完成后，您可以将 `dump.pcap` 文件拷贝到有图形界面的机器，用 Wireshark 或 Kubeshark 的 Web 界面打开查看（`kubeshark view dump.pcap`）。

#### 方式 C：通过端口转发在本地浏览器查看

如果您可以从本地电脑 SSH 到这台服务器，可以使用 SSH 端口转发：

1. 在服务器上启动 Kubeshark（但不打开浏览器）：

   bash

   ```
   kubeshark tap -n default deployment/chatbot-client-common-deploy
   ```

   

   它会默认监听 `localhost:8899`。

2. 在您的本地电脑（有浏览器）打开另一个终端，建立 SSH 隧道：

   bash

   ```
   ssh -L 8899:localhost:8899 worker@服务器IP
   ```

   

3. 然后在本地浏览器中访问 `http://localhost:8899`，即可看到 Kubeshark 的图形界面，轻松查看每个请求的 Header。





10.31.165.240





## 🔧 方案：使用 `tcpdump` + `kubectl` 抓取请求 Header

### 前提条件

- 您的业务 Pod 中运行的是 **HTTP（非 HTTPS）** 服务。如果是 HTTPS，则需要额外处理 TLS 解密，稍后补充。
- 您有权限执行 `kubectl exec` 或 `kubectl debug`（通常都具备）。

### 方法一：在 Pod 内直接抓包（最精确）

1. **获取目标 Pod 名称**：

   bash

   ```
   POD_NAME=$(kubectl get pods -n default -l app=chatbot-client-common -o jsonpath='{.items[0].metadata.name}')
   echo $POD_NAME
   ```

   

   如果您的 Deployment 标签不是 `app=chatbot-client-common`，请调整 `-l` 参数，或者直接用：

   bash

   ```
   kubectl get pods -n default | grep chatbot-client-common
   ```

   

2. **检查 Pod 内是否有 `tcpdump` 命令**：

   bash

   ```
   kubectl exec -n default $POD_NAME -- which tcpdump
   ```

   

   - 如果有输出路径（如 `/usr/sbin/tcpdump`），则直接执行下一步。
   - 如果没有，请使用 **方法二**（`kubectl debug` 临时容器）。

3. **开始抓取 HTTP 请求并过滤 Header**（假设服务监听 8080 端口）：

   bash

   ```
   kubectl exec -n default $POD_NAME -- tcpdump -i any -A -s 0 'tcp port 8080' | grep -E "^(GET|POST|PUT|DELETE|PATCH|HEAD|Host:|User-Agent:|Authorization:|X-|Cookie:)"
   ```

   

   - `-A`：以 ASCII 显示数据包内容。
   - `-s 0`：抓取完整包。
   - `'tcp port 8080'`：替换为您的服务实际监听端口（如 80、8080、3000 等）。
   - `grep` 过滤出 HTTP 方法和常见 Header 关键字。

   该命令会持续输出，直到您按 `Ctrl+C` 停止。您可以看到类似：

   text

   ```
   GET /api/users HTTP/1.1
   Host: example.com
   User-Agent: curl/7.68.0
   Authorization: Bearer eyJ...
   ```

   

### 方法二：使用 `kubectl debug` 添加临时抓包容器（当 Pod 内无 tcpdump 时）

Kubernetes 1.23+ 支持 `kubectl debug` 创建临时容器，挂载到目标 Pod 的网络命名空间，从而共享网络栈。

bash

```
kubectl debug -it -n default $POD_NAME --image=nicolaka/netshoot --target=<container-name> -- tcpdump -i any -A -s 0 'tcp port 8080' | grep -E "^(GET|POST|Host:|User-Agent:)"
```



- `--target=<container-name>` 指定要共享网络命名空间的目标容器名称（如果 Pod 只有一个容器，可以省略）。
- `nicolaka/netshoot` 是一个包含 `tcpdump`、`curl` 等网络工具的常用镜像。

> **注意**：`kubectl debug` 需要集群支持 `EphemeralContainers` 特性（默认在 K8s 1.23+ 启用）。如果不支持，可以考虑使用 `kubectl run` 临时 Pod，但会复杂一些。

### 方法三：从节点层面抓包（需要 root 权限）

如果您有节点 root 权限，可以 SSH 到节点，然后：

bash

```
# 获取 Pod 的 IP
POD_IP=$(kubectl get pod -n default $POD_NAME -o jsonpath='{.status.podIP}')
# 在节点上抓取该 IP 的流量
sudo tcpdump -i any -A -s 0 "host $POD_IP and tcp port 8080" | grep -E "^(GET|POST|Host:|User-Agent:)"
```



但您是 `worker` 用户，可能没有 `sudo` 权限，所以优先使用方法一或二。



```
# 抓取 60 秒，保存到 /tmp/capture.pcap
kubectl exec -n dev-lane chatbot-client-common-deploy-558f966d88-qsc6x -- \
  tcpdump -i any -s 0 -c 1000 -w /tmp/capture.pcap &
```

等待抓包完成后，将 pcap 文件复制到本地

```
kubectl cp -n dev-lane \
  chatbot-client-common-deploy-558f966d88-qsc6x:/tmp/capture.pcap \
  ./capture.pcap
```

### 用 `tshark` 查看 HTTP 请求头

现在文件已在本地，执行之前你尝试的命令即可（但之前用了不存在的字段，建议改用下面的方式）：

#### 完整查看 HTTP 协议树（推荐）

bash

```
tshark -r capture.pcap -Y "http.request" -O http
```



#### 或只提取关键字段（以表格形式输出）

bash

```
tshark -r capture.pcap -Y "http.request" -T fields \
  -e http.request.method \
  -e http.request.uri \
  -e http.host \
  -e http.user_agent \
  -e http.authorization \
  -e http.x_forwarded_for
```

### 每个请求的关键信息（以第一个为例）

- **请求方法**：`POST`
- **URI**：`/v1/conf_data/get`
- **HTTP 版本**：`HTTP/1.1`
- **目标主机**：`conf-engine-service:8080`
- **User-Agent**：`go-resty/1.12.0`（说明这个应用使用了 Go 的 Resty HTTP 客户端）
- **Content-Type**：`application/json; charset=utf-8`
- **请求体大小**：`102 字节`（具体 JSON 内容未在输出中显示，但存在）
- **自定义头部**（可能是你关心的）：
  - `Xdheader-Hostname: chatbot-client-common-deploy-558f966d88-qsc6x` —— 当前 Pod 名称
  - `Xdheader-Svcname: chatbot-client-common` —— 当前服务名
  - `Xd-Metadata:` （为空，可能是预留字段）
- **分布式追踪头部**（B3 协议）：
  - `X-B3-Traceid`、`X-B3-Spanid`、`X-B3-Sampled` —— 用于链路追踪（如 Zipkin/Jaeger）

------

### 删除容器内的临时文件以释放空间

kubectl exec -n dev-lane chatbot-client-common-deploy-558f966d88-qsc6x -- rm /tmp/capture.pcap





# tcpdump

## 进pod

```
tcpdump -i any -s 0 -w /tmp/capture.pcap
```

![image-20260818195448165](C:/Users/16658/AppData/Roaming/Typora/typora-user-images/image-20260818195448165.png)

只看 chatbot-api-gate 返回给快手的响应头

```bash
tcpdump -i eth0 -s 0 -A 'tcp port 8080' | grep -A 8 'chatbot-api-gate.*8080 > .*traefik.*HTTP/1\.1 200'
```

只看请求头（快手发过来的）

```bash
tcpdump -i eth0 -s 0 -A 'tcp port 8080' | grep -A 15 'POST /spi/chatbotevent'
```

同时看请求头 + 对应的响应头（推荐）

```bash
tcpdump -i eth0 -s 0 -A 'tcp port 8080' | grep -E -A 15 'POST /spi/chatbotevent|chatbot-api-gate.*8080 > .*traefik.*HTTP/1\.1'
```

示例

![image-20260818195554701](C:/Users/16658/AppData/Roaming/Typora/typora-user-images/image-20260818195554701.png)
