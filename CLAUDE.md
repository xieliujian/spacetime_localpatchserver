# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`spacetime_localpatchserver` 是一个 Unity 资源热更服务器，使用 Go + Gin 实现。

## Commands

**开发运行：**
```bash
go run cmd/server/main.go -config config.yaml
```

**构建：**
```bash
go build -o patchserver cmd/server/main.go
```

**测试：**
```bash
go test ./...
```

**单元测试（特定包）：**
```bash
go test ./internal/storage -v
go test ./internal/auth -v
```

## Architecture

**核心组件：**
- `internal/config` — YAML 配置加载
- `internal/storage` — 版本管理和文件存储（文件系统 + metadata.json）
- `internal/auth` — API Key 认证中间件
- `internal/api` — REST API 处理器（config, version, upload, download）
- `internal/web` — Web 管理界面
- `cmd/server` — 服务器入口

**数据存储：**
- `data/versions/{version}/` — 按版本号存储资源文件
- `data/metadata.json` — 版本元数据（版本号、上传时间、文件列表）

**API 端点：**
- `GET /api/config` — 客户端获取服务器地址和当前版本
- `GET /api/version/latest` — 获取最新版本信息
- `GET /api/download/:version/*filepath` — 下载资源文件
- `POST /api/upload` — 上传资源（需要 API Key）
- `GET /api/versions` — 获取所有版本列表（需要 API Key）
- `DELETE /api/versions/:id` — 删除版本（需要 API Key）

**Web 界面：**
- 访问 `http://localhost:8080/` 打开管理界面
- 功能：查看版本列表、上传资源、删除版本
