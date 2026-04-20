# 资源热更服务器实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 实现 Unity 资源热更服务器，支持编辑器上传和客户端下载

**Architecture:** 单体 Go 服务器，使用 Gin 框架提供 REST API，文件系统存储资源，内嵌 Web 管理界面

**Tech Stack:** Go 1.21+, Gin, embed, YAML 配置

---

## File Structure

**New Files:**
- `cmd/server/main.go` — 服务器入口
- `internal/config/config.go` — 配置加载
- `internal/storage/manager.go` — 存储管理器
- `internal/storage/metadata.go` — 元数据结构
- `internal/auth/middleware.go` — API Key 认证中间件
- `internal/api/config.go` — 配置 API
- `internal/api/version.go` — 版本查询 API
- `internal/api/download.go` — 下载 API
- `internal/api/upload.go` — 上传 API
- `internal/web/handler.go` — Web 界面处理
- `web/index.html` — Web 管理界面
- `config.yaml` — 配置文件示例
- `go.mod` — Go 模块定义

**Test Files:**
- `internal/storage/manager_test.go`
- `internal/auth/middleware_test.go`
- `internal/api/upload_test.go`
- `internal/api/download_test.go`

---

### Task 1: 项目初始化和配置模块

**Files:**
- Create: `go.mod`
- Create: `config.yaml`
- Create: `internal/config/config.go`

- [ ] **Step 1: 初始化 Go 模块**

```bash
cd d:/xieliujian/spacetime_localpatchserver
go mod init spacetime_localpatchserver
go get github.com/gin-gonic/gin@latest
go get gopkg.in/yaml.v3@latest
```

Expected: `go.mod` 和 `go.sum` 创建成功

- [ ] **Step 2: 创建配置文件示例**

Create `config.yaml`:

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

- [ ] **Step 3: 编写配置加载模块**

Create `internal/config/config.go`:

```go
package config

import (
	"os"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Auth    AuthConfig    `yaml:"auth"`
	Storage StorageConfig `yaml:"storage"`
}

type ServerConfig struct {
	Port            int    `yaml:"port"`
	PatchServerURL  string `yaml:"patch_server_url"`
}

type AuthConfig struct {
	APIKey string `yaml:"api_key"`
}

type StorageConfig struct {
	DataDir         string `yaml:"data_dir"`
	MaxUploadSizeMB int    `yaml:"max_upload_size_mb"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	
	return &cfg, nil
}
```

- [ ] **Step 4: 提交**

```bash
git add go.mod go.sum config.yaml internal/config/config.go
git commit -m "feat: initialize project and add config module"
```

---

### Task 2: 存储管理器 - 元数据结构

**Files:**
- Create: `internal/storage/metadata.go`
- Create: `internal/storage/metadata_test.go`

- [ ] **Step 1: 编写元数据结构测试**

Create `internal/storage/metadata_test.go`:

```go
package storage

import (
	"encoding/json"
	"testing"
	"time"
)

func TestMetadata_Marshal(t *testing.T) {
	meta := &Metadata{
		Versions: []VersionInfo{
			{
				Version:    1,
				UploadTime: time.Date(2026, 4, 20, 10, 0, 0, 0, time.UTC),
				TotalSize:  1024,
				FileCount:  2,
				Files: []FileInfo{
					{Path: "scene1.unity3d", Size: 512},
					{Path: "ui/main.prefab", Size: 512},
				},
			},
		},
	}
	
	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	
	var decoded Metadata
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	
	if len(decoded.Versions) != 1 {
		t.Errorf("expected 1 version, got %d", len(decoded.Versions))
	}
	if decoded.Versions[0].Version != 1 {
		t.Errorf("expected version 1, got %d", decoded.Versions[0].Version)
	}
}
```

- [ ] **Step 2: 运行测试确认失败**

```bash
go test ./internal/storage -v
```

Expected: FAIL - undefined: Metadata

- [ ] **Step 3: 实现元数据结构**

Create `internal/storage/metadata.go`:

```go
package storage

