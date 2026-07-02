package main

import (
	"fmt"
	"kira-go/internal/config"
	"kira-go/internal/database"
	"kira-go/internal/router"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化加载配置
	config.Load()

	// 连接数据库
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	defer func() {
		if sqlDB, err := database.GetDB().DB(); err == nil && sqlDB != nil {
			sqlDB.Close()
		}
	}()

	// // 数据库自动迁移
	// if err := database.AutoMigrate(); err != nil {
	// 	log.Fatalf("Failed to auto migrate database: %v", err)
	// }

	// 初始化路由
	r := router.Router()

	// 配置 CORS
	setupCORS(r, config.C.CORSOrigins)

	// 挂载上传文件目录
	uploadsDir := "./uploads"
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		os.MkdirAll(uploadsDir, 0755)
	}
	r.Static("/uploads", uploadsDir)

	// 挂载 Vue 管理后台（HTML5 模式支持）
	adminDist := "./admin/dist"
	if _, err := os.Stat(adminDist); err == nil {
		// 静态文件服务
		r.Static("/admin", adminDist)
		// HTML5 模式：404 时重定向到 index.html（支持 Vue Router history 模式）
		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			// 只处理 /admin 路径的 404，其他路径直接返回 404
			if len(path) > 6 && path[:6] == "/admin" {
				indexHTML, err := os.ReadFile(filepath.Join(adminDist, "index.html"))
				if err == nil {
					c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
					return
				}
			}
			// 其他路径不处理，Gin 会自动返回 404
		})
	}

	// 启动服务
	addr := fmt.Sprintf(":%d", config.C.Port)
	log.Printf("Server starting on port %d", config.C.Port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// setupCORS 配置 CORS 中间件
func setupCORS(r *gin.Engine, corsOrigins []string) {
	// 如果没有配置 CORS 来源，默认允许所有
	if len(corsOrigins) == 0 || (len(corsOrigins) == 1 && corsOrigins[0] == "") {
		corsOrigins = []string{"*"}
	}

	r.Use(func(c *gin.Context) {
		// 设置允许的来源（如果有多个，取第一个或返回 *）
		if len(corsOrigins) == 1 && corsOrigins[0] == "*" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else if len(corsOrigins) > 0 {
			c.Writer.Header().Set("Access-Control-Allow-Origin", corsOrigins[0])
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})
}
