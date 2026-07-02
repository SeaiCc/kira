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

// SetupTagsRoutes 设置标签路由
func SetupTagsRoutes(r *gin.RouterGroup) {
	// 公开路由（无需认证）
	{
		// GET /tags 列出所有标签
		r.GET("", listTags)
	}

	// 需要认证的路由
	{
		auth := r.Use(middleware.JWTMiddleware())
		// POST /tags 创建标签
		auth.POST("", createTag)
		// PUT /tags/{tag_id} 更新标签
		auth.PUT("/:tag_id", updateTag)
		// DELETE /tags/{tag_id} 删除标签
		auth.DELETE("/:tag_id", deleteTag)
	}
}

// listTags 获取所有标签（按文章数降序）
func listTags(c *gin.Context) {
	// 调用 Service
	tags, err := service.GetTags()
	if err != nil {
		utils.InternalServerError(c, "获取标签列表失败："+err.Error())
		return
	}

	// 直接返回数组（与 Python 一致）
	c.JSON(http.StatusOK, tags)
}

// createTag 创建标签
func createTag(c *gin.Context) {
	var data schemas.TagCreate

	// 解析请求体
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.BadRequest(c, "请求参数错误："+err.Error())
		return
	}

	// 调用 Service
	result, err := service.CreateTag(data)
	if err != nil {
		utils.InternalServerError(c, "创建标签失败："+err.Error())
		return
	}

	// 直接返回对象（与 Python 一致，FastAPI POST 默认 200）
	c.JSON(http.StatusOK, result)
}

// updateTag 更新标签
func updateTag(c *gin.Context) {
	tagIDStr := c.Param("tag_id")
	tagID, err := strconv.ParseUint(tagIDStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误：无效的 tag_id")
		return
	}

	var data schemas.TagUpdate
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.BadRequest(c, "请求参数错误："+err.Error())
		return
	}

	// 调用 Service
	result, err := service.UpdateTag(uint(tagID), data)
	if err != nil {
		if err.Error() == "标签不存在" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalServerError(c, "更新标签失败："+err.Error())
		return
	}

	// 直接返回对象（与 Python 一致）
	c.JSON(http.StatusOK, result)
}

// deleteTag 删除标签
func deleteTag(c *gin.Context) {
	tagIDStr := c.Param("tag_id")
	tagID, err := strconv.ParseUint(tagIDStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误：无效的 tag_id")
		return
	}

	// 调用 Service
	err = service.DeleteTag(uint(tagID))
	if err != nil {
		if err.Error() == "标签不存在" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalServerError(c, "删除标签失败："+err.Error())
		return
	}

	// 返回响应
	utils.Success(c, gin.H{"ok": true})
}