import (
	"encoding/json"
	"os"
	"time"
)

type Metadata struct {
	Versions []VersionInfo `json:"versions"`
}

type VersionInfo struct {
	Version    int        `json:"version"`
	UploadTime time.Time  `json:"upload_time"`
	TotalSize  int64      `json:"total_size"`
	FileCount  int        `json:"file_count"`
	Files      []FileInfo `json:"files"`
}

type FileInfo struct {
	Path string `json:"path"`
	Size int64  `json:"size"`
}

func LoadMetadata(path string) (*Metadata, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Metadata{Versions: []VersionInfo{}}, nil
		}
		return nil, err
	}
	
	var meta Metadata
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}
	
	return &meta, nil
}

func (m *Metadata) Save(path string) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(path, data, 0644)
}
```

- [ ] **Step 4: 运行测试确认通过**

```bash
go test ./internal/storage -v
```

Expected: PASS

- [ ] **Step 5: 提交**

```bash
git add internal/storage/
git commit -m "feat: add metadata structure and persistence"
```

---

### Task 3: 存储管理器 - 核心逻辑

**Files:**
- Create: `internal/storage/manager.go`
- Create: `internal/storage/manager_test.go`

- [ ] **Step 1: 编写存储管理器测试**

Create `internal/storage/manager_test.go`:

```go
package storage

import (
	"os"
	"path/filepath"
	"testing"
)

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
	
	v2 := mgr.NextVersion()
	if v2 != 2 {
		t.Errorf("expected version 2, got %d", v2)
	}
}

func TestManager_GetVersionPath(t *testing.T){
	tmpDir := t.TempDir()
	mgr, _ := NewManager(tmpDir)
	
	path := mgr.GetVersionPath(1)
	expected := filepath.Join(tmpDir, "versions", "1")
	if path != expected {
		t.Errorf("expected %s, got %s", expected, path)
	}
}
```

- [ ] **Step 2: 运行测试确认失败**

```bash
go test ./internal/storage -v -run TestManager
```

Expected: FAIL - undefined: NewManager

- [ ] **Step 3: 实现存储管理器**

Create `internal/storage/manager.go`:

```go
package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Manager struct {
	dataDir      string
	metadataPath string
	metadata     *Metadata
	mu           sync.Mutex
}

func NewManager(dataDir string) (*Manager, error) {
	versionsDir := filepath.Join(dataDir, "versions")
	if err := os.MkdirAll(versionsDir, 0755); err != nil {
		return nil, err
	}
	
	metadataPath := filepath.Join(dataDir, "metadata.json")
	metadata, err := LoadMetadata(metadataPath)
	if err != nil {
		return nil, err
	}
	
	return &Manager{
		dataDir:      dataDir,
		metadataPath: metadataPath,
		metadata:     metadata,
	}, nil
}

func (m *Manager) NextVersion() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	maxVer := 0
	for _, v := range m.metadata.Versions {
		if v.Version > maxVer {
			maxVer = v.Version
		}
	}
	return maxVer + 1
}

func (m *Manager) GetVersionPath(version int) string {
	return filepath.Join(m.dataDir, "versions", fmt.Sprintf("%d", version))
}

func (m *Manager) GetLatestVersion() *VersionInfo {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if len(m.metadata.Versions) == 0 {
		return nil
	}
	
	latest := &m.metadata.Versions[0]
	for i := range m.metadata.Versions {
		if m.metadata.Versions[i].Version > latest.Version {
			latest = &m.metadata.Versions[i]
		}
	}
	return latest
}

func (m *Manager) AddVersion(info VersionInfo) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.metadata.Versions = append(m.metadata.Versions, info)
	return m.metadata.Save(m.metadataPath)
}

