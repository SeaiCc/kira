package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"kira-go/internal/middleware"
	"kira-go/internal/schemas"
	"kira-go/internal/service"
	"kira-go/internal/utils"
)

// SetupProjectRoutes 设置项目路由
func SetupProjectRoutes(r *gin.RouterGroup) {
	// 公开路由
	{
		// GET /projects 列出所有项目
		r.GET("", listProjects)
		// GET /projects/{slug} 根据 slug 获取项目
		r.GET("/:slug", getProjectBySlug)
	}

	// 需要认证的路由
	{
		auth := r.Use(middleware.JWTMiddleware())
		// POST /projects 创建项目
		auth.POST("", createProject)
		// PUT /projects/{project_id} 更新项目
		auth.PUT("/:project_id", updateProject)
		// DELETE /projects/{project_id} 删除项目
		auth.DELETE("/:project_id", deleteProject)
	}
}

// listProjects 列出所有项目
func listProjects(c *gin.Context) {
	projects, err := service.GetProjects()
	if err != nil {
		utils.InternalServerError(c, "获取项目列表失败："+err.Error())
		return
	}

	c.JSON(http.StatusOK, projects)
}

// getProjectBySlug 根据 slug 获取项目
func getProjectBySlug(c *gin.Context) {
	slug := c.Param("slug")

	project, err := service.GetProjectBySlug(slug)
	if err != nil {
		if err.Error() == "项目不存在" {
			utils.NotFound(c, "项目不存在")
			return
		}
		utils.InternalServerError(c, "获取项目失败："+err.Error())
		return
	}

	c.JSON(http.StatusOK, project)
}

// createProject 创建项目
func createProject(c *gin.Context) {
	var data schemas.ProjectCreate

	if err := c.ShouldBindJSON(&data); err != nil {
		utils.BadRequest(c, "请求参数错误："+err.Error())
		return
	}

	project, err := service.CreateProject(data)
	if err != nil {
		utils.InternalServerError(c, "创建项目失败："+err.Error())
		return
	}

	c.JSON(http.StatusOK, project)
}

// updateProject 更新项目
func updateProject(c *gin.Context) {
	projectIDStr := c.Param("project_id")

	// 解析 project_id
	var projectID uint
	if _, err := scanInt(projectIDStr, &projectID); err != nil {
		utils.BadRequest(c, "无效的项目 ID")
		return
	}

	// 解析请求体
	var data schemas.ProjectUpdate

	if err := c.ShouldBindJSON(&data); err != nil {
		utils.BadRequest(c, "请求参数错误："+err.Error())
		return
	}

	project, err := service.UpdateProject(projectID, data)
	if err != nil {
		if err.Error() == "项目不存在" {
			utils.NotFound(c, "项目不存在")
			return
		}
		utils.InternalServerError(c, "更新项目失败："+err.Error())
		return
	}

	c.JSON(http.StatusOK, project)
}

// deleteProject 删除项目
func deleteProject(c *gin.Context) {
	projectIDStr := c.Param("project_id")

	// 解析 project_id
	var projectID uint
	if _, err := scanInt(projectIDStr, &projectID); err != nil {
		utils.BadRequest(c, "无效的项目 ID")
		return
	}

	if err := service.DeleteProject(projectID); err != nil {
		if err.Error() == "项目不存在" {
			utils.NotFound(c, "项目不存在")
			return
		}
		utils.InternalServerError(c, "删除项目失败："+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// scanInt 辅助函数，将字符串转为 uint
func scanInt(s string, dest *uint) (int, error) {
	v, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}
	*dest = uint(v)
	return 1, nil
}
