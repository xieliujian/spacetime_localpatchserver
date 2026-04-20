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

	versionInfo := h.storage.GetVersion(version)
	if versionInfo == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "version not found"})
		return
	}

	c.JSON(http.StatusOK, versionInfo)
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
