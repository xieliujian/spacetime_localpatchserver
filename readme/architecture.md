[← 返回总览](../README.md)

# 架构设计

## 整体架构

```
┌─────────────────┐
│  Unity 编辑器   │ ──── HTTP POST ───→ ┌──────────────┐
└─────────────────┘                     │              │
                                        │  Gin Server  │
┌─────────────────┐                     │              │
│  游戏客户端     │ ──── HTTP GET ────→ │  (Go 1.21+)  │
└─────────────────┘                     │              │
                                        └──────┬───────┘
┌─────────────────┐                            │
│  Web 管理界面   │ ──── HTTP ────────→        │
└─────────────────┘                            ↓
                                        ┌──────────────┐
                                        │  文件系统    │
                                        │  metadata.json│
                                        └──────────────┘
```

## 模块划分

### 1. internal/config
- 职责：加载和验证 YAML 配置
- 核心类型：`Config`, `ServerConfig`, `AuthConfig`, `StorageConfig`
- 关键方法：`Load(path)`, `validate()`

### 2. internal/storage
- 职责：版本管理、文件存储、元数据持久化
- 核心类型：`Manager`, `Metadata`, `VersionInfo`, `FileInfo`
- 关键方法：
  - `NextVersion()` — 获取下一个版本号
  - `AddVersion(info)` — 添加新版本
  - `GetLatestVersion()` — 获取最新版本
  - `DeleteVersion(version)` — 删除指定版本

### 3. internal/auth
- 职责：API Key 认证
- 核心类型：`APIKeyMiddleware`
- 验证方式：检查 HTTP Header `X-API-Key`

### 4. internal/api
- 职责：HTTP API 处理
- 模块：
  - `config.go` — 配置查询
  - `version.go` — 版本查询
  - `upload.go` — 资源上传
  - `download.go` — 资源下载

### 5. internal/web
- 职责：Web 管理界面
- 实现：静态 HTML + JavaScript

### 6. cmd/server
- 职责：服务器入口
- 功能：初始化各模块、注册路由、启动服务

## 数据流

### 上传流程
```
Unity 编辑器
  ↓ POST /api/upload (multipart/form-data)
  ↓ Header: X-API-Key
认证中间件验证
  ↓
UploadHandler
  ↓ 1. 获取下一个版本号
  ↓ 2. 创建版本目录 data/versions/{version}/
  ↓ 3. 保存所有文件
  ↓ 4. 更新 metadata.json
  ↓
返回 {"version": N, "uploaded_files": [...]}
```

### 下载流程
```
游戏客户端
  ↓ GET /api/config
  ↓ 获取 patch_server_url 和 current_version
  ↓
  ↓ GET /api/version/latest
  ↓ 对比本地版本号
  ↓
  ↓ GET /api/download/{version}/{filepath}
  ↓ 下载资源文件
  ↓
应用资源
```

## 存储结构

```
data/
├── metadata.json          # 版本元数据
└── versions/
    ├── 1/                 # 版本 1
    │   ├── scene1.unity3d
    │   └── ui/
    │       └── main.prefab
    ├── 2/                 # 版本 2
    │   └── ...
    └── 3/                 # 版本 3
        └── ...
```

### metadata.json 格式

```json
{
  "versions": [
    {
      "version": 1,
      "upload_time": "2026-04-20T10:00:00Z",
      "total_size": 1048576,
      "file_count": 5,
      "files": [
        {"path": "scene1.unity3d", "size": 524288},
        {"path": "ui/main.prefab", "size": 524288}
      ]
    }
  ]
}
```

## 并发安全

- `storage.Manager` 使用 `sync.Mutex` 保护元数据读写
- 版本号递增通过锁保证原子性
- 文件上传失败时自动回滚（删除版本目录）

## 扩展性

当前为单体架构，未来可扩展：
- 对象存储（OSS/S3）替代文件系统
- Redis 缓存元数据
- 多实例部署 + 负载均衡
- CDN 加速下载

[← 返回总览](../README.md)
