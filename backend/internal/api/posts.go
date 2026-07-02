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

// SetupPostsRoutes 设置文章路由
func SetupPostsRoutes(r *gin.RouterGroup) {
	// 公开路由（无需认证）
	{
		// GET /posts 文章列表
		r.GET("", listPosts)
		// GET /posts/count 文章数量统计
		r.GET("/count", postCount)
		// GET /posts/detail/{post_id} 根据 ID 获取文章详情
		r.GET("/detail/:post_id", getPostByID)
		// GET /posts/{slug} 根据 slug 获取文章详情
		r.GET("/:slug", getPostBySlug)
	}

	// 需要认证的路由
	{
		auth := r.Use(middleware.JWTMiddleware())
		// POST /posts 创建文章
		auth.POST("", createPost)
		// PUT /posts/{post_id} 更新文章
		auth.PUT("/:post_id", updatePost)
		// DELETE /posts/{post_id} 删除文章
		auth.DELETE("/:post_id", deletePost)
		// POST /posts/{post_id}/like 点赞
		auth.POST("/:post_id/like", likePost)
		// POST /posts/{post_id}/unlike 取消点赞
		auth.POST("/:post_id/unlike", unlikePost)
	}
}

// listPosts 获取文章列表
func listPosts(c *gin.Context) {
	// 解析查询参数
	status := c.Query("status")
	category := c.Query("category")
	tag := c.Query("tag")

	// 解析分页参数
	page := 1
	if p := c.Query("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}

	size := 10
	if s := c.Query("size"); s != "" {
		if val, err := strconv.Atoi(s); err == nil && val > 0 && val <= 200 {
			size = val
		}
	}

	// 调用 Service
	posts, err := service.GetPosts(status, category, tag, page, size)
	if err != nil {
		utils.InternalServerError(c, "获取文章列表失败："+err.Error())
		return
	}

	// 直接返回数组（与 Python 一致）
	c.JSON(http.StatusOK, posts)
}

// postCount 获取文章数量
func postCount(c *gin.Context) {
	status := c.Query("status")

	count, err := service.CountPosts(status)
	if err != nil {
		utils.InternalServerError(c, "统计文章数量失败："+err.Error())
		return
	}

	// 直接返回对象（与 Python 一致）
	c.JSON(http.StatusOK, gin.H{"count": count})
}

// getPostByID 根据 ID 获取文章详情
func getPostByID(c *gin.Context) {
	// 解析 post_id 参数
	postIDStr := c.Param("post_id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误：无效的 post_id")
		return
	}

	// 调用 Service
	result, err := service.GetPostByID(uint(postID))
	if err != nil {
		if err.Error() == "文章不存在" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalServerError(c, "获取文章失败："+err.Error())
		return
	}

	// 直接返回对象（与 Python 一致）
	c.JSON(http.StatusOK, result)
}

// getPostBySlug 根据 slug 获取文章详情
func getPostBySlug(c *gin.Context) {
	slug := c.Param("slug")

	// 调用 Service
	result, err := service.GetPostBySlug(slug)
	if err != nil {
		if err.Error() == "文章不存在" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalServerError(c, "获取文章失败："+err.Error())
		return
	}

	// 直接返回对象（与 Python 一致）
	c.JSON(http.StatusOK, result)
}

// createPost 创建文章
func createPost(c *gin.Context) {
	var data schemas.PostCreate

	// 解析请求体
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.BadRequest(c, "请求参数错误："+err.Error())
		return
	}

	// 调用 Service
	result, err := service.CreatePost(data)
	if err != nil {
		utils.InternalServerError(c, "创建文章失败："+err.Error())
		return
	}

	// 直接返回对象（与 Python 一致，FastAPI POST 默认 200）
	c.JSON(http.StatusOK, result)
}

// updatePost 更新文章
func updatePost(c *gin.Context) {
	postIDStr := c.Param("post_id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误：无效的 post_id")
		return
	}

	var data schemas.PostUpdate
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.BadRequest(c, "请求参数错误："+err.Error())
		return
	}

	// 调用 Service
	result, err := service.UpdatePost(uint(postID), data)
	if err != nil {
		if err.Error() == "文章不存在" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalServerError(c, "更新文章失败："+err.Error())
		return
	}

	// 直接返回对象（与 Python 一致）
	c.JSON(http.StatusOK, result)
}

// deletePost 删除文章
func deletePost(c *gin.Context) {
	postIDStr := c.Param("post_id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误：无效的 post_id")
		return
	}

	// 调用 Service
	err = service.DeletePost(uint(postID))
	if err != nil {
		if err.Error() == "文章不存在" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalServerError(c, "删除文章失败："+err.Error())
		return
	}

	// 直接返回对象（与 Python 一致）
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// likePost 点赞文章
func likePost(c *gin.Context) {
	postIDStr := c.Param("post_id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误：无效的 post_id")
		return
	}

	// 调用 Service
	count, err := service.ToggleLike(uint(postID), false)
	if err != nil {
		if err.Error() == "文章不存在" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalServerError(c, "点赞失败："+err.Error())
		return
	}

	// 直接返回对象（与 Python 一致）
	c.JSON(http.StatusOK, gin.H{"likes": count})
}

// unlikePost 取消点赞
func unlikePost(c *gin.Context) {
	postIDStr := c.Param("post_id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误：无效的 post_id")
		return
	}

	// 调用 Service
	count, err := service.ToggleLike(uint(postID), true)
	if err != nil {
		if err.Error() == "文章不存在" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalServerError(c, "取消点赞失败："+err.Error())
		return
	}

	// 直接返回对象（与 Python 一致）
	c.JSON(http.StatusOK, gin.H{"likes": count})
}
