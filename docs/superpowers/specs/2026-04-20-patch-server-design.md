# 资源热更服务器设计文档

**日期：** 2026-04-20  
**项目：** spacetime_localpatchserver  
**技术栈：** Go + Gin

---

## 概述

为 Unity 游戏项目提供资源热更服务器。Unity 编辑器打包资源后通过 HTTP API 上传到服务器，游戏客户端启动时从远程配置获取服务器地址，对比版本号后按需下载最新资源。

---

## 需求

- Unity 编辑器打包后通过 HTTP API 上传资源（API Key 认证）
- 服务器按整数版本号（1, 2, 3...）管理资源，每个版本独立文件夹
- 游戏客户端通过远程配置 API 获取服务器地址，对比版本号后下载资源
- 提供 Web 管理界面查看版本列表、上传资源、删除版本
- 单体 Go 服务器，无需外部数据库

---

## 项目结构

```
spacetime_localpatchserver/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── config.go
│   │   ├── upload.go
│   │   ├── download.go
│   │   └── version.go
│   ├── auth/
│   │   └── middleware.go
│   ├── storage/
│   │   └── manager.go
│   └── web/
│       └── handler.go
├── web/
│   ├── index.html
│   ├── css/
│   └── js/
├── data/
│   ├── versions/
│   │   ├── 1/
│   │   ├── 2/
│   │   └── ...
│   └── metadata.json
├── config.yaml
└── go.mod
```

---

## 核心组件

### HTTP Server
- 使用 Gin 框架
- 端口可配置（默认 8080）
- 路由分组：`/api/config`、`/api/upload`、`/api/download`、`/api/versions`、`/web`

### Storage Manager
- 管理 `data/versions/` 目录
- 维护 `metadata.json`（记录每个版本的上传时间、文件列表、大小）
- 版本号自增逻辑：读取当前最大版本号 +1，加锁保证原子性

### Auth Middleware
- 检查 HTTP Header `X-API-Key`
- 上传和版本管理 API 需要认证
- 下载和配置 API 无需认证

### Web UI
- 单页应用，纯静态 HTML/CSS/JS，内嵌到 Go 二进制（embed）
- 功能：版本列表、上传资源（拖拽/选择）、删除版本、版本详情

---

## API 设计

### 配置 API（无需认证）
```
GET /api/config
Response: {
  "patch_server_url": "http://localhost:8080",
  "current_version": 3
}
```

### 最新版本查询（无需认证）
```
GET /api/version/latest
Response: {
  "version": 3,
  "upload_time": "2026-04-20T10:30:00Z",
  "total_size": 1048576,
  "file_count": 15
}
```

### 下载资源（无需认证）
```
GET /api/download/:version/:filepath
例如: GET /api/download/3/scenes/level1.unity3d
```
- 直接返回文件流
- 支持断点续传（Range Header）

### 上传资源（需要 API Key）
```
POST /api/upload
Header: X-API-Key: your-api-key
Body: multipart/form-data
  - files: 多个文件
  - version: (可选) 指定版本号，不传则自动递增
Response: {
  "version": 4,
  "uploaded_files": ["scene1.unity3d", "ui/main.prefab"],
  "total_size": 2097152
}
```

### 版本管理（需要 API Key）
```
GET    /api/versions        # 获取所有版本列表
DELETE /api/versions/:id    # 删除指定版本
GET    /api/versions/:id    # 获取版本详情（文件列表）
```

---

## 配置文件

```yaml
server:
  port: 8080
  patch_server_url: "http://localhost:8080"

auth:
  api_key: "your-secret-api-key"

storage:
  data_dir: "./data"
```

---

## 数据流

### 客户端启动热更流程
1. 请求 `GET /api/config` 获取热更服务器地址
2. 请求 `GET /api/version/latest` 获取最新版本号
3. 对比本地版本号（存储在 PlayerPrefs 或本地文件）
4. 如果服务器版本更新，下载 `GET /api/download/:version/:filepath`
5. 下载完成后更新本地版本号

### Unity 编辑器上传流程
1. 编辑器打包资源到临时目录
2. 调用 `POST /api/upload`，Header 携带 API Key
3. 批量上传所有文件（multipart/form-data）
4. 服务器创建新版本文件夹，保存文件，更新 metadata.json
5. 返回新版本号给编辑器

---

## 错误处理

| 场景 | HTTP 状态码 |
|------|------------|
| API Key 错误 | 401 Unauthorized |
| 文件过大 | 413 Payload Too Large |
| 磁盘空间不足 | 507 Insufficient Storage |
| 版本/文件不存在 | 404 Not Found |
| 部分文件上传失败 | 回滚，删除未完成版本文件夹 |

并发上传时使用文件锁保证版本号原子性递增。

---

## 测试策略

**单元测试：**
- Storage Manager：版本号生成、文件保存、metadata 更新
- Auth Middleware：API Key 验证
- API Handlers：请求/响应处理

**集成测试：**
- 完整上传流程（模拟 Unity 编辑器）
- 完整下载流程（模拟客户端检查版本并下载）
- 并发上传测试

---

## 部署

**开发：**
```bash
go run cmd/server/main.go -config config.yaml
```

**生产：**
```bash
go build -o patchserver cmd/server/main.go
./patchserver -config config.yaml
```

**Docker：**
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o patchserver cmd/server/main.go

FROM alpine:latest
COPY --from=builder /app/patchserver /patchserver
COPY config.yaml /config.yaml
EXPOSE 8080
CMD ["/patchserver", "-config", "/config.yaml"]
```

---

## 依赖

- `github.com/gin-gonic/gin` — HTTP 框架
- 标准库（无需外部数据库）
