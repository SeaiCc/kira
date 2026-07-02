package api

import (
	"net/http"
	"strconv"

	"kira-go/internal/middleware"
	"kira-go/internal/schemas"
	"kira-go/internal/service"
	"kira-go/internal/utils"

	"github.com/gin-gonic/gin"
)

// SetupCategoriesRoutes 设置分类路由
func SetupCategoriesRoutes(r *gin.RouterGroup) {
	// 公开路由（无需认证）
	{
		// GET /categories 列出所有分类
		r.GET("", listCategories)
	}

	// 需要认证的路由
	{
		auth := r.Use(middleware.JWTMiddleware())
		// POST /categories 创建分类
		auth.POST("", createCategory)
		// PUT /categories/{cat_id} 更新分类
		auth.PUT("/:cat_id", updateCategory)
		// DELETE /categories/{cat_id} 删除分类
		auth.DELETE("/:cat_id", deleteCategory)
	}
}

// listCategories 获取所有分类
func listCategories(c *gin.Context) {
	// 调用 Service
	categories, err := service.GetCategories()
	if err != nil {
		utils.InternalServerError(c, "获取分类列表失败："+err.Error())
		return
	}

	// 直接返回数组
	c.JSON(http.StatusOK, categories)
}

// createCategory 创建分类
func createCategory(c *gin.Context) {
	var data schemas.CategoryCreate

	// 解析请求体
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.BadRequest(c, "请求参数错误："+err.Error())
		return
	}

	// 调用 Service
	result, err := service.CreateCategory(data)
	if err != nil {
		utils.InternalServerError(c, "创建分类失败："+err.Error())
		return
	}

	// 直接返回对象（与 Python 一致，FastAPI POST 默认 200）
	c.JSON(http.StatusOK, result)
}

// updateCategory 更新分类
func updateCategory(c *gin.Context) {
	catIDStr := c.Param("cat_id")
	catID, err := strconv.ParseUint(catIDStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误：无效的 cat_id")
		return
	}

	var data schemas.CategoryUpdate
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.BadRequest(c, "请求参数错误："+err.Error())
		return
	}

	// 调用 Service
	result, err := service.UpdateCategory(uint(catID), data)
	if err != nil {
		if err.Error() == "分类不存在" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalServerError(c, "更新分类失败："+err.Error())
		return
	}

	// 直接返回对象（与 Python 一致）
	c.JSON(http.StatusOK, result)
}

// deleteCategory 删除分类
func deleteCategory(c *gin.Context) {
	catIDStr := c.Param("cat_id")
	catID, err := strconv.ParseUint(catIDStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误：无效的 cat_id")
		return
	}

	// 调用 Service
	err = service.DeleteCategory(uint(catID))
	if err != nil {
		if err.Error() == "分类不存在" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalServerError(c, "删除分类失败："+err.Error())
		return
	}

	// 返回响应
	utils.Success(c, gin.H{"ok": true})
}
