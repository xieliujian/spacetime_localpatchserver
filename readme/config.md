[← 返回总览](../README.md)

# 配置说明

配置文件：`config.yaml`

---

## 完整配置示例

```yaml
server:
  port: 8080
  patch_server_url: "http://localhost:8080"

auth:
  api_key: "dev-api-key-change-in-production"

storage:
  data_dir: "./data"
  max_upload_size_mb: 500
```

---

## 配置项说明

### server

服务器相关配置。

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `port` | int | 是 | HTTP 服务端口（1-65535） |
| `patch_server_url` | string | 是 | 返回给客户端的服务器地址，支持配置不同环境 |

**示例场景：**
- 开发环境：`http://localhost:8080`
- 测试环境：`http://test-patch.example.com`
- 生产环境：`http://patch.example.com`

---

### auth

认证相关配置。

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `api_key` | string | 是 | API Key，用于上传和版本管理接口认证 |

**安全建议：**
- 生产环境必须修改默认值
- 使用强密码（建议 32 位随机字符串）
- 定期轮换 API Key
- 不要将 API Key 提交到版本控制

---

### storage

存储相关配置。

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `data_dir` | string | 是 | 资源存储根目录（相对或绝对路径） |
| `max_upload_size_mb` | int | 是 | 单次上传最大文件大小（MB） |

**存储目录结构：**
```
{data_dir}/
├── metadata.json
└── versions/
    ├── 1/
    ├── 2/
    └── 3/
```

**注意事项：**
- 确保目录有读写权限
- 定期清理旧版本释放空间
- 建议使用绝对路径避免路径问题

---

## 环境变量覆盖（可选）

支持通过环境变量覆盖配置：

```bash
export PATCH_SERVER_PORT=9090
export PATCH_SERVER_API_KEY="production-key"
go run cmd/server/main.go
```

---

## 配置验证

启动时会自动验证配置：
- 端口范围检查
- 必填字段检查
- 目录权限检查

验证失败会输出错误并退出。

[← 返回总览](../README.md)
