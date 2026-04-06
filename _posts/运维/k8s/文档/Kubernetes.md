# 基本概念

<svg viewBox="0 0 1400 1100" xmlns="http://www.w3.org/2000/svg">   <defs>     <!-- 渐变定义 -->     <linearGradient id="clusterGrad" x1="0%" y1="0%" x2="0%" y2="100%">       <stop offset="0%" style="stop-color:#1a237e;stop-opacity:1" />       <stop offset="100%" style="stop-color:#0d1642;stop-opacity:1" />     </linearGradient>     <linearGradient id="masterGrad" x1="0%" y1="0%" x2="0%" y2="100%">       <stop offset="0%" style="stop-color:#c62828;stop-opacity:1" />       <stop offset="100%" style="stop-color:#8e0000;stop-opacity:1" />     </linearGradient>     <linearGradient id="workerGrad" x1="0%" y1="0%" x2="0%" y2="100%">       <stop offset="0%" style="stop-color:#1565c0;stop-opacity:1" />       <stop offset="100%" style="stop-color:#0d47a1;stop-opacity:1" />     </linearGradient>     <linearGradient id="podGrad" x1="0%" y1="0%" x2="0%" y2="100%">       <stop offset="0%" style="stop-color:#2e7d32;stop-opacity:1" />       <stop offset="100%" style="stop-color:#1b5e20;stop-opacity:1" />     </linearGradient>     <linearGradient id="containerGrad" x1="0%" y1="0%" x2="0%" y2="100%">       <stop offset="0%" style="stop-color:#f57c00;stop-opacity:1" />       <stop offset="100%" style="stop-color:#e65100;stop-opacity:1" />     </linearGradient>     <linearGradient id="serviceGrad" x1="0%" y1="0%" x2="0%" y2="100%">       <stop offset="0%" style="stop-color:#6a1b9a;stop-opacity:1" />       <stop offset="100%" style="stop-color:#4a148c;stop-opacity:1" />     </linearGradient>     <linearGradient id="configGrad" x1="0%" y1="0%" x2="0%" y2="100%">       <stop offset="0%" style="stop-color:#006064;stop-opacity:1" />       <stop offset="100%" style="stop-color:#00363a;stop-opacity:1" />     </linearGradient>     <linearGradient id="externalGrad" x1="0%" y1="0%" x2="0%" y2="100%">       <stop offset="0%" style="stop-color:#37474f;stop-opacity:1" />       <stop offset="100%" style="stop-color:#263238;stop-opacity:1" />     </linearGradient>     <filter id="shadow" x="-20%" y="-20%" width="140%" height="140%">       <feDropShadow dx="0" dy="2" stdDeviation="3" flood-opacity="0.3"/>     </filter>     <marker id="arrowhead" markerWidth="10" markerHeight="7" refX="9" refY="3.5" orient="auto">       <polygon points="0 0, 10 3.5, 0 7" fill="#90caf9" />     </marker>     <marker id="arrowhead2" markerWidth="10" markerHeight="7" refX="9" refY="3.5" orient="auto">       <polygon points="0 0, 10 3.5, 0 7" fill="#ffcc80" />     </marker>   </defs>      <!-- 背景 -->   <rect width="1400" height="1100" fill="#eceff1"/>      <!-- 标题 -->   <text x="700" y="35" text-anchor="middle" font-family="Arial, sans-serif" font-size="26" font-weight="bold" fill="#333">     Kubernetes 完整架构图（含所有核心组件）   </text>      <!-- ========== 外部访问层 ========== -->   <rect x="30" y="55" width="1340" height="120" rx="12" fill="url(#externalGrad)" filter="url(#shadow)"/>   <text x="700" y="85" text-anchor="middle" fill="white" font-family="Arial, sans-serif" font-size="16" font-weight="bold">     🌐 外部访问层 (External Access)   </text>      <!-- Ingress -->   <rect x="60" y="100" width="300" height="60" rx="8" fill="rgba(255,152,0,0.3)" stroke="#ff9800" stroke-width="2"/>   <text x="210" y="125" text-anchor="middle" fill="#ff9800" font-family="Arial, sans-serif" font-size="13" font-weight="bold">🔀 Ingress</text>   <text x="210" y="145" text-anchor="middle" fill="#ffe0b2" font-family="Arial, sans-serif" font-size="10">域名路由 + SSL证书 + 负载均衡</text>      <!-- 外部用户 -->   <rect x="400" y="100" width="200" height="60" rx="8" fill="rgba(255,255,255,0.1)" stroke="white" stroke-width="1"/>   <text x="500" y="125" text-anchor="middle" fill="white" font-family="Arial, sans-serif" font-size="12">👤 外部用户</text>   <text x="500" y="145" text-anchor="middle" fill="#b0bec5" font-family="Arial, sans-serif" font-size="10">https://api.example.com</text>      <!-- 外部服务类型 -->   <rect x="630" y="100" width="700" height="60" rx="8" fill="rgba(255,255,255,0.05)" stroke="#90caf9" stroke-width="1" stroke-dasharray="3,3"/>   <text x="980" y="125" text-anchor="middle" fill="#90caf9" font-family="Arial, sans-serif" font-size="11" font-weight="bold">Service Types: ClusterIP(内部) | NodePort(节点端口) | LoadBalancer(云负载) | ExternalName(外部域名)</text>   <text x="980" y="145" text-anchor="middle" fill="#bbdefb" font-family="Arial, sans-serif" font-size="10">NodePort: http://NodeIP:30080 → 映射到 Service:80</text>      <!-- ========== Cluster 层级 ========== -->   <rect x="30" y="190" width="1340" height="890" rx="15" fill="url(#clusterGrad)" filter="url(#shadow)"/>   <text x="700" y="220" text-anchor="middle" fill="white" font-family="Arial, sans-serif" font-size="18" font-weight="bold">     🔷 Kubernetes Cluster (集群边界)   </text>      <!-- ========== Master Node (控制平面) ========== -->   <rect x="60" y="240" width="1300" height="180" rx="12" fill="url(#masterGrad)" filter="url(#shadow)"/>   <text x="700" y="270" text-anchor="middle" fill="white" font-family="Arial, sans-serif" font-size="15" font-weight="bold">     🎛️ Master Node (控制平面) - 集群大脑，管理所有资源   </text>      <!-- API Server -->   <rect x="80" y="290" width="200" height="115" rx="8" fill="rgba(255,255,255,0.15)" stroke="white" stroke-width="2"/>   <text x="180" y="315" text-anchor="middle" fill="white" font-family="Arial, sans-serif" font-size="12" font-weight="bold">🔷 kube-apiserver</text>   <text x="180" y="335" text-anchor="middle" fill="#ffccbc" font-family="Arial, sans-serif" font-size="10">集群API网关/入口</text>   <text x="180" y="355" text-anchor="middle" fill="#ffccbc" font-family="Arial, sans-serif" font-size="9">认证、授权、访问控制</text>   <text x="180" y="375" text-anchor="middle" fill="#ffccbc" font-family="Arial, sans-serif" font-size="9">所有组件通信枢纽</text>   <text x="180" y="395" text-anchor="middle" fill="#ffccbc" font-family="Arial, sans-serif" font-size="9">kubectl → apiserver</text>      <!-- Scheduler -->   <rect x="300" y="290" width="180" height="115" rx="8" fill="rgba(255,255,255,0.15)" stroke="white" stroke-width="2"/>   <text x="390" y="315" text-anchor="middle" fill="white" font-family="Arial, sans-serif" font-size="12" font-weight="bold">📅 Scheduler</text>   <text x="390" y="335" text-anchor="middle" fill="#ffccbc" font-family="Arial, sans-serif" font-size="10">调度器</text>   <text x="390" y="355" text-anchor="middle" fill="#ffccbc" font-family="Arial, sans-serif" font-size="9">监控节点资源</text>   <text x="390" y="375" text-anchor="middle" fill="#ffccbc" font-family="Arial, sans-serif" font-size="9">Pod → 最优Node</text>   <text x="390" y="395" text-anchor="middle" fill="#ffccbc" font-family="Arial, sans-serif" font-size="9">调度策略：资源最多</text>      <!-- Controller Manager -->   <rect x="500" y="290" width="200" height="115" rx="8" fill="rgba(255,255,255,0.15)" stroke="white" stroke-width="2"/>   <text x="600" y="315" text-anchor="middle" fill="white" font-family="Arial, sans-serif" font-size="12" font-weight="bold">🎮 Controller Manager</text>   <text x="600" y="335" text-anchor="middle" fill="#ffccbc" font-family="Arial, sans-serif" font-size="10">控制器管理器</text>   <text x="600" y="355" text-anchor="middle" fill="#ffccbc" font-family="Arial, sans-serif" font-size="9">管理资源状态</text>   <text x="600" y="375" text-anchor="middle" fill="#ffccbc" font-family="Arial, sans-serif" font-size="9">故障检测与恢复</text>   <text x="600" y="395" text-anchor="middle" fill="#ffccbc" font-family="Arial, sans-serif" font-size="9">Deployment/StatefulSet</text>      <!-- etcd -->   <rect x="720" y="290" width="180" height="115" rx="8" fill="rgba(255,255,255,0.15)" stroke="white" stroke-width="2"/>   <text x="810" y="315" text-anchor="middle" fill="white" font-family="Arial, sans-serif" font-size="12" font-weight="bold">💾 etcd</text>   <text x="810" y="335" text-anchor="middle" fill="#ffccbc" font-family="Arial, sans-serif" font-size="10">键值存储 (类似Redis)</text>   <text x="810" y="355" text-anchor="middle" fill="#ffccbc" font-family="Arial, sans-serif" font-size="9">存储集群所有状态</text>   <text x="810" y="375" text-anchor="middle" fill="#ffccbc" font-family="Arial, sans-serif" font-size="9">Pod/Node/Service配置</text>   <text x="810" y="395" text-anchor="middle" fill="#ffccbc" font-family="Arial, sans-serif" font-size="9">不存应用数据</text>      <!-- Cloud Controller -->   <rect x="920" y="290" width="200" height="115" rx="8" fill="rgba(255,255,255,0.1)" stroke="#ffcc80" stroke-width="1" stroke-dasharray="3,3"/>   <text x="1020" y="315" text-anchor="middle" fill="#ffcc80" font-family="Arial, sans-serif" font-size="12" font-weight="bold">☁️ Cloud Controller</text>   <text x="1020" y="335" text-anchor="middle" fill="#ffe0b2" font-family="Arial, sans-serif" font-size="10">云控制器 (可选)</text>   <text x="1020" y="355" text-anchor="middle" fill="#ffe0b2" font-family="Arial, sans-serif" font-size="9">对接云平台API</text>   <text x="1020" y="375" text-anchor="middle" fill="#ffe0b2" font-family="Arial, sans-serif" font-size="9">AWS/Azure/GCP</text>   <text x="1020" y="395" text-anchor="middle" fill="#ffe0b2" font-family="Arial, sans-serif" font-size="9">LoadBalancer provision</text>      <!-- 箭头：Master 到 Worker -->   <path d="M 180 405 L 180 440" stroke="#90caf9" stroke-width="2" fill="none" marker-end="url(#arrowhead)"/>   <path d="M 390 405 L 390 440" stroke="#90caf9" stroke-width="2" fill="none" marker-end="url(#arrowhead)"/>   <path d="M 600 405 L 600 440" stroke="#90caf9" stroke-width="2" fill="none" marker-end="url(#arrowhead)"/>   <path d="M 810 405 L 810 440" stroke="#90caf9" stroke-width="2" fill="none" marker-end="url(#arrowhead)"/>      <!-- ========== Worker Node 1 ========== -->   <rect x="60" y="450" width="630" height="610" rx="12" fill="url(#workerGrad)" filter="url(#shadow)"/>   <text x="375" y="480" text-anchor="middle" fill="white" font-family="Arial, sans-serif" font-size="14" font-weight="bold">     🖥️ Worker Node 1 (工作节点) - cd-032112-*   </text>      <!-- Node 组件 -->   <rect x="80" y="495" width="590" height="50" rx="6" fill="rgba(255,255,255,0.1)" stroke="white" stroke-width="1"/>   <text x="120" y="515" fill="#ffeb3b" font-family="Arial, sans-serif" font-size="11" font-weight="bold">🏷️ 标签:</text>   <text x="180" y="515" fill="white" font-family="Arial, sans-serif" font-size="11">nodetype=common | usage=db | number=112</text>   <text x="120" y="535" fill="#81d4fa" font-family="Arial, sans-serif" font-size="10">🔧 kubelet (管理Pod) | 🔀 kube-proxy (网络代理/负载均衡) | 🐳 container runtime (Docker/containerd)</text>      <!-- Pod 1: Nginx (Deployment) -->   <rect x="80" y="560" width="280" height="230" rx="8" fill="url(#podGrad)" filter="url(#shadow)"/>   <text x="220" y="585" text-anchor="middle" fill="white" font-family="Arial, sans-serif" font-size="12" font-weight="bold">📦 Pod: nginx-xxx</text>   <text x="220" y="605" text-anchor="middle" fill="#c8e6c9" font-family="Arial, sans-serif" font-size="10">Controller: Deployment</text>      <!-- Pod 1 标签 -->   <rect x="95" y="620" width="250" height="35" rx="4" fill="rgba(0,0,0,0.2)"/>   <text x="220" y="635" text-anchor="middle" fill="#ffeb3b" font-family="Arial, sans-serif" font-size="9">🏷️ app=nginx, env=prod, version=v1.2</text>   <text x="220" y="650" text-anchor="middle" fill="#ffeb3b" font-family="Arial, sans-serif" font-size="9">pod-template-hash=55f89d4fd5</text>      <!-- Container -->   <rect x="100" y="665" width="240" height="50" rx="5" fill="url(#containerGrad)"/>   <text x="220" y="685" text-anchor="middle" fill="white" font-family="Arial, sans-serif" font-size="11" font-weight="bold">🐳 Container: nginx</text>   <text x="220" y="705" text-anchor="middle" fill="#ffe0b2" font-family="Arial, sans-serif" font-size="10">Image: nginx:latest</text>      <!-- Pod 1 说明 -->   <rect x="95" y="725" width="250" height="55" rx="4" fill="rgba(255,255,255,0.1)"/>   <text x="220" y="745" text-anchor="middle" fill="#b3e5fc" font-family="Arial, sans-serif" font-size="9">🔁 自动重启: 故障时重建</text>   <text x="220" y="760" text-anchor="middle" fill="#b3e5fc" font-family="Arial, sans-serif" font-size="9">📈 副本控制: 保持3个副本</text>   <text x="220" y="775" text-anchor="middle" fill="#b3e5fc" font-family="Arial, sans-serif" font-size="9">🔄 滚动更新: 平滑升级</text>      <!-- Pod 2: MySQL (StatefulSet) -->   <rect x="380" y="560" width="280" height="280" rx="8" fill="url(#podGrad)" filter="url(#shadow)"/>   <text x="520" y="585" text-anchor="middle" fill="white" font-family="Arial, sans-serif" font-size="12" font-weight="bold">📦 Pod: mysql-0</text>   <text x="520" y="605" text-anchor="middle" fill="#ce93d8" font-family="Arial, sans-serif" font-size="11" font-weight="bold">Controller: StatefulSet</text>      <!-- Pod 2 标签 -->   <rect x="395" y="620" width="250" height="35" rx="4" fill="rgba(0,0,0,0.2)"/>   <text x="520" y="635" text-anchor="middle" fill="#ffeb3b" font-family="Arial, sans-serif" font-size="9">🏷️ app=mysql, statefulset=mysql, ordinal=0</text>   <text x="520" y="650" text-anchor="middle" fill="#ffeb3b" font-family="Arial, sans-serif" font-size="9">稳定网络ID: mysql-0.mysql-svc</text>      <!-- Container -->   <rect x="400" y="665" width="240" height="50" rx="5" fill="url(#containerGrad)"/>   <text x="520" y="685" text-anchor="middle" fill="white" font-family="Arial, sans-serif" font-size="11" font-weight="bold">🐳 Container: mysql</text>   <text x="520" y="705" text-anchor="middle" fill="#ffe0b2" font-family="Arial, sans-serif" font-size="10">Image: mysql:8.0</text>      <!-- Volume -->   <rect x="395" y="725" width="250" height="45" rx="4" fill="rgba(3,169,244,0.3)" stroke="#03a9f4" stroke-width="1"/>   <text x="520" y="745" text-anchor="middle" fill="#81d4fa" font-family="Arial, sans-serif" font-size="10" font-weight="bold">💾 Volume (持久化存储)</text>   <text x="520" y="762" text-anchor="middle" fill="#b3e5fc" font-family="Arial, sans-serif" font-size="9">PVC → PV → 远程存储(OSS/NFS)</text>      <!-- Pod 2 说明 -->   <rect x="395" y="780" width="250" height="50" rx="4" fill="rgba(255,255,255,0.1)"/>   <text x="520" y="800" text-anchor="middle" fill="#ce93d8" font-family="Arial, sans-serif" font-size="9">🎯 有状态应用: 数据持久化</text>   <text x="520" y="815" text-anchor="middle" fill="#ce93d8" font-family="Arial, sans-serif" font-size="9">🔒 稳定标识: 网络ID和存储不变</text>   <text x="520" y="830" text-anchor="middle" fill="#ce93d8" font-family="Arial, sans-serif" font-size="9">⚠️ 注意: Redis内存数据不适合StatefulSet</text>      <!-- ConfigMap & Secret -->   <rect x="80" y="810" width="280" height="90" rx="6" fill="url(#configGrad)"/>   <text x="220" y="835" text-anchor="middle" fill="white" font-family="Arial, sans-serif" font-size="11" font-weight="bold">📋 ConfigMap + 🔐 Secret</text>   <text x="220" y="855" text-anchor="middle" fill="#80deea" font-family="Arial, sans-serif" font-size="9">配置分离: 修改配置 → 重载Pod (无需重新编译)</text>   <text x="220" y="875" text-anchor="middle" fill="#80deea" font-family="Arial, sans-serif" font-size="9">ConfigMap: 明文配置 (DB_HOST, PORT)</text>   <text x="220" y="890" text-anchor="middle" fill="#ff8a80" font-family="Arial, sans-serif" font-size="9">Secret: Base64编码敏感信息 (PASSWORD)</text>      <!-- ========== Worker Node 2 ========== -->   <rect x="710" y="450" width="630" height="610" rx="12" fill="url(#workerGrad)" filter="url(#shadow)"/>   <text x="1025" y="480" text-anchor="middle" fill="white" font-family="Arial, sans-serif" font-size="14" font-weight="bold">     🖥️ Worker Node 2 (工作节点) - cd-038118-*   </text>      <!-- Node 组件 -->   <rect x="730" y="495" width="590" height="50" rx="6" fill="rgba(255,255,255,0.1)" stroke="white" stroke-width="1"/>   <text x="770" y="515" fill="#ffeb3b" font-family="Arial, sans-serif" font-size="11" font-weight="bold">🏷️ 标签:</text>   <text x="830" y="515" fill="white" font-family="Arial, sans-serif" font-size="11">nodetype=common | app=goc | kubernetes.io/hostname=cd-038118-*</text>   <text x="770" y="535" fill="#81d4fa" font-family="Arial, sans-serif" font-size="10">🔧 kubelet | 🔀 kube-proxy (同节点Pod路由) | 🐳 container runtime</text>      <!-- Pod 3: GOC App -->   <rect x="730" y="560" width="280" height="200" rx="8" fill="url(#podGrad)" filter="url(#shadow)"/>   <text x="870" y="585" text-anchor="middle" fill="white" font-family="Arial, sans-serif" font-size="12" font-weight="bold">📦 Pod: goc-app-xxx</text>   <text x="870" y="605" text-anchor="middle" fill="#c8e6c9" font-family="Arial, sans-serif" font-size="10">Controller: Deployment</text>      <!-- Pod 3 标签 -->   <rect x="745" y="620" width="250" height="30" rx="4" fill="rgba(0,0,0,0.2)"/>   <text x="870" y="640" text-anchor="middle" fill="#ffeb3b" font-family="Arial, sans-serif" font-size="9">🏷️ app=goc, env=test, version=v1.2</text>      <!-- Multi-Containers -->   <rect x="750" y="660" width="240" height="40" rx="4" fill="url(#containerGrad)"/>   <text x="870" y="680" text-anchor="middle" fill="white" font-family="Arial, sans-serif" font-size="10" font-weight="bold">🐳 Main: goc:v1.2</text>   <text x="870" y="695" text-anchor="middle" fill="#ffe0b2" font-family="Arial, sans-serif" font-size="9">业务容器</text>      <rect x="750" y="710" width="240" height="35" rx="4" fill="rgba(255,152,0,0.5)"/>   <text x="870" y="730" text-anchor="middle" fill="white" font-family="Arial, sans-serif" font-size="9">🐳 Sidecar: log-collector</text>      <!-- Pod 4: Another GOC -->   <rect x="1030" y="560" width="280" height="200" rx="8" fill="url(#podGrad)" filter="url(#shadow)"/>   <text x="1170" y="585" text-anchor="middle" fill="white" font-family="Arial, sans-serif" font-size="12" font-weight="bold">📦 Pod: goc-app-yyy</text>   <text x="1170" y="605" text-anchor="middle" fill="#c8e6c9" font-family="Arial, sans-serif" font-size="10">Controller: Deployment</text>      <rect x="1045" y="620" width="250" height="30" rx="4" fill="rgba(0,0,0,0.2)"/>   <text x="1170" y="640" text-anchor="middle" fill="#ffeb3b" font-family="Arial, sans-serif" font-size="9">🏷️ app=goc, env=test, version=v1.2</text>      <rect x="1050" y="660" width="240" height="40" rx="4" fill="url(#containerGrad)"/>   <text x="1170" y="680" text-anchor="middle" fill="white" font-family="Arial, sans-serif" font-size="10" font-weight="bold">🐳 Main: goc:v1.2</text>      <rect x="1050" y="710" width="240" height="35" rx="4" fill="rgba(255,152,0,0.5)"/>   <text x="1170" y="730" text-anchor="middle" fill="white" font-family="Arial, sans-serif" font-size="9">🐳 Sidecar: log-collector</text>      <!-- Service 关联 -->   <rect x="730" y="780" width="580" height="80" rx="6" fill="url(#serviceGrad)" filter="url(#shadow)"/>   <text x="1020" y="805" text-anchor="middle" fill="white" font-family="Arial, sans-serif" font-size="12" font-weight="bold">🔗 Service: goc-svc (ClusterIP)</text>   <text x="1020" y="825" text-anchor="middle" fill="#e1bee7" font-family="Arial, sans-serif" font-size="10">Selector: app=goc → 自动关联这两个Pod作为后端</text>   <text x="1020" y="845" text-anchor="middle" fill="#e1bee7" font-family="Arial, sans-serif" font-size="10">负载均衡: 流量分发到 goc-app-xxx 和 goc-app-yyy</text>      <!-- ========== 网络流量说明 ========== -->   <rect x="60" y="920" width="1280" height="140" rx="8" fill="rgba(0,0,0,0.3)"/>   <text x="80" y="945" fill="white" font-family="Arial, sans-serif" font-size="13" font-weight="bold">🔄 网络流量路径 & 核心机制:</text>      <text x="80" y="975" fill="#90caf9" font-family="Arial, sans-serif" font-size="11">     外部访问: 用户 → Ingress(域名) → Service(NodePort/LoadBalancer) → Pod   </text>   <text x="80" y="995" fill="#81c784" font-family="Arial, sans-serif" font-size="11">     内部通信: Pod A → Service(ClusterIP) → kube-proxy(负载均衡) → Pod B (同节点优先路由)   </text>   <text x="80" y="1015" fill="#ffcc80" font-family="Arial, sans-serif" font-size="11">     Pod IP不稳定: Pod重建后IP变化 → 通过Service稳定访问 → DNS: mysql-svc.production.svc.cluster.local   </text>   <text x="80" y="1035" fill="#ce93d8" font-family="Arial, sans-serif" font-size="11">     有状态应用: StatefulSet保证每个Pod有稳定网络ID(mysql-0.mysql-svc)和持久化存储   </text>   <text x="80" y="1055" fill="#80deea" font-family="Arial, sans-serif" font-size="11">     配置管理: ConfigMap/Secret挂载为环境变量或文件 → 修改后重载Pod → 无需重新编译镜像   </text> </svg>

## Kubernetes 组件
### Node
Node：节点，一个物理机或者一台虚拟机。

### Pod
Pod 是 Kubernetes 的最小调度单元，可以理解为容器的抽象。一个 Pod 就是一个或者多个应用容器的组合。它创建了一个容器的运行环境，在这个环境中容器可以共享一些资源，比如网络、存储和运行时的一些配置等等。

假设我们系统包括一个应用程序和一个数据库，就可以将应用程序和数据库分别放到两个不同的 Pod 中，一般情况下一个 Pod 中只运行一个容器，这样可以更好地实现应用程序的解耦和扩展。

> 一个 Pod 中也是可以运行多个容器的，一般仅限于这些容器是高度耦合的情况，它们之间为了共享一些配置或者资源，不得不将它们放到一个容器中
>

应用程序要访问数据库的话，只需要知道数据库的 IP 地址，这里的 IP 地址是 Pod 在创建的时候自动创建的，是一个集群内部的 IP 地址（也就是无法从集群外部访问），Pod 之间通过这些 IP 地址进行通信。

#### Pod IP 不稳定问题
Pod 并不是稳定的实体，它们非常容易被创建或者销毁，比如发生故障的时候，Kubernetes会自动将发生故障的 Pod 销毁掉然后创建一个新的 Pod 替代它，这时候 IP 也会重新分配，如果应用程序还用原来的 IP 来访问的话就找不到。

