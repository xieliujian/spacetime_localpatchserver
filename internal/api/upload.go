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
