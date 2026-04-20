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
	filePath := c.Param("filepath")

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid version"})
		return
	}

	if h.storage.GetVersion(version) == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "version not found"})
		return
	}

	versionPath := h.storage.GetVersionPath(version)
	fullPath := filepath.Join(versionPath, filePath)

	c.File(fullPath)
}
