# 测试说明

## 集成测试

测试所有 API 端点的上传和下载功能。

### 运行测试

**方式一：使用脚本（推荐）**
```bash
# 双击运行
test-integration.bat
```

**方式二：手动运行**
```bash
# 1. 启动服务器
go run cmd/server/main.go -config config.yaml

# 2. 在另一个终端运行测试
go run test/integration/main.go
```

### 测试覆盖

- ✅ GET /api/config — 获取配置
- ✅ POST /api/upload — 上传文件（需要 API Key）
- ✅ POST /api/upload — 无认证测试（应返回 401）
- ✅ GET /api/version/latest — 获取最新版本
- ✅ GET /api/versions — 获取版本列表
- ✅ GET /api/download/:version/:filepath — 下载文件
- ✅ GET /api/download/99999/notexist — 下载不存在的文件（应返回 404）
- ✅ DELETE /api/versions/:id — 删除版本

### 测试输出示例

```
=== Patch Server Test ===
Server: http://localhost:8080
Time:   2026-04-21 10:30:00

[Test] GET /api/config
  [PASS] request
  [PASS] status (status 200)
  [PASS] has patch_server_url: http://localhost:8080

[Test] POST /api/upload (no API key)
  [PASS] request
  [PASS] status 401 (status 401)

[Test] POST /api/upload
  [PASS] create form file scene1.unity3d
  [PASS] create form file main.prefab
  [PASS] request
  [PASS] status (status 200)
  [PASS] uploaded version: 1

[Test] GET /api/version/latest
  [PASS] request
  [PASS] status (status 200)
  [PASS] latest version: 1

[Test] GET /api/versions
  [PASS] request
  [PASS] status (status 200)
  [PASS] version count: 1

[Test] GET /api/download/1/scene1.unity3d
  [PASS] request
  [PASS] status (status 200)
  [PASS] downloaded 27 bytes

[Test] GET /api/download/99999/notexist.unity3d
  [PASS] request
  [PASS] status 404 (status 404)

[Test] DELETE /api/versions/1
  [PASS] request
  [PASS] status (status 200)

=== Results: 18 passed, 0 failed ===
```

### 注意事项

- 测试会创建临时文件并上传
- 测试结束后会删除上传的版本
- 确保 `config.yaml` 中的 API Key 是 `dev-api-key-change-in-production`
