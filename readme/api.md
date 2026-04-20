[← 返回总览](../README.md)

# API 文档

所有接口基础路径：`http://{server_host}:{port}`

---

## 公开接口（无需认证）

### GET /api/config

客户端启动时获取热更服务器地址和当前版本号。

**响应**
```json
{
  "patch_server_url": "http://localhost:8080",
  "current_version": 3
}
```

---

### GET /api/version/latest

获取最新版本详情。

**响应**
```json
{
  "version": 3,
  "upload_time": "2026-04-20T10:30:00Z",
  "total_size": 1048576,
  "file_count": 15
}
```

**错误响应**
```json
{"error": "no versions available"}  // 404
```

---

### GET /api/download/:version/*filepath

下载指定版本的资源文件。

**路径参数**
- `version` — 版本号（整数）
- `filepath` — 文件路径（支持子目录）

**示例**
```
GET /api/download/3/scenes/level1.unity3d
GET /api/download/3/ui/main.prefab
```

**响应**
- 成功：文件流（支持 Range 断点续传）
- 失败：`{"error": "version not found"}` 404

---

## 需要认证的接口

所有需要认证的接口必须在 Header 中携带 API Key：

```
X-API-Key: your-api-key
```

---

### POST /api/upload

上传资源文件，自动创建新版本。

**请求**
```
Content-Type: multipart/form-data

files: [文件1, 文件2, ...]   (必填，支持多文件)
version: 5                   (可选，不传则自动递增)
```

**响应**
```json
{
  "version": 4,
  "uploaded_files": ["scene1.unity3d", "ui/main.prefab"],
  "total_size": 2097152
}
```

**错误响应**
| 状态码 | 原因 |
|--------|------|
| 400 | 无文件或表单格式错误 |
| 401 | API Key 无效 |
| 413 | 文件超过大小限制 |
| 507 | 磁盘空间不足 |

---

### GET /api/versions

获取所有版本列表。

**响应**
```json
{
  "versions": [
    {
      "version": 1,
      "upload_time": "2026-04-20T10:00:00Z",
      "total_size": 1048576,
      "file_count": 5,
      "files": [
        {"path": "scene1.unity3d", "size": 524288}
      ]
    }
  ]
}
```

---

### GET /api/versions/:id

获取指定版本详情（含文件列表）。

**路径参数**
- `id` — 版本号

**响应**
```json
{
  "version": 3,
  "upload_time": "2026-04-20T10:30:00Z",
  "total_size": 1048576,
  "file_count": 2,
  "files": [
    {"path": "scene1.unity3d", "size": 524288},
    {"path": "ui/main.prefab", "size": 524288}
  ]
}
```

---

### DELETE /api/versions/:id

删除指定版本（同时删除版本目录下所有文件）。

**路径参数**
- `id` — 版本号

**响应**
```json
{"message": "version deleted"}
```

---

## Unity 客户端集成示例

```csharp
// 1. 获取配置
var configRes = await UnityWebRequest.Get(configServerUrl + "/api/config").SendWebRequest();
var config = JsonUtility.FromJson<ServerConfig>(configRes.downloadHandler.text);

// 2. 检查版本
int localVersion = PlayerPrefs.GetInt("patch_version", 0);
if (config.current_version > localVersion) {
    // 3. 下载资源
    string downloadUrl = $"{config.patch_server_url}/api/download/{config.current_version}/";
    // ... 下载并应用资源
    PlayerPrefs.SetInt("patch_version", config.current_version);
}
```

## Unity 编辑器上传示例

```csharp
// 上传打包好的资源
var form = new WWWForm();
foreach (var file in bundleFiles) {
    form.AddBinaryData("files", File.ReadAllBytes(file), Path.GetFileName(file));
}

var req = UnityWebRequest.Post(serverUrl + "/api/upload", form);
req.SetRequestHeader("X-API-Key", apiKey);
await req.SendWebRequest();
```

[← 返回总览](../README.md)
