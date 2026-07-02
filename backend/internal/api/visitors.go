package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"kira-go/internal/middleware"
	"kira-go/internal/service"
	"kira-go/internal/utils"
)

// SetupVisitorRoutes 设置访客记录路由
func SetupVisitorRoutes(r *gin.RouterGroup) {
	// 公开路由
	{
		// GET /visitors 列出最近访客列表
		r.GET("", listRecentVisitors)
		// GET /visitors/count 总访客数
		r.GET("/count", visitorCount)
		// GET /visitors/location 当前访问者地理位置
		r.GET("/location", getVisitorLocation)
		// POST /visitors/record 记录当前访问
		r.POST("/record", recordVisitor)
	}

	// 需要认证的路由（管理员操作）
	{
		auth := r.Use(middleware.JWTMiddleware())
		// DELETE /visitors/{visitor_id} 删除单条访客记录
		auth.DELETE("/:visitor_id", deleteVisitor)
		// DELETE /visitors 清空所有访客记录
		auth.DELETE("", clearVisitors)
	}
}

// listRecentVisitors 获取最近访客列表
func listRecentVisitors(c *gin.Context) {
	// 解析分页参数
	page := 1
	size := 20

	if p := c.DefaultQuery("page", "1"); p != "" {
		if val, err := strconv.ParseUint(p, 10, 64); err == nil {
			page = int(val)
		}
	}
	if s := c.DefaultQuery("size", "20"); s != "" {
		if val, err := strconv.ParseUint(s, 10, 64); err == nil {
			size = int(val)
		}
	}

	// 限制 size 范围
	if size > 100 {
		size = 100
	}
	if size < 1 {
		size = 1
	}
	if page < 1 {
		page = 1
	}

	visitors, err := service.GetRecentVisitors(page, size)
	if err != nil {
		utils.InternalServerError(c, "获取访客列表失败："+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": visitors})
}

// visitorCount 获取总访客数
func visitorCount(c *gin.Context) {
	count, err := service.GetVisitorCount()
	if err != nil {
		utils.InternalServerError(c, "获取访客数失败："+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "count": count})
}

// getVisitorLocation 获取当前访问者的地理位置
func getVisitorLocation(c *gin.Context) {
	// 获取 IP
	ip := c.Request.Header.Get("x-forwarded-for")
	if ip == "" {
		ip = c.Request.Header.Get("x-real-ip")
	}
	if ip == "" {
		// 从 RemoteAddr 提取 IP（可能包含端口）
		ip = c.ClientIP()
	}

	// 处理 x-forwarded-for 可能包含多个 IP 的情况
	if len(ip) > 0 && strings.Contains(ip, ",") {
		parts := strings.Split(ip, ",")
		if len(parts) > 0 {
			ip = parts[0]
		}
	}

	// 调用 Service 层获取地理位置
	geo := service.FetchGeo(ip)

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": geo})
}

// recordVisitor 记录当前访问
func recordVisitor(c *gin.Context) {
	// 获取 IP
	ip := c.Request.Header.Get("x-forwarded-for")
	if ip == "" {
		ip = c.Request.Header.Get("x-real-ip")
	}
	if ip == "" {
		ip = c.ClientIP()
	}

	// 处理 x-forwarded-for 可能包含多个 IP 的情况
	if len(ip) > 0 && strings.Contains(ip, ",") {
		parts := strings.Split(ip, ",")
		if len(parts) > 0 {
			ip = parts[0]
		}
	}

	path := c.Request.Header.Get("x-path")
	userAgent := c.Request.Header.Get("user-agent")

	_, err := service.RecordVisit(ip, path, userAgent)
	if err != nil {
		utils.InternalServerError(c, "记录访问失败："+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// deleteVisitor 删除单条访客记录
func deleteVisitor(c *gin.Context) {
	visitorIDStr := c.Param("visitor_id")

	val, err := strconv.ParseUint(visitorIDStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "无效的访客 ID")
		return
	}
	visitorID := uint(val)

	if err := service.DeleteVisitor(visitorID); err != nil {
		utils.InternalServerError(c, "删除访客记录失败："+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// clearVisitors 清空所有访客记录
func clearVisitors(c *gin.Context) {
	if err := service.ClearVisitors(); err != nil {
		utils.InternalServerError(c, "清空访客记录失败："+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}