func (m *Manager) DeleteVersion(version int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	for i, v := range m.metadata.Versions {
		if v.Version == version {
			m.metadata.Versions = append(m.metadata.Versions[:i], m.metadata.Versions[i+1:]...)
			if err := m.metadata.Save(m.metadataPath); err != nil {
				return err
			}
			return os.RemoveAll(m.GetVersionPath(version))
		}
	}
	return fmt.Errorf("version %d not found", version)
}

func (m *Manager) GetAllVersions() []VersionInfo {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	return append([]VersionInfo{}, m.metadata.Versions...)
}
```

- [ ] **Step 4: 运行测试确认通过**

```bash
go test ./internal/storage -v
```

Expected: PASS

- [ ] **Step 5: 提交**

```bash
git add internal/storage/
git commit -m "feat: implement storage manager with version management"
```

---

### Task 4: 认证中间件

**Files:**
- Create: `internal/auth/middleware.go`
- Create: `internal/auth/middleware_test.go`

- [ ] **Step 1: 编写认证中间件测试**

Create `internal/auth/middleware_test.go`:

```go
package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	
	"github.com/gin-gonic/gin"
)

func TestAPIKeyMiddleware_Valid(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	r := gin.New()
	r.Use(APIKeyMiddleware("test-key"))
	r.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-API-Key", "test-key")
	w := httptest.NewRecorder()
	
	r.ServeHTTP(w, req)
	
	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestAPIKeyMiddleware_Invalid(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	r := gin.New()
	r.Use(APIKeyMiddleware("test-key"))
	r.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-API-Key", "wrong-key")
	w := httptest.NewRecorder()
	
	r.ServeHTTP(w, req)
	
	if w.Code != 401 {
		t.Errorf("expected 401, got %d", w.Code)
	}
}
```

- [ ] **Step 2: 运行测试确认失败**

```bash
go test ./internal/auth -v
```

Expected: FAIL - undefined: APIKeyMiddleware

- [ ] **Step 3: 实现认证中间件**

Create `internal/auth/middleware.go`:

```go
package auth

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
)

func APIKeyMiddleware(expectedKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey != expectedKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid API key"})
			c.Abort()
			return
		}
		c.Next()
	}
}
```

- [ ] **Step 4: 运行测试确认通过**

```bash
go test ./internal/auth -v
```

Expected: PASS

- [ ] **Step 5: 提交**

```bash
git add internal/auth/
git commit -m "feat: add API key authentication middleware"
```

---

### Task 5: 配置和版本查询 API

**Files:**
- Create: `internal/api/config.go`
- Create: `internal/api/version.go`

- [ ] **Step 1: 实现配置 API**

Create `internal/api/config.go`:

```go
package api

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
	"spacetime_localpatchserver/internal/storage"
)

type ConfigHandler struct {
	patchServerURL string
	storage        *storage.Manager
}

func NewConfigHandler(patchServerURL string, storage *storage.Manager) *ConfigHandler {
	return &ConfigHandler{
		patchServerURL: patchServerURL,
		storage:        storage,
	}
}

func (h *ConfigHandler) GetConfig(c *gin.Context) {
	latest := h.storage.GetLatestVersion()
	currentVersion := 0
	if latest != nil {
		currentVersion = latest.Version
	}
	
	c.JSON(http.StatusOK, gin.H{
		"patch_server_url": h.patchServerURL,
		"current_version":  currentVersion,
	})
}
```

- [ ] **Step 2: 实现版本查询 API**

Create `internal/api/version.go`:

```go
package api

import (
	"net/http"
	"strconv"
	
	"github.com/gin-gonic/gin"
	"spacetime_localpatchserver/internal/storage"
)

type VersionHandler struct {
	storage *storage.Manager
}

func NewVersionHandler(storage *storage.Manager) *VersionHandler {
	return &VersionHandler{storage: storage}
}

