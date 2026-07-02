package api

import (
	"net/http"
	"time"

	"kira-go/internal/config"
	"kira-go/internal/database"
	"kira-go/internal/middleware"
	"kira-go/internal/models"
	"kira-go/internal/schemas"
	"kira-go/internal/utils"

	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes 设置认证路由
func SetupAuthRoutes(r *gin.RouterGroup) {
	// 公开路由
	{
		// POST /auth/login 登录
		r.POST("/login", login)
	}

	// 需要认证的路由
	{
		auth := r.Use(middleware.JWTMiddleware())
		// GET /auth/me 获取当前用户信息
		auth.GET("/me", getCurrentUser)
		// PUT /auth/me 更新当前用户信息
		auth.PUT("/me", updateCurrentUser)
	}
}

// login 登录处理函数
func login(c *gin.Context) {
	var req schemas.LoginRequest

	// 解析请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误："+err.Error())
		return
	}

	// 查询用户
	db := database.GetDB()
	user := &models.User{}
	result := db.Where("username = ?", req.Username).First(user)

	if result.Error != nil {
		utils.Unauthorized(c, "用户名或密码错误")
		return
	}

	// 验证密码
	if !utils.VerifyPassword(req.Password, user.HashedPassword) {
		utils.Unauthorized(c, "用户名或密码错误")
		return
	}

	// 生成 token
	tokenData := map[string]any{
		"sub":   user.Username,
		"admin": user.IsAdmin,
	}
	token, err := utils.CreateToken(tokenData)
	if err != nil {
		utils.InternalServerError(c, "生成 token 失败")
		return
	}

	// 计算过期时间
	expires := time.Now().UTC().Add(time.Duration(config.C.JWTExpireHours) * time.Hour)

	// 返回响应
	utils.SuccessWithMsg(c, "success", gin.H{
		"accessToken":  token,
		"refreshToken": "",
		"expires":      expires.Format(time.RFC3339),
		"avatar":       user.Avatar,
		"username":     user.Username,
		"nickname":     getNickname(user),
		"roles":        getRoles(user.IsAdmin),
		"permissions":  getPermissions(user.IsAdmin),
	})
}

// getCurrentUser 获取当前用户信息
func getCurrentUser(c *gin.Context) {
	username := middleware.GetCurrentUser(c)
	if username == "" {
		utils.Unauthorized(c, "未认证")
		return
	}

	// 查询用户
	db := database.GetDB()
	user := &models.User{}
	result := db.Where("username = ?", username).First(user)

	if result.Error != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	// 返回响应
	utils.Success(c, gin.H{
		"avatar":      user.Avatar,
		"username":    user.Username,
		"nickname":    getNickname(user),
		"email":       user.Email,
		"description": user.Bio,
		"phone":       "",
		"roles":       getRoles(user.IsAdmin),
		"permissions": getPermissions(user.IsAdmin),
	})
}

// updateCurrentUser 更新当前用户信息
func updateCurrentUser(c *gin.Context) {
	username := middleware.GetCurrentUser(c)
	if username == "" {
		utils.Unauthorized(c, "未认证")
		return
	}

	// 解析请求体
	var data map[string]any
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.BadRequest(c, "请求参数错误："+err.Error())
		return
	}

	// 查询用户
	db := database.GetDB()
	user := &models.User{}
	result := db.Where("username = ?", username).First(user)

	if result.Error != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	// 更新字段
	if v, ok := data["nickname"]; ok {
		if str, ok := v.(string); ok {
			user.Nickname = str
		}
	}
	if v, ok := data["email"]; ok {
		if str, ok := v.(string); ok {
			user.Email = str
		}
	}
	// 支持 bio 或 description 字段
	if v, ok := data["bio"]; ok {
		if str, ok := v.(string); ok {
			user.Bio = str
		}
	} else if v, ok := data["description"]; ok {
		if str, ok := v.(string); ok {
			user.Bio = str
		}
	}
	if v, ok := data["avatar"]; ok {
		if str, ok := v.(string); ok {
			user.Avatar = str
		}
	}

	// 更新数据库
	user.UpdatedAt = time.Now()
	if err := db.Save(user).Error; err != nil {
		utils.InternalServerError(c, "更新失败："+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// getNickname 获取昵称（如果为空则返回用户名）
func getNickname(user *models.User) string {
	if user.Nickname == "" {
		return user.Username
	}
	return user.Nickname
}

// getRoles 获取角色列表
func getRoles(isAdmin bool) []string {
	if isAdmin {
		return []string{"admin"}
	}
	return []string{}
}

// getPermissions 获取权限列表
func getPermissions(isAdmin bool) []string {
	if isAdmin {
		return []string{"*:*:*"}
	}
	return []string{}
}
