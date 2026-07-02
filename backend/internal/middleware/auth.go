package middleware

import (
	"net/http"

	"kira-go/internal/utils"

	"github.com/gin-gonic/gin"
)

// authContextKey 上下文 key 类型
type authContextKey string

const (
	// UserContextKey 用户信息上下文 key
	UserContextKey authContextKey = "sub"
)

// JWTMiddleware JWT 认证中间件（对应 Python 的 Depends(get_current_user)）
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// 检查 Authorization: Bearer <token> 格式（对应 HTTPBearer + Depends(security)）
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "认证失败",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 解析 token（对应 decode_token）
		claims, err := utils.DecodeToken(authHeader[7:])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "认证失败",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 只注入 sub 字段（username），对应 Python 的 user["sub"]
		username := claims["sub"].(string)
		c.Set(string(UserContextKey), username)
		c.Next()
	}
}

// GetCurrentUser 从 Context 获取当前用户名（对应 Python 的 user["sub"]）
// 返回 username 字符串，使用方式：username := GetCurrentUser(c)
func GetCurrentUser(c *gin.Context) string {
	user, exists := c.Get(string(UserContextKey))
	if !exists {
		return ""
	}
	return user.(string)
}