func (h *VersionHandler) GetLatest(c *gin.Context) {
	latest := h.storage.GetLatestVersion()
	if latest == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no versions available"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"version":     latest.Version,
		"upload_time": latest.UploadTime,
		"total_size":  latest.TotalSize,
		"file_count":  latest.FileCount,
	})
}

func (h *VersionHandler) GetAll(c *gin.Context) {
	versions := h.storage.GetAllVersions()
	c.JSON(http.StatusOK, gin.H{"versions": versions})
}

func (h *VersionHandler) GetDetail(c *gin.Context) {
	versionStr := c.Param("id")
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid version"})
		return
	}
	
	versions := h.storage.GetAllVersions()
	for _, v := range versions {
		if v.Version == version {
			c.JSON(http.StatusOK, v)
			return
		}
	}
	
	c.JSON(http.StatusNotFound, gin.H{"error": "version not found"})
}

func (h *VersionHandler) Delete(c *gin.Context) {
	versionStr := c.Param("id")
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid version"})
		return
	}
	
	if err := h.storage.DeleteVersion(version); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "version deleted"})
}
```

- [ ] **Step 3: 提交**

```bash
git add internal/api/
git commit -m "feat: add config and version query APIs"
```

---

### Task 6: 下载 API

**Files:**
- Create: `internal/api/download.go`

- [ ] **Step 1: 实现下载 API**

Create `internal/api/download.go`:

```go
package api

import (
	"net/http"
	"path/filepath"
	"strconv"
	
	"github.com/gin-gonic/gin"
	"spacetime_localpatchserver/internal/storage"
)

type DownloadHandler struct {
	storage *storage.Manager
}

func NewDownloadHandler(storage *storage.Manager) *DownloadHandler {
	return &DownloadHandler{storage: storage}
}

func (h *DownloadHandler) Download(c *gin.Context) {
	versionStr := c.Param("version")
	filepath_param := c.Param("filepath")
	
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid version"})
		return
	}
	
	versionPath := h.storage.GetVersionPath(version)
	fullPath := filepath.Join(versionPath, filepath_param)
	
	c.File(fullPath)
}
```

- [ ] **Step 2: 提交**

```bash
git add internal/api/download.go
git commit -m "feat: add download API with file serving"
```

---

### Task 7: 上传 API

**Files:**
- Create: `internal/api/upload.go`

- [ ] **Step 1: 实现上传 API**

Create `internal/api/upload.go`:

```go
package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	
	"github.com/gin-gonic/gin"
	"spacetime_localpatchserver/internal/storage"
)

type UploadHandler struct {
	storage         *storage.Manager
	maxUploadSizeMB int
}

