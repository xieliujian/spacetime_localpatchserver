[← 返回总览](../README.md)

# 部署指南

---

## 开发环境部署

### 前置要求
- Go 1.21+
- Git

### 步骤

1. **克隆代码**
```bash
git clone <repository>
cd spacetime_localpatchserver
```

2. **安装依赖**
```bash
go mod tidy
```

3. **配置**
```bash
cp config.yaml config.local.yaml
# 编辑 config.local.yaml 修改配置
```

4. **运行**
```bash
go run cmd/server/main.go -config config.local.yaml
```

5. **验证**
```bash
curl http://localhost:8080/api/config
```

---

## 生产环境部署

### 方式一：二进制部署

1. **构建**
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o patchserver cmd/server/main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o patchserver.exe cmd/server/main.go
```

2. **部署**
```bash
# 上传到服务器
scp patchserver user@server:/opt/patchserver/
scp config.yaml user@server:/opt/patchserver/

# 启动
ssh user@server
cd /opt/patchserver
./patchserver -config config.yaml
```

3. **配置 systemd 服务（Linux）**

创建 `/etc/systemd/system/patchserver.service`：

```ini
[Unit]
Description=Patch Server
After=network.target

[Service]
Type=simple
User=patchserver
WorkingDirectory=/opt/patchserver
ExecStart=/opt/patchserver/patchserver -config /opt/patchserver/config.yaml
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

启动服务：
```bash
sudo systemctl daemon-reload
sudo systemctl enable patchserver
sudo systemctl start patchserver
sudo systemctl status patchserver
```

---

### 方式二：Docker 部署

1. **创建 Dockerfile**

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o patchserver cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/patchserver .
COPY config.yaml .
COPY web ./web
EXPOSE 8080
CMD ["./patchserver", "-config", "config.yaml"]
```

2. **构建镜像**
```bash
docker build -t patchserver:latest .
```

3. **运行容器**
```bash
docker run -d \
  --name patchserver \
  -p 8080:8080 \
  -v $(pwd)/data:/root/data \
  -v $(pwd)/config.yaml:/root/config.yaml \
  patchserver:latest
```

4. **使用 docker-compose**

创建 `docker-compose.yml`：

```yaml
version: '3.8'
services:
  patchserver:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./data:/root/data
      - ./config.yaml:/root/config.yaml
    restart: unless-stopped
```

启动：
```bash
docker-compose up -d
```

---

## Nginx 反向代理

```nginx
server {
    listen 80;
    server_name patch.example.com;

    client_max_body_size 500M;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

---

## 生产环境检查清单

- [ ] 修改默认 API Key
- [ ] 配置 HTTPS（Let's Encrypt）
- [ ] 设置防火墙规则
- [ ] 配置日志轮转
- [ ] 设置磁盘空间监控
- [ ] 配置自动备份 metadata.json
- [ ] 限制上传文件大小
- [ ] 配置 CDN（可选）

---

## 监控和日志

### 日志查看
```bash
# systemd
sudo journalctl -u patchserver -f

# Docker
docker logs -f patchserver
```

### 健康检查
```bash
curl http://localhost:8080/api/config
```

### 磁盘空间监控
```bash
df -h /opt/patchserver/data
```

[← 返回总览](../README.md)