### Service
为了解决这个问题，Kubernetes 提供了一个名为 Service 的资源对象，它可以将一组 Pod 封装成一个服务，这个服务通过一个统一的入口来访问。

就比如上面的场景，我们分别将应用程序和数据库两组 Pod 封装成两个 Service，这样应用程序就可以通过 Service 的 IP 地址访问数据库（有点像路由器和反向代理），即使 Pod 的 IP 地址发生了变化，Service 的 IP 地址也不会发生变化，Service 会自动将请求转发到其它健康的 Pod 上。

<img src="C:/Users/16658/AppData/Roaming/Typora/typora-user-images/image-20260328175353912.png" alt="image-20260328175353912"  />

> 正向代理：代理的是 C 端（客户端），S 端不知道 C 端的 IP。C 端发送请求给代理服务器，再由代理服务器向 S 端（服务端）发送请求，S 端就会将数据响应给代理服务器，再由代理服务器将数据传输给 C 端。（科学上网使用的就是正向代理）
>
> 反向代理：代理的是 S 端，C 端不知道 S 端的 IP。S 端可能有多台服务器，C 端向代理服务器发送请求，代理服务器向 S 端的其中一台服务器发送请求，服务器将数据响应给代理服务器，再由代理服务器将数据传输给 C 端。（比如 Nginx 的负载均衡）
>
> 使用代理模式的好处：
>
> 1. 隐藏真实 IP，隐私
> 2. 在代理服务器上设置缓存可以加快访问
> 3. 突破网络限制，有些网站会限制某些地区的 IP 访问，使用化茧成蝶可以突破限制。（还是科学上网）
>
> 坏处：
>
> 1. C 端和 S 端的 IP 对彼此隐藏了，但会将 IP 暴露给代理服务器
> 2. 网络链路多了代理服务器的节点，降低访问速度
>