func NewUploadHandler(storage *storage.Manager, maxUploadSizeMB int) *UploadHandler {
	return &UploadHandler{
		storage:         storage,
		maxUploadSizeMB: maxUploadSizeMB,
	}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form data"})
		return
	}
	
	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no files uploaded"})
		return
	}
	
	version := h.storage.NextVersion()
	if versionStr := c.PostForm("version"); versionStr != "" {
		if v, err := strconv.Atoi(versionStr); err == nil {
			version = v
		}
	}
	
	versionPath := h.storage.GetVersionPath(version)
	if err := os.MkdirAll(versionPath, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create version directory"})
		return
	}
	
	var uploadedFiles []string
	var totalSize int64
	var fileInfos []storage.FileInfo
	
	for _, file := range files {
		dst := filepath.Join(versionPath, file.Filename)
		dstDir := filepath.Dir(dst)
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			os.RemoveAll(versionPath)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create directory"})
			return
		}
		
		if err := c.SaveUploadedFile(file, dst); err != nil {
			os.RemoveAll(versionPath)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
			return
		}
		
		uploadedFiles = append(uploadedFiles, file.Filename)
		totalSize += file.Size
		fileInfos = append(fileInfos, storage.FileInfo{
			Path: file.Filename,
			Size: file.Size,
		})
	}
	
	versionInfo := storage.VersionInfo{
		Version:    version,
		UploadTime: time.Now(),
		TotalSize:  totalSize,
		FileCount:  len(files),
		Files:      fileInfos,
	}
	
	if err := h.storage.AddVersion(versionInfo); err != nil {
		os.RemoveAll(versionPath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save metadata"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"version":        version,
		"uploaded_files": uploadedFiles,
		"total_size":     totalSize,
	})
}
```

- [ ] **Step 2: 提交**

```bash
git add internal/api/upload.go
git commit -m "feat: add upload API with multipart form support"
```

---

### Task 8: Web 管理界面

**Files:**
- Create: `web/index.html`
- Create: `internal/web/handler.go`

- [ ] **Step 1: 创建 Web 界面 HTML**

Create `web/index.html`:

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>资源热更服务器管理</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: Arial, sans-serif; padding: 20px; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background: white; padding: 20px; border-radius: 8px; }
        h1 { margin-bottom: 20px; color: #333; }
        .upload-section { margin-bottom: 30px; padding: 20px; background: #f9f9f9; border-radius: 4px; }
        .upload-section input[type="text"] { width: 300px; padding: 8px; margin-right: 10px; }
        .upload-section input[type="file"] { margin: 10px 0; }
        .upload-section button { padding: 10px 20px; background: #007bff; color: white; border: none; border-radius: 4px; cursor: pointer; }
        .upload-section button:hover { background: #0056b3; }
        table { width: 100%; border-collapse: collapse; margin-top: 20px; }
        th, td { padding: 12px; text-align: left; border-bottom: 1px solid #ddd; }
        th { background: #f0f0f0; font-weight: bold; }
        .delete-btn { padding: 6px 12px; background: #dc3545; color: white; border: none; border-radius: 4px; cursor: pointer; }
        .delete-btn:hover { background: #c82333; }
        .status { margin-top: 10px; padding: 10px; border-radius: 4px; }
        .status.success { background: #d4edda; color: #155724; }
        .status.error { background: #f8d7da; color: #721c24; }
    </style>
</head>
<body>
    <div class="container">
        <h1>资源热更服务器管理</h1>
        
        <div class="upload-section">
            <h2>上传资源</h2>
            <input type="text" id="apiKey" placeholder="API Key" />
            <br>
            <input type="file" id="fileInput" multiple />
            <br>
            <button onclick="uploadFiles()">上传</button>
            <div id="uploadStatus" class="status" style="display:none;"></div>
        </div>
        
        <h2>版本列表</h2>
        <table id="versionTable">
            <thead>
                <tr>
                    <th>版本号</th>
                    <th>上传时间</th>
                    <th>文件数量</th>
                    <th>总大小</th>
                    <th>操作</th>
                </tr>
            </thead>
            <tbody></tbody>
        </table>
    </div>
    
    <script>
        function formatBytes(bytes) {
            if (bytes === 0) return '0 B';
            const k = 1024;
            const sizes = ['B', 'KB', 'MB', 'GB'];
            const i = Math.floor(Math.log(bytes) / Math.log(k));
            return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
        }
        
        function formatDate(dateStr) {
            const date = new Date(dateStr);
            return date.toLocaleString('zh-CN');
        }
        
        async function loadVersions() {
            try {
                const res = await fetch('/api/versions');
                const data = await res.json();
                const tbody = document.querySelector('#versionTable tbody');
                tbody.innerHTML = '';
                
                data.versions.forEach(v => {
                    const row = tbody.insertRow();
                    row.innerHTML = `
                        <td>${v.version}</td>
                        <td>${formatDate(v.upload_time)}</td>
                        <td>${v.file_count}</td>
                        <td>${formatBytes(v.total_size)}</td>
                        <td><button class="delete-btn" onclick="deleteVersion(${v.version})">删除</button></td>
                    `;
                });
            } catch (err) {
                console.error('Failed to load versions:', err);
            }
        }
        
        async function uploadFiles() {
            const apiKey = document.getElementById('apiKey').value;
            const fileInput = document.getElementById('fileInput');
            const status = document.getElementById('uploadStatus');
            
            if (!apiKey) {
                status.className = 'status error';
                status.textContent = '请输入 API Key';
                status.style.display = 'block';
                return;
            }
            
            if (fileInput.files.length === 0) {
                status.className = 'status error';
                status.textContent = '请选择文件';
                status.style.display = 'block';
                return;
            }
            
            const formData = new FormData();
            for (let file of fileInput.files) {
                formData.append('files', file);
            }
            
            try {
                const res = await fetch('/api/upload', {
                    method: 'POST',
                    headers: { 'X-API-Key': apiKey },
                    body: formData
                });
                
                const data = await res.json();
                
                if (res.ok) {
                    status.className = 'status success';
                    status.textContent = `上传成功！版本号: ${data.version}`;
                    fileInput.value = '';
                    loadVersions();
                } else {
                    status.className = 'status error';
                    status.textContent = `上传失败: ${data.error}`;
                }
            } catch (err) {
                status.className = 'status error';
                status.textContent = `上传失败: ${err.message}`;
            }
            
            status.style.display = 'block';
        }
        
        async function deleteVersion(version) {
            const apiKey = document.getElementById('apiKey').value;
            
            if (!apiKey) {
                alert('请输入 API Key');
                return;
            }
            
            if (!confirm(`确定删除版本 ${version}？`)) {
                return;
            }
            
            try {
                const res = await fetch(`/api/versions/${version}`, {
                    method: 'DELETE',
                    headers: { 'X-API-Key': apiKey }
                });
                
                if (res.ok) {
                    alert('删除成功');
                    loadVersions();
                } else {
                    const data = await res.json();
                    alert(`删除失败: ${data.error}`);
                }
            } catch (err) {
                alert(`删除失败: ${err.message}`);
            }
        }
        
        loadVersions();
    </script>
</body>
</html>
```

