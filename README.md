# spacetime_localpatchserver

Unity 资源热更服务器 — 支持编辑器上传资源、客户端按版本下载、Web 管理界面。

## 技术文档

| 文档 | 说明 |
|------|------|
| [架构设计](readme/architecture.md) | 整体架构、模块划分、数据流 |
| [API 文档](readme/api.md) | 所有 HTTP 接口说明 |
| [配置说明](readme/config.md) | config.yaml 各字段说明 |
| [部署指南](readme/deployment.md) | 开发、生产、Docker 部署方式 |
| [开发指南](readme/development.md) | 本地开发、测试、构建命令 |

## 快速开始

```bash
# 安装 Go 1.21+
go mod tidy
go run cmd/server/main.go -config config.yaml
```

访问 `http://localhost:8080` 打开管理界面。
