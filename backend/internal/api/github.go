package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"kira-go/internal/config"
	"kira-go/internal/database"
	"kira-go/internal/models"
	"kira-go/internal/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// SetupGitHubRoutes 设置 GitHub 认证路由
func SetupGitHubRoutes(r *gin.RouterGroup) {
	// GET /auth/github/login 重定向到 GitHub 授权页
	r.GET("/login", githubLogin)
	// GET /auth/github/callback GitHub OAuth 回调
	r.GET("/callback", githubCallback)
	// GET /auth/github/me 获取当前 GitHub 用户信息
	r.GET("/me", githubCurrentUser)
}

// getGitHubConfig 获取 GitHub OAuth 配置（懒加载）
func getGitHubConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.C.GitHubClientID,
		ClientSecret: config.C.GitHubClientSecret,
		RedirectURL:  "",
		Scopes:       []string{"read:user"},
		Endpoint:     github.Endpoint,
	}
}

// githubLogin GitHub 登录（重定向到授权页）
func githubLogin(c *gin.Context) {
	if config.C.GitHubClientID == "" {
		utils.InternalServerError(c, "未配置 GITHUB_CLIENT_ID")
		return
	}

	// 生成 state（可选，用于 CSRF 保护）
	state := fmt.Sprintf("%d", time.Now().Unix())

	// 重定向到 GitHub 授权页
	ghConfig := getGitHubConfig()
	authURL := ghConfig.AuthCodeURL("", oauth2.SetAuthURLParam("state", state))
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// githubCallback GitHub OAuth 回调
func githubCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		utils.BadRequest(c, "缺少 code 参数")
		return
	}

	// 获取配置
	ghConfig := getGitHubConfig()

	// 1. 用 code 换取 access_token
	token, err := ghConfig.Exchange(c.Request.Context(), code)
	if err != nil {
		utils.BadRequest(c, "GitHub 授权失败："+err.Error())
		return
	}

	// 2. 获取 GitHub 用户信息
	client := ghConfig.Client(c.Request.Context(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		utils.BadRequest(c, "获取 GitHub 用户信息失败："+err.Error())
		return
	}
	defer resp.Body.Close()

	ghUser := &GitHubUser{}
	if err := json.NewDecoder(resp.Body).Decode(ghUser); err != nil {
		utils.BadRequest(c, "解析 GitHub 用户信息失败："+err.Error())
		return
	}

	// 3. 查找或创建 github_user
	db := database.GetDB()
	user := &models.GitHubUser{}
	result := db.Where("github_id = ?", ghUser.ID).First(user)

	if result.Error != nil {
		// 用户不存在，创建新用户
		user = &models.GitHubUser{
			GitHubID: ghUser.ID,
			Login:    ghUser.Login,
			Avatar:   ghUser.AvatarURL,
			Bio:      ghUser.Bio,
		}
		if err := db.Create(user).Error; err != nil {
			utils.InternalServerError(c, "创建用户失败："+err.Error())
			return
		}
	} else {
		// 用户已存在，更新信息
		user.Login = ghUser.Login
		user.Avatar = ghUser.AvatarURL
		user.Bio = ghUser.Bio
		if err := db.Save(user).Error; err != nil {
			utils.InternalServerError(c, "更新用户失败："+err.Error())
			return
		}
	}

	// 4. 签发 JWT
	tokenData := map[string]any{
		"sub":   user.ID,
		"login": user.Login,
		"type":  "github",
	}
	jwtToken, err := utils.CreateToken(tokenData)
	if err != nil {
		utils.InternalServerError(c, "生成 token 失败："+err.Error())
		return
	}

	// 5. 重定向回前端统一回调页，带上 token
	frontendOrigin := getEnvOrDefault("FRONTEND_ORIGIN", "https://boke.hiromu.top")
	redirectURL := frontendOrigin + "/auth/callback?token=" + url.QueryEscape(jwtToken)
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// githubCurrentUser 获取当前 GitHub 用户信息
func githubCurrentUser(c *gin.Context) {
	// 从 Authorization header 解析 JWT
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		utils.Unauthorized(c, "未登录")
		return
	}

	// 解析 token
	claims, err := utils.DecodeToken(authHeader[7:])
	if err != nil {
		utils.Unauthorized(c, "登录已过期，请重新登录")
		return
	}

	// 获取 user_id (sub)
	sub, ok := claims["sub"].(float64) // JWT 中数字默认解析为 float64
	if !ok {
		utils.Unauthorized(c, "登录已过期，请重新登录")
		return
	}
	userID := int(sub)

	// 查询用户
	db := database.GetDB()
	user := &models.GitHubUser{}
	result := db.Where("id = ?", userID).First(user)
	if result.Error != nil {
		utils.Unauthorized(c, "用户不存在")
		return
	}

	// 返回用户信息
	c.JSON(http.StatusOK, gin.H{
		"id":     user.ID,
		"login":  user.Login,
		"avatar": user.Avatar,
		"bio":    user.Bio,
	})
}

// GitHubUser GitHub API 响应结构
type GitHubUser struct {
	ID          int    `json:"id"`
	Login       string `json:"login"`
	AvatarURL   string `json:"avatar_url"`
	Bio         string `json:"bio"`
	HTMLURL     string `json:"html_url"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Company     string `json:"company"`
	Blog        string `json:"blog"`
	Location    string `json:"location"`
	Email       string `json:"email"`
	Followers   int    `json:"followers"`
	Following   int    `json:"following"`
	PublicRepos int    `json:"public_repos"`
}

// getEnvOrDefault 获取环境变量，如果不存在则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