- [ ] **Step 2: 创建 Web 处理器**

Create `internal/web/handler.go`:

```go
package web

import (
	"embed"
	"io/fs"
	"net/http"
	
	"github.com/gin-gonic/gin"
)

//go:embed ../../../web
var webFS embed.FS

func SetupRoutes(r *gin.Engine) error {
	webContent, err := fs.Sub(webFS, "web")
	if err != nil {
		return err
	}
	
	r.GET("/", func(c *gin.Context) {
		c.FileFromFS("index.html", http.FS(webContent))
	})
	
	return nil
}
```

- [ ] **Step 3: 提交**

```bash
git add web/ internal/web/
git commit -m "feat: add web management interface"
```

---

### Task 9: 主服务器入口

**Files:**
- Create: `cmd/server/main.go`

- [ ] **Step 1: 实现服务器入口**

Create `cmd/server/main.go`:

```go
package main

import (
	"flag"
	"fmt"
	"log"
	
	"github.com/gin-gonic/gin"
	"spacetime_localpatchserver/internal/api"
	"spacetime_localpatchserver/internal/auth"
	"spacetime_localpatchserver/internal/config"
	"spacetime_localpatchserver/internal/storage"
	"spacetime_localpatchserver/internal/web"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()
	
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	
	storageMgr, err := storage.NewManager(cfg.Storage.DataDir)
	if err != nil {
		log.Fatalf("Failed to create storage manager: %v", err)
	}
	
	r := gin.Default()
	
	configHandler := api.NewConfigHandler(cfg.Server.PatchServerURL, storageMgr)
	versionHandler := api.NewVersionHandler(storageMgr)
	downloadHandler := api.NewDownloadHandler(storageMgr)
	uploadHandler := api.NewUploadHandler(storageMgr, cfg.Storage.MaxUploadSizeMB)
	
	r.GET("/api/config", configHandler.GetConfig)
	r.GET("/api/version/latest", versionHandler.GetLatest)
	r.GET("/api/download/:version/*filepath", downloadHandler.Download)
	
	authGroup := r.Group("/api")
	authGroup.Use(auth.APIKeyMiddleware(cfg.Auth.APIKey))
	{
		authGroup.POST("/upload", uploadHandler.Upload)
		authGroup.GET("/versions", versionHandler.GetAll)
		authGroup.GET("/versions/:id", versionHandler.GetDetail)
		authGroup.DELETE("/versions/:id", versionHandler.Delete)
	}
	
	if err := web.SetupRoutes(r); err != nil {
		log.Fatalf("Failed to setup web routes: %v", err)
	}
	
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
```

