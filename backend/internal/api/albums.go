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

// SetupAlbumsRoutes 设置相册路由
// 对应 Python: app/api/albums.py
func SetupAlbumsRoutes(r *gin.RouterGroup) {
	// 公开路由（无需认证）
	{
		// GET /albums 列出所有相册
		r.GET("", listAlbums)
		// GET /albums/{album_id} 根据 id 获取相册信息
		r.GET("/:album_id", getAlbum)
		// GET /albums/{album_id}/photos 获取相册内所有图片
		r.GET("/:album_id/photos", getAlbumPhotos)
	}

	// 需要认证的路由
	{
		auth := r.Use(middleware.JWTMiddleware())
		// POST /albums 创建相册
		auth.POST("", createAlbum)
		// PUT /albums/{album_id} 更新相册
		auth.PUT("/:album_id", updateAlbum)
		// DELETE /albums/{album_id} 删除相册
		auth.DELETE("/:album_id", deleteAlbum)
		// POST /albums/photos 给相册添加图片
		auth.POST("/photos", addPhoto)
		// DELETE /albums/photos/{photo_id} 删除相册内的图片
		auth.DELETE("/photos/:photo_id", deletePhoto)
	}
}

// listAlbums 获取所有相册
// GET /api/albums
func listAlbums(c *gin.Context) {
	// 调用 Service
	albums, err := service.GetAlbums()
	if err != nil {
		utils.InternalServerError(c, "获取相册列表失败："+err.Error())
		return
	}

	// 直接返回数组（与 Python 一致）
	c.JSON(http.StatusOK, albums)
}

// getAlbum 根据 ID 获取相册
// GET /api/albums/{album_id}
func getAlbum(c *gin.Context) {
	albumIDStr := c.Param("album_id")
	albumID, err := strconv.ParseUint(albumIDStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误：无效的 album_id")
		return
	}

	// 调用 Service
	album, err := service.GetAlbumByID(uint(albumID))
	if err != nil {
		if err.Error() == "相册不存在" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalServerError(c, "获取相册失败："+err.Error())
		return
	}

	// 直接返回对象（与 Python 一致）
	c.JSON(http.StatusOK, album)
}

// getAlbumPhotos 获取相册中的所有照片
// GET /api/albums/{album_id}/photos
func getAlbumPhotos(c *gin.Context) {
	albumIDStr := c.Param("album_id")
	albumID, err := strconv.ParseUint(albumIDStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误：无效的 album_id")
		return
	}

	// 调用 Service
	photos, err := service.GetPhotos(uint(albumID))
	if err != nil {
		utils.InternalServerError(c, "获取照片列表失败："+err.Error())
		return
	}

	// 直接返回数组（与 Python 一致）
	c.JSON(http.StatusOK, photos)
}

// createAlbum 创建相册
// POST /api/albums
func createAlbum(c *gin.Context) {
	var data schemas.AlbumCreate

	// 解析请求体
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.BadRequest(c, "请求参数错误："+err.Error())
		return
	}

	// 调用 Service
	result, err := service.CreateAlbum(data)
	if err != nil {
		utils.InternalServerError(c, "创建相册失败："+err.Error())
		return
	}

	// 直接返回对象（与 Python 一致，FastAPI POST 默认 200）
	c.JSON(http.StatusOK, result)
}

// updateAlbum 更新相册
// PUT /api/albums/{album_id}
func updateAlbum(c *gin.Context) {
	albumIDStr := c.Param("album_id")
	albumID, err := strconv.ParseUint(albumIDStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误：无效的 album_id")
		return
	}

	var data schemas.AlbumUpdate
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.BadRequest(c, "请求参数错误："+err.Error())
		return
	}

	// 调用 Service
	result, err := service.UpdateAlbum(uint(albumID), data)
	if err != nil {
		if err.Error() == "相册不存在" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalServerError(c, "更新相册失败："+err.Error())
		return
	}

	// 直接返回对象（与 Python 一致）
	c.JSON(http.StatusOK, result)
}

// deleteAlbum 删除相册
// DELETE /api/albums/{album_id}
func deleteAlbum(c *gin.Context) {
	albumIDStr := c.Param("album_id")
	albumID, err := strconv.ParseUint(albumIDStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误：无效的 album_id")
		return
	}

	// 调用 Service
	err = service.DeleteAlbum(uint(albumID))
	if err != nil {
		if err.Error() == "相册不存在" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalServerError(c, "删除相册失败："+err.Error())
		return
	}

	// 返回响应（与 Python 一致：{"ok": true}）
	utils.Success(c, gin.H{"ok": true})
}

// addPhoto 给相册添加图片
// POST /api/albums/photos
func addPhoto(c *gin.Context) {
	var data schemas.PhotoCreate

	// 解析请求体
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.BadRequest(c, "请求参数错误："+err.Error())
		return
	}

	// 调用 Service
	result, err := service.AddPhoto(data)
	if err != nil {
		utils.InternalServerError(c, "添加照片失败："+err.Error())
		return
	}

	// 直接返回对象（与 Python 一致）
	c.JSON(http.StatusOK, result)
}

// deletePhoto 删除相册内的图片
// DELETE /api/albums/photos/{photo_id}
func deletePhoto(c *gin.Context) {
	photoIDStr := c.Param("photo_id")
	photoID, err := strconv.ParseUint(photoIDStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误：无效的 photo_id")
		return
	}

	// 调用 Service
	err = service.DeletePhoto(uint(photoID))
	if err != nil {
		if err.Error() == "照片不存在" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalServerError(c, "删除照片失败："+err.Error())
		return
	}

	// 返回响应（与 Python 一致：{"ok": true}）
	utils.Success(c, gin.H{"ok": true})
}
