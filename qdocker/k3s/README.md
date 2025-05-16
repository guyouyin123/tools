[toc]

# k3s启动教程

## 步骤 1: 创建 docker-compose.yml 文件

```yaml
version: '3'
services:
  k3s:
    image: rancher/k3s:v1.27.4-k3s1
    container_name: k3s
    privileged: true
    volumes:
      - ./k3s-data:/var/lib/rancher/k3s
      - /etc/machine-id:/etc/machine-id
      - /dev/mapper:/dev/mapper
    ports:
      - "6443:6443"
      - "8080:80"
      - "8443:443"
    environment:
      - K3S_KUBECONFIG_MODE=644
      - K3S_KUBECONFIG_OUTPUT=/output/kubeconfig.yaml
    command: server --disable traefik --disable metrics-server
```

## 步骤 2: 启动 K3s

```bash
mkdir k3s-data
docker-compose up -d
```

## 步骤 3: 获取 kubeconfig

等待容器启动后（约30秒），从容器中复制 kubeconfig 文件：

```bash
docker cp k3s:/output/kubeconfig.yaml ./kubeconfig.yaml
```

## 步骤 4: 配置 kubectl

```bash
export KUBECONFIG=$(pwd)/kubeconfig.yaml
export KUBECONFIG=/Users/jeff/myself/tools/qdocker/k3s/kubeconfig.yaml
kubectl get nodes -A
```

## 步骤 5：开启kubectl 代理，解决 Kubernetes API 未授权问题
便于本地调试
```bash
kubectl proxy --port=8080
访问:GET http://localhost:8080/api/v1/pods
```

## 注意事项

1. 这个配置禁用了 Traefik 和 metrics-server 以简化部署
2. 数据会持久化在本地 `k3s-data` 目录中
3. 如果需要 ARM64 兼容的镜像，K3s 会自动处理
4. 端口映射：
   - 6443: Kubernetes API
   - 8080: 将来可以用于 Ingress HTTP
   - 8443: 将来可以用于 Ingress HTTPS