package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"spacetime_localpatchserver/internal/api"
	"spacetime_localpatchserver/internal/auth"
	"spacetime_localpatchserver/internal/config"
	"spacetime_localpatchserver/internal/storage"
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

	// 静态文件和页面
	r.StaticFS("/web", http.Dir("./web"))
	r.GET("/", func(c *gin.Context) {
		c.File("./web/index.html")
	})
	r.GET("/version.html", func(c *gin.Context) {
		c.File("./web/version.html")
	})

	configHandler := api.NewConfigHandler(cfg.Server.PatchServerURL, storageMgr)
	versionHandler := api.NewVersionHandler(storageMgr)
	downloadHandler := api.NewDownloadHandler(storageMgr)
	uploadHandler := api.NewUploadHandler(storageMgr, cfg.Storage.MaxUploadSizeMB)

	// 公开接口（无需认证）
	r.GET("/api/config", configHandler.GetConfig)
	r.GET("/api/version/latest", versionHandler.GetLatest)
	r.GET("/api/versions", versionHandler.GetAll)
	r.GET("/api/versions/:id", versionHandler.GetDetail)
	r.GET("/api/download/:version/*filepath", downloadHandler.Download)

	// 需要认证的接口
	authGroup := r.Group("/api")
	authGroup.Use(auth.APIKeyMiddleware(cfg.Auth.APIKey))
	{
		authGroup.POST("/upload", uploadHandler.Upload)
		authGroup.DELETE("/versions/:id", versionHandler.Delete)
	}

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
