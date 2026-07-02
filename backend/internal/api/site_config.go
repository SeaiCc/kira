package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"kira-go/internal/middleware"
	"kira-go/internal/schemas"
	"kira-go/internal/service"
	"kira-go/internal/utils"
)

// SetupSiteConfigRoutes 设置站点配置路由
func SetupSiteConfigRoutes(r *gin.RouterGroup) {
	// 公开路由
	{
		// GET /site-config 获取所有配置，key-value 字典
		r.GET("", getAllConfig)
		// GET /site-config/list 获取所有配置，列表
		r.GET("/list", getAllConfigList)
		// GET /site-config/{key} 根据 key 获取配置
		r.GET("/:key", getConfigByKey)
	}

	// 需要认证的路由
	{
		auth := r.Use(middleware.JWTMiddleware())
		// POST /site-config 创建站点配置
		auth.POST("", createConfig)
		// PUT /site-config/{key} 修改站点配置
		auth.PUT("/:key", updateConfig)
		// PUT /site-config 批量更新配置
		auth.PUT("", batchUpdateConfig)
		// DELETE /site-config/{key} 删除站点配置
		auth.DELETE("/:key", deleteConfig)
	}
}

// getAllConfig 获取所有配置（key-value 字典）
func getAllConfig(c *gin.Context) {
	configs, err := service.GetAllConfig()
	if err != nil {
		utils.InternalServerError(c, "获取配置失败："+err.Error())
		return
	}

	c.JSON(http.StatusOK, configs)
}

// getAllConfigList 获取所有配置列表
func getAllConfigList(c *gin.Context) {
	configs, err := service.GetAllConfigList()
	if err != nil {
		utils.InternalServerError(c, "获取配置列表失败："+err.Error())
		return
	}

	c.JSON(http.StatusOK, configs)
}

// getConfigByKey 根据 key 获取配置
func getConfigByKey(c *gin.Context) {
	key := c.Param("key")

	value, err := service.GetConfig(key)
	if err != nil {
		if err.Error() == "配置不存在" {
			utils.NotFound(c, "配置不存在")
			return
		}
		utils.InternalServerError(c, "获取配置失败："+err.Error())
		return
	}

	c.JSON(http.StatusOK, value)
}

// createConfig 创建配置
func createConfig(c *gin.Context) {
	var data struct {
		Key         string `json:"key" binding:"required"`
		Value       string `json:"value"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		utils.BadRequest(c, "请求参数错误："+err.Error())
		return
	}

	config, err := service.CreateConfig(data.Key, data.Value, data.Description)
	if err != nil {
		if err.Error() == "配置已存在" {
			utils.BadRequest(c, "配置已存在")
			return
		}
		utils.InternalServerError(c, "创建配置失败："+err.Error())
		return
	}

	c.JSON(http.StatusOK, config)
}

// updateConfig 更新配置
func updateConfig(c *gin.Context) {
	key := c.Param("key")

	var data schemas.SiteConfigUpdate

	if err := c.ShouldBindJSON(&data); err != nil {
		utils.BadRequest(c, "请求参数错误："+err.Error())
		return
	}

	config, err := service.UpdateConfig(key, data)
	if err != nil {
		utils.InternalServerError(c, "更新配置失败："+err.Error())
		return
	}

	c.JSON(http.StatusOK, config)
}

// batchUpdateConfig 批量更新配置
func batchUpdateConfig(c *gin.Context) {
	var configs map[string]interface{}

	if err := c.ShouldBindJSON(&configs); err != nil {
		utils.BadRequest(c, "请求参数错误："+err.Error())
		return
	}

	result, err := service.BatchUpdateConfig(configs)
	if err != nil {
		utils.InternalServerError(c, "批量更新配置失败："+err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// deleteConfig 删除配置
func deleteConfig(c *gin.Context) {
	key := c.Param("key")

	if err := service.DeleteConfig(key); err != nil {
		if err.Error() == "配置不存在" {
			utils.NotFound(c, "配置不存在")
			return
		}
		utils.InternalServerError(c, "删除配置失败："+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
