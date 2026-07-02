package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"kira-go/internal/middleware"
	"kira-go/internal/service"
	"kira-go/internal/utils"
)

// SetupUploadRoutes 设置上传路由
func SetupUploadRoutes(r *gin.RouterGroup) {
	// 需要认证的路由
	{
		auth := r.Use(middleware.JWTMiddleware())
		// POST /upload/image 上传图片
		auth.POST("/image", uploadImage)
	}
}

// uploadImage 上传图片
func uploadImage(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "未找到文件")
		return
	}

	// 调用 OSS 服务上传
	result, err := service.UploadImage(file)
	if err != nil {
		if err.Error() == "不支持的文件类型" {
			utils.BadRequest(c, err.Error())
			return
		}
		if err.Error() == "文件大小不能超过 10MB" {
			utils.BadRequest(c, err.Error())
			return
		}
		utils.InternalServerError(c, "上传失败："+err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}