#### 内部服务和外部服务
内部服务：不能暴露或者不需要暴露给外部的服务，比如数据库、缓存、消息队列等，这些服务只需要在集群内部访问就可以了。

外部服务：后端的 API 接口或者前端界面等等，这些就是需要暴露给用户的服务。

外部服务常见的类型有 ExternalName、LoadBalancer、NodePort、ClusterIP，其中 NodePort 是我们常用的类型，它会在节点上开放一个端口，然后将这个端口映射到 Service 的 IP 地址和端口上，这样就可以通过节点的 IP 地址和端口来访问 Service 了。是不是感觉有点熟悉？http://localhost:8080，想起来了吗？

在开发和测试阶段使用 IP 和端口号的方式是没有问题的，但在生产环境中通常是用域名来访问服务的，这时就用到了另一个资源对象 Ingress。

### Ingress
Ingress 是用于管理从集群外部访问集群内部服务的入口和方式，可以通过 Ingress 配置不同的转发规则，从而根据不同的规则来访问不同的 Service 以及 Service 所对应的 Pod。还可以通过 Ingress 来配置域名，这样就可以从集群外部使用域名和访问 Service。

> Ingress 也可以用来配置 SSL 证书或者负载均衡。
>

![](https://cdn.nlark.com/yuque/0/2024/png/38959865/1713182820931-ed5d436a-9601-4cfa-b212-b8f8b1d867cd.png)

### ConfigMap
原来我们的应用程序需要访问数据库的话，一般的做法是将数据库的地址和端口等连接信息写到配置文件或者环境变量中，然后在应用程序中读取这些配置信息，这样配置信息就和应用程序耦合在一起了，当数据库的地址或者端口发生变化，我们就得修改应用程序的配置信息然后重新编译部署，这样不仅麻烦，而且对于一些需要不间断运行的服务来说是不能接受的（比如你深夜肚子饿了想点外卖而服务器却在重新编译部署）。

为了解决这个问题，Kubernetes 提供了一个 ConfigMap 组件，它可以将配置信息封装起来，然后就可以在应用程序中读取和使用。将配置信息和应用程序的镜像内容分离开，当数据库的地址和端口发生变化的时候，只需要修改 ConfigMap对象中的配置信息，**然后重新加载 Pod**，不需要重新编译和部署应用程序。

> 可以理解为给数据库加了个反向代理
>

### Secret
但 ConfigMap 有个问题，就是它的配置信息是明文的，如果用户名和密码存在 ConfigMap 中是有风险的，于是 Kubernetes 提供了 Secret 组件，用于封装敏感信息，会将配置信息 Base64 编码，但 Base64 编码很容易解码，浏览器随便搜个解码器就能得到原文，所以还需要配合其它手段来确保安全性。

> 真是一层套一层
>

### Volume
Pod 被销毁或重启时数据也跟着消失，这对于需要持久化存储的应用程序比如数据库肯定是不行的，Kubernetes 提供了 Volume 组件，它可以将一些持久化存储的资源挂载到集群中的本地磁盘上，或者挂载到集群外部的远程存储上（比如 OSS）。

### Deployment
如果服务端只有一个节点的话，这个节点发生故障就会导致服务宕机（即单点故障），无法实现高可用性。这个好解决，既然一个不够，那就多复制几个，当一个节点发生故障的时候，Service 就会自动将请求转发到另一个节点。（这里的 Service 指的是上面的 [Service](#T9uOD) 组件）

Deployment 就是用来定义和管理应用程序的副本数量以及应用程序的更新策略，将一个或者多个 Pod 组合在一起，简化应用程序的部署和更新操作，还可以副本控制、滚动更新、自动扩缩容等。

副本控制：定义和管理应用程序的副本数量，比如定义一个应用程序副本数量为 3，当其中一个发生故障时，就会生成一个新的副本来替换坏副本，始终保持有 3 个副本在集群中运行。

滚动更新：定义和管理应用程序的更新策略，使用新版本替换旧版本，确保应用程序的平滑升级。

> 平滑升级：不对用户的使用造成中断或不便的升级
>

### StatefulSet
除了应用程序，数据库也有故障、升级和更新维护的时候，数据库停了服务也停了，所以数据库也需要采取多副本的方式来保证高可用性，但一般不使用 Deployment 来实现数据库的多副本，因为**数据库的多副本之间是有状态的**，就是每个副本的数据存在差异（状态不同），需要确保数据的一致性，可以把数据写到一个共享的存储中或者同步不同副本之间的数据。

对于这一类有状态的应用程序，Kubernetes 提供了 StatefulSet 组件来管理，StatefulSet 跟 Deployment 非常类似，也提供了定义和管理应用程序副本数量和动态扩缩容等功能，**此外它还保证了每个副本都有自己稳定的网络标识符和持久化存储**，数据库、缓存、消息队列等这些**有状态的应用**以及**保留了会话状态的应用**一般都需要使用 StatefulSet 而不是 Deployment。

但 StatefulSet 部署的过程比较复杂和繁琐，而且并不是所有的有状态应用都适合用 StatefulSet 来部署，比如Redis，当 Redis 实例需要使用内存进行数据存储，并且数据存储在**实例的内存**中而不是持久性存储中时，如果这时候实例重启，Redis 数据丢失，这时候 StatefulSet 就无法保证数据的持久性（还是有点不理解）。更加通用和简单的方式是把有状态的应用程序从 Kubernetes 中剥离出来，在集群外单独部署。（ToLearn：集群外单独部署有状态应用的方法）

##  Kubernetes 架构
Kubernetes 是典型的 Master-Worker 架构， Master Node 负责管理整个集群，Worker Node 负责运行应用程序和服务。

> 官方定义：Kubernetes 通过将容器放入在节点上运行的 Pod 中来执行工作负载。
>

### Worker Node
下面这个 Node 就运行了两个 Pod，这里的 Node 就是集群中实际完成工作的节点，也就是 Worker-Node。

![](https://cdn.nlark.com/yuque/0/2024/png/38959865/1713189041853-6cfe420c-5a24-4fcd-8871-0168cc910887.png)

为了对外提供服务，每个 Node 上都会包含三个组件，kubelet、kube-proxy 和 container runtime。

#### container runtime
翻译为“运行环境”，可以理解为一个运行容器的实例，负责拉取容器镜像、创建容器、启动或者停止容器等等，所有的容器都需要通过 container runtime 来运行，所以每个节点上都需要安装 container runtime。

#### kubelet
负责管理和维护每个节点上的 Pod 并确保它们按照预期运行，它也会定期从 apiserver 组件接收新的或者修改后的 Pod 规范，同时会监控工作节点的运行情况，然后将信息汇报给 apiserver。

#### kube-proxy
负责为 Pod 对象提供网络代理和负载均衡服务。



一般 Kubernetes 集群包含多个节点，节点之间通过 Service 通信，这就需要一个负载均衡器来接收请求，然后再将请求发送到不同节点上。kube-proxy 就是负责这个功能的组件，它在每个 Node 上启动一个网络代理，使发往 Service 的流量**路由**到正确的后端 Pod，比如 Node A 的应用程序要访问数据库，应用程序对数据库的访问请求并不会被随机路由到别的节点的数据库中，kube-proxy 会把请求路由到与应用程序同一个节点（也就是 Node A）的数据库 Pod 中。

> “路由”也是一个动词
>

### Master Node
![](https://cdn.nlark.com/yuque/0/2024/png/38959865/1713198479570-094946e8-3750-4848-863c-4bc26373ffb6.png)

#### kube-apiserver
负责提供 Kubernetes 集群的 API 接口服务，所有的组件通过 apiserver 通信，用户在集群中部署应用时也需要使用客户端与 apiserver 进行交互。apiserver 就像一个集群的网关，是整个系统的入口，所有的请求都会先经过它，再由它分发给不同的组件进行处理。

除了提供 API 接口之外，apiserver 还负责对所有对象的增删改查等操作进行认证、授权和访问控制。apiserver 接收到请求时会先验证请求的合法性，验证通过之后才会将请求转发给 **Scheduler** 处理。

#### Scheduler 调度器
负责监控集群中所有节点的资源使用情况，然后根据调度策略将 Pod 调度到合适的节点上运行。

> 调度策略：比如新增一个新 Pod 时，将 Pod 调度到空闲资源最多的节点上
>

#### Controller Manager 控制器管理器
负责管理集群中各种资源对象（比如 Node、Pod、Service）的状态，然后根据状态做出相应的响应。比如集群中有一个节点发生故障时，得有一个机制来监控和检测这个故障然后处理故障，比如重启 Pod 或者使用新 Pod 替代，这就是 Controller Manager 要做的事。 

#### etcd
高可用**键-值存储系统**，类似 Redis，用于存储集群中所有资源对象的状态信息，每新增或者崩掉一个 Pod，这些信息都会被记录到 etcd 中，使用命令行查询集群状态时就是通过 etcd 来获取的。

> etcd 一般只存储集群中应用程序的状态信息，不存储应用程序的数据
>

#### Cloud Controller Manager
一个云平台相关的控制器，负责与云平台（例如 Google GKE、Microsoft AKS、Amazon EKS）的 API 进行交互，并且对不同平台都提供一致的管理接口。

