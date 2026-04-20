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
