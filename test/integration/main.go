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

func pause(message string) {
	fmt.Printf("\n%s\n", message)
	fmt.Print("按回车继续...")
	fmt.Scanln()
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
	file3 := filepath.Join(tmpDir, "audio", "bgm.mp3")
	os.MkdirAll(filepath.Dir(file2), 0755)
	os.MkdirAll(filepath.Dir(file3), 0755)
	os.WriteFile(file1, []byte("This is a fake Unity scene bundle for testing"), 0644)
	os.WriteFile(file2, []byte("This is a fake UI prefab for testing"), 0644)
	os.WriteFile(file3, []byte("This is a fake audio file for testing"), 0644)

	// 构建 multipart form
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	for _, path := range []string{file1, file2, file3} {
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
		fmt.Printf("  [INFO] uploaded files: %v\n", result["uploaded_files"])
		fmt.Printf("  [INFO] total size: %v bytes\n", result["total_size"])
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
		fmt.Printf("  [INFO] upload_time: %v\n", result["upload_time"])
		fmt.Printf("  [INFO] file_count: %v\n", result["file_count"])
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
		fmt.Printf("  [INFO] content: %s\n", string(body))
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
	fmt.Println("=== Patch Server Integration Test ===")
	fmt.Printf("Server: %s\n", serverURL)
	fmt.Printf("Time:   %s\n", time.Now().Format("2006-01-02 15:04:05"))

	// 检查服务器是否在线
	_, err := http.Get(serverURL + "/api/config")
	if err != nil {
		fmt.Printf("\n[ERROR] Server not reachable at %s\n", serverURL)
		fmt.Println("Please start the server first: go run cmd/server/main.go -config config.yaml")
		pause("按回车退出...")
		os.Exit(1)
	}

	// 基础测试
	testGetConfig()
	testUploadUnauthorized()

	// 上传测试
	version := testUpload()
	testGetLatestVersion()
	testGetVersions()

	// 暂停，让用户查看 web 界面
	fmt.Println("\n" + "============================================================")
	fmt.Printf("✅ 上传成功！版本号: %d\n", version)
	fmt.Printf("📂 数据目录: ./data/versions/%d/\n", version)
	fmt.Printf("🌐 Web 界面: %s\n", serverURL)
	fmt.Println("============================================================")
	pause("现在可以打开浏览器访问 Web 界面查看上传的版本。")

	// 下载测试
	testDownload(version)
	testDownloadNotFound()

	// 询问是否删除
	fmt.Println("\n" + "============================================================")
	fmt.Printf("⚠️  即将删除版本 %d\n", version)
	fmt.Println("============================================================")
	fmt.Print("是否删除此版本？(y/n): ")
	var answer string
	fmt.Scanln(&answer)

	if answer == "y" || answer == "Y" {
		testDeleteVersion(version)
		fmt.Println("\n✅ 版本已删除，可以刷新 Web 界面确认")
	} else {
		fmt.Println("\n⏭️  跳过删除，版本保留在服务器上")
	}

	fmt.Printf("\n=== Results: %d passed, %d failed ===\n", pass, fail)
	pause("\n测试完成，按回车退出...")

	if fail > 0 {
		os.Exit(1)
	}
}
