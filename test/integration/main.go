package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	serverURL = "http://localhost:8080"
	apiKey    = "dev-api-key-change-in-production"
)

var pass, fail int

func check(name string, err error) bool {
	if err != nil {
		fmt.Printf("  [FAIL] %s: %v\n", name, err)
		fail++
		return false
	}
	fmt.Printf("  [PASS] %s\n", name)
	pass++
	return true
}

func checkStatus(name string, got, want int) bool {
	if got != want {
		fmt.Printf("  [FAIL] %s: status %d, want %d\n", name, got, want)
		fail++
		return false
	}
	fmt.Printf("  [PASS] %s (status %d)\n", name, got)
	pass++
	return true
}

// --- 测试：获取配置 ---
func testGetConfig() {
	fmt.Println("\n[Test] GET /api/config")
	resp, err := http.Get(serverURL + "/api/config")
	if !check("request", err) {
		return
	}
	defer resp.Body.Close()
	checkStatus("status", resp.StatusCode, 200)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	if _, ok := result["patch_server_url"]; ok {
		fmt.Printf("  [PASS] has patch_server_url: %v\n", result["patch_server_url"])
		pass++
	} else {
		fmt.Println("  [FAIL] missing patch_server_url")
		fail++
	}
}

// --- 测试：上传文件 ---
func testUpload() int {
	fmt.Println("\n[Test] POST /api/upload")

	// 创建临时测试文件
	tmpDir, _ := os.MkdirTemp("", "patch-test-*")
	defer os.RemoveAll(tmpDir)

	file1 := filepath.Join(tmpDir, "scene1.unity3d")
	file2 := filepath.Join(tmpDir, "ui", "main.prefab")
	os.MkdirAll(filepath.Dir(file2), 0755)
	os.WriteFile(file1, []byte("fake unity bundle content 1"), 0644)
	os.WriteFile(file2, []byte("fake prefab content 2"), 0644)

	// 构建 multipart form
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	for _, path := range []string{file1, file2} {
		fw, err := w.CreateFormFile("files", filepath.Base(path))
		if !check("create form file "+filepath.Base(path), err) {
			return 0
		}
		f, _ := os.Open(path)
		io.Copy(fw, f)
		f.Close()
	}
	w.Close()

	req, _ := http.NewRequest("POST", serverURL+"/api/upload", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("X-API-Key", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if !check("request", err) {
		return 0
	}
	defer resp.Body.Close()
	checkStatus("status", resp.StatusCode, 200)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	version := 0
	if v, ok := result["version"]; ok {
		version = int(v.(float64))
		fmt.Printf("  [PASS] uploaded version: %d\n", version)
		pass++
	} else {
		fmt.Println("  [FAIL] missing version in response")
		fail++
	}
	return version
}

// --- 测试：上传需要认证 ---
func testUploadUnauthorized() {
	fmt.Println("\n[Test] POST /api/upload (no API key)")

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, _ := w.CreateFormFile("files", "test.txt")
	fw.Write([]byte("test"))
	w.Close()

	req, _ := http.NewRequest("POST", serverURL+"/api/upload", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if !check("request", err) {
		return
	}
	defer resp.Body.Close()
	checkStatus("status 401", resp.StatusCode, 401)
}

// --- 测试：获取最新版本 ---
func testGetLatestVersion() {
	fmt.Println("\n[Test] GET /api/version/latest")
	resp, err := http.Get(serverURL + "/api/version/latest")
	if !check("request", err) {
		return
	}
	defer resp.Body.Close()
	checkStatus("status", resp.StatusCode, 200)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	if v, ok := result["version"]; ok {
		fmt.Printf("  [PASS] latest version: %v\n", v)
		pass++
	} else {
		fmt.Println("  [FAIL] missing version")
		fail++
	}
}

// --- 测试：获取版本列表 ---
func testGetVersions() {
	fmt.Println("\n[Test] GET /api/versions")
	req, _ := http.NewRequest("GET", serverURL+"/api/versions", nil)
	req.Header.Set("X-API-Key", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if !check("request", err) {
		return
	}
	defer resp.Body.Close()
	checkStatus("status", resp.StatusCode, 200)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	if versions, ok := result["versions"]; ok {
		list := versions.([]interface{})
		fmt.Printf("  [PASS] version count: %d\n", len(list))
		pass++
	} else {
		fmt.Println("  [FAIL] missing versions")
		fail++
	}
}

// --- 测试：下载文件 ---
func testDownload(version int) {
	fmt.Printf("\n[Test] GET /api/download/%d/scene1.unity3d\n", version)
	if version == 0 {
		fmt.Println("  [SKIP] no version to download")
		return
	}

	url := fmt.Sprintf("%s/api/download/%d/scene1.unity3d", serverURL, version)
	resp, err := http.Get(url)
	if !check("request", err) {
		return
	}
	defer resp.Body.Close()
	checkStatus("status", resp.StatusCode, 200)

	body, _ := io.ReadAll(resp.Body)
	if len(body) > 0 {
		fmt.Printf("  [PASS] downloaded %d bytes\n", len(body))
		pass++
	} else {
		fmt.Println("  [FAIL] empty response body")
		fail++
	}
}

// --- 测试：下载不存在的版本 ---
func testDownloadNotFound() {
	fmt.Println("\n[Test] GET /api/download/99999/notexist.unity3d")
	resp, err := http.Get(serverURL + "/api/download/99999/notexist.unity3d")
	if !check("request", err) {
		return
	}
	defer resp.Body.Close()
	checkStatus("status 404", resp.StatusCode, 404)
}

// --- 测试：删除版本 ---
func testDeleteVersion(version int) {
	fmt.Printf("\n[Test] DELETE /api/versions/%d\n", version)
	if version == 0 {
		fmt.Println("  [SKIP] no version to delete")
		return
	}

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/api/versions/%d", serverURL, version), nil)
	req.Header.Set("X-API-Key", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if !check("request", err) {
		return
	}
	defer resp.Body.Close()
	checkStatus("status", resp.StatusCode, 200)
}

func main() {
	fmt.Println("=== Patch Server Test ===")
	fmt.Printf("Server: %s\n", serverURL)
	fmt.Printf("Time:   %s\n", time.Now().Format("2006-01-02 15:04:05"))

	// 检查服务器是否在线
	_, err := http.Get(serverURL + "/api/config")
	if err != nil {
		fmt.Printf("\n[ERROR] Server not reachable at %s\n", serverURL)
		fmt.Println("Please start the server first: go run cmd/server/main.go -config config.yaml")
		os.Exit(1)
	}

	testGetConfig()
	testUploadUnauthorized()
	version := testUpload()
	testGetLatestVersion()
	testGetVersions()
	testDownload(version)
	testDownloadNotFound()
	testDeleteVersion(version)

	fmt.Printf("\n=== Results: %d passed, %d failed ===\n", pass, fail)
	if fail > 0 {
		os.Exit(1)
	}
}