- [ ] **Step 2: 测试服务器启动**

```bash
go run cmd/server/main.go -config config.yaml
```

Expected: Server starting on :8080

- [ ] **Step 3: 测试配置 API**

在另一个终端：

```bash
curl http://localhost:8080/api/config
```

Expected: `{"current_version":0,"patch_server_url":"http://localhost:8080"}`

- [ ] **Step 4: 停止服务器并提交**

```bash
git add cmd/server/main.go
git commit -m "feat: add main server entry point with all routes"
```

---

### Task 10: 集成测试和文档

**Files:**
- Modify: `CLAUDE.md`
- Create: `README.md`

- [ ] **Step 1: 更新 CLAUDE.md**

Update `CLAUDE.md`:

```markdown
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
- `internal/web` — Web 管理界面（embed）
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
```

- [ ] **Step 2: 创建 README**

Create `README.md`:

```markdown
# spacetime_localpatchserver

Unity 资源热更服务器

## 功能

- Unity 编辑器通过 HTTP API 上传资源
- 游戏客户端通过远程配置获取服务器地址
- 版本号管理（整数递增：1, 2, 3...）
- Web 管理界面

## 快速开始

1. 配置 `config.yaml`
2. 运行服务器：`go run cmd/server/main.go`
3. 访问 `http://localhost:8080` 打开管理界面

## API 文档

详见设计文档：`docs/superpowers/specs/2026-04-20-patch-server-design.md`
```

- [ ] **Step 3: 集成测试**

手动测试完整流程：

1. 启动服务器
2. 打开 Web 界面 `http://localhost:8080`
3. 输入 API Key: `dev-api-key-change-in-production`
4. 上传测试文件
5. 验证版本列表显示
6. 测试下载：`curl http://localhost:8080/api/download/1/test.txt`
7. 测试删除版本

- [ ] **Step 4: 提交**

```bash
git add CLAUDE.md README.md
git commit -m "docs: update CLAUDE.md and add README"
```

---

### Task 11: 最终验证和构建

- [ ] **Step 1: 运行所有测试**

```bash
go test ./... -v
```

Expected: 所有测试通过

- [ ] **Step 2: 构建生产二进制**

```bash
go build -o patchserver cmd/server/main.go
```

Expected: `patchserver` 或 `patchserver.exe` 生成

- [ ] **Step 3: 测试生产二进制**

```bash
./patchserver -config config.yaml
```

Expected: 服务器正常启动

- [ ] **Step 4: 最终提交**

```bash
git add .
git commit -m "chore: final build and verification"
```

---

## 完成

所有任务完成后，服务器具备以下功能：

✅ 配置加载（YAML）  
✅ 版本管理（整数递增，文件系统存储）  
✅ API Key 认证  
✅ 配置 API（客户端获取服务器地址）  
✅ 版本查询 API  
✅ 上传 API（multipart/form-data）  
✅ 下载 API（文件流 + 断点续传支持）  
✅ Web 管理界面  
✅ 单元测试覆盖核心模块  

**下一步：** Unity 客户端和编辑器集成（需要在 UnityFramework 项目中实现）
