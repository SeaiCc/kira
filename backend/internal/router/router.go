package router

import (
	"kira-go/internal/api"

	"github.com/gin-gonic/gin"
)

// Router 创建主路由并聚合所有子路由
// 对应 Python: app/api/router.py
func Router() *gin.Engine {
	r := gin.New()

	// 健康检查
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 动态路由接口（对应 Python main.py:48）
	r.GET("/api/routes", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code":    0,
			"message": "success",
			"data":    []interface{}{},
		})
	})

	// API 路由组
	apiGroup := r.Group("/api")
	{
		// 认证路由
		auth := apiGroup.Group("/auth")
		api.SetupAuthRoutes(auth)

		// GitHub OAuth 路由
		github := apiGroup.Group("/auth/github")
		api.SetupGitHubRoutes(github)

		// 文章路由
		posts := apiGroup.Group("/posts")
		api.SetupPostsRoutes(posts)

		// 分类路由
		categories := apiGroup.Group("/categories")
		api.SetupCategoriesRoutes(categories)

		// 标签路由
		tags := apiGroup.Group("/tags")
		api.SetupTagsRoutes(tags)

		// 相册路由
		albums := apiGroup.Group("/albums")
		api.SetupAlbumsRoutes(albums)

		// 项目路由
		projects := apiGroup.Group("/projects")
		api.SetupProjectRoutes(projects)

		// 站点配置路由
		siteConfig := apiGroup.Group("/site-config")
		api.SetupSiteConfigRoutes(siteConfig)

		// 上传路由
		upload := apiGroup.Group("/upload")
		api.SetupUploadRoutes(upload)

		// 访客路由
		visitors := apiGroup.Group("/visitors")
		api.SetupVisitorRoutes(visitors)

		// 仪表盘路由
		dashboard := apiGroup.Group("/dashboard")
		api.SetupDashboardRoutes(dashboard)
	}

	return r
}
