[← 返回总览](../README.md)

# 开发指南

---

## 环境准备

### 必需工具
- Go 1.21+
- Git
- 代码编辑器（推荐 VS Code + Go 插件）

### 可选工具
- curl / Postman（API 测试）
- Docker（容器化测试）

---

## 项目结构

```
spacetime_localpatchserver/
├── cmd/
│   └── server/
│       └── main.go              # 服务器入口
├── internal/
│   ├── api/                     # HTTP API 处理
│   │   ├── config.go
│   │   ├── version.go
│   │   ├── upload.go
│   │   └── download.go
│   ├── auth/
│   │   └── middleware.go        # API Key 认证
│   ├── config/
│   │   └── config.go            # 配置加载
│   ├── storage/
│   │   ├── manager.go           # 存储管理器
│   │   └── metadata.go          # 元数据结构
│   └── web/
│       └── handler.go           # Web 界面
├── web/
│   └── index.html               # 管理界面
├── config.yaml                  # 配置文件
├── go.mod
└── README.md
```

---

## 常用命令

### 安装依赖
```bash
go mod tidy
```

### 运行开发服务器
```bash
go run cmd/server/main.go -config config.yaml
```

### 运行测试
```bash
# 所有测试
go test ./...

# 特定包
go test ./internal/storage -v
go test ./internal/auth -v

# 带覆盖率
go test ./... -cover
```

### 代码格式化
```bash
go fmt ./...
```

### 代码检查
```bash
go vet ./...
```

### 构建
```bash
# 当前平台
go build -o patchserver cmd/server/main.go

# 跨平台
GOOS=linux GOARCH=amd64 go build -o patchserver cmd/server/main.go
GOOS=windows GOARCH=amd64 go build -o patchserver.exe cmd/server/main.go
```

---

## 开发流程

### 1. 添加新功能

1. 在对应模块创建文件
2. 编写测试（TDD）
3. 实现功能
4. 运行测试确保通过
5. 提交代码

### 2. 修改 API

1. 更新 `internal/api/` 下的处理器
2. 更新 `cmd/server/main.go` 路由
3. 更新 API 文档 `readme/api.md`
4. 测试接口

### 3. 修改存储逻辑

1. 更新 `internal/storage/manager.go`
2. 更新测试 `internal/storage/manager_test.go`
3. 确保元数据兼容性

---

## 测试

### 单元测试示例

```go
// internal/storage/manager_test.go
func TestManager_NextVersion(t *testing.T) {
    tmpDir := t.TempDir()
    mgr, err := NewManager(tmpDir)
    if err != nil {
        t.Fatalf("NewManager failed: %v", err)
    }

    v1 := mgr.NextVersion()
    if v1 != 1 {
        t.Errorf("expected version 1, got %d", v1)
    }
}
```

### API 测试

```bash
# 获取配置
curl http://localhost:8080/api/config

# 上传文件
curl -X POST http://localhost:8080/api/upload \
  -H "X-API-Key: dev-api-key-change-in-production" \
  -F "files=@test.txt"

# 获取版本列表
curl -H "X-API-Key: dev-api-key-change-in-production" \
  http://localhost:8080/api/versions

# 下载文件
curl http://localhost:8080/api/download/1/test.txt -o downloaded.txt
```

---

## 调试

### 启用 Gin 调试模式

```go
// cmd/server/main.go
gin.SetMode(gin.DebugMode)  // 开发环境
gin.SetMode(gin.ReleaseMode) // 生产环境
```

### 查看日志

Gin 默认输出请求日志：
```
[GIN] 2026/04/20 - 10:30:00 | 200 |   1.234567ms |   127.0.0.1 | GET      "/api/config"
```

### 使用 Delve 调试器

```bash
# 安装
go install github.com/go-delve/delve/cmd/dlv@latest

# 调试
dlv debug cmd/server/main.go -- -config config.yaml
```

---

## 代码规范

### 命名约定
- 包名：小写，单数
- 文件名：小写，下划线分隔
- 类型名：大驼峰（导出）或小驼峰（私有）
- 函数名：大驼峰（导出）或小驼峰（私有）

### 注释
- 导出的类型和函数必须有注释
- 注释以类型/函数名开头

```go
// Manager 管理资源版本和文件存储
type Manager struct {
    // ...
}

// NextVersion 返回下一个可用的版本号
func (m *Manager) NextVersion() int {
    // ...
}
```

### 错误处理
- 使用 `fmt.Errorf` 包装错误
- 提供上下文信息

```go
if err != nil {
    return fmt.Errorf("failed to load metadata: %w", err)
}
```

---

## 常见问题

### Q: 端口被占用
```bash
# 查找占用端口的进程
lsof -i :8080  # macOS/Linux
netstat -ano | findstr :8080  # Windows

# 修改 config.yaml 中的端口
```

### Q: 依赖下载失败
```bash
# 使用代理
export GOPROXY=https://goproxy.cn,direct
go mod tidy
```

### Q: 测试失败
```bash
# 清理缓存重新测试
go clean -testcache
go test ./...
```

[← 返回总览](../README.md)
