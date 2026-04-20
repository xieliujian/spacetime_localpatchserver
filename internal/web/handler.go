package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"spacetime_localpatchserver/internal/storage"
)

type Handler struct {
	storage *storage.Manager
}

func NewHandler(storage *storage.Manager) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Patch Server Management",
	})
}

func (h *Handler) GetVersions(c *gin.Context) {
	versions := h.storage.GetAllVersions()
	c.JSON(http.StatusOK, gin.H{"versions": versions})
}
