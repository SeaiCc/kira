package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"kira-go/internal/database"
	"kira-go/internal/models"
)

// SetupDashboardRoutes 设置仪表盘路由
func SetupDashboardRoutes(r *gin.RouterGroup) {
	// 公开路由
	{
		// GET /dashboard/stats 获取状态统计
		r.GET("/stats", getDashboardStats)
	}
}

// getDashboardStats 获取仪表盘统计数据
func getDashboardStats(c *gin.Context) {
	now := time.Now()
	thirtyDaysAgo := now.AddDate(0, 0, -30)

	// ── 总数统计 ──
	var postCount, draftCount, categoryCount, tagCount, visitorCount int64
	_ = database.GetDB().Model(&models.Post{}).Where("status = ?", "published").Count(&postCount).Error
	_ = database.GetDB().Model(&models.Post{}).Where("status = ?", "draft").Count(&draftCount).Error
	_ = database.GetDB().Model(&models.Category{}).Count(&categoryCount).Error
	_ = database.GetDB().Model(&models.Tag{}).Count(&tagCount).Error
	_ = database.GetDB().Model(&models.Visitor{}).Count(&visitorCount).Error

	// ── 文章发布趋势（近 30 天） ──
	postTrend := getPostTrend(thirtyDaysAgo)

	// ── 访客趋势（近 30 天） ──
	visitorTrend := getVisitorTrend(thirtyDaysAgo)

	// ── 分类分布 ──
	categoryDistribution := getCategoryDistribution()

	// ── 浏览器分布 ──
	browserDistribution := getBrowserDistribution()

	c.JSON(http.StatusOK, gin.H{
		"counts": gin.H{
			"posts":      postCount,
			"drafts":     draftCount,
			"categories": categoryCount,
			"tags":       tagCount,
			"visitors":   visitorCount,
		},
		"post_trend":            postTrend,
		"visitor_trend":         visitorTrend,
		"category_distribution": categoryDistribution,
		"browser_distribution":  browserDistribution,
	})
}

// getPostTrend 获取文章发布趋势
func getPostTrend(since time.Time) []map[string]interface{} {
	type PostCount struct {
		Date  string `gorm:"column:date"`
		Count int    `gorm:"column:count"`
	}

	var rows []PostCount
	db := database.GetDB()
	err := db.Model(&models.Post{}).
		Where("status = ?", "published").
		Where("created_at >= ?", since).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Group("DATE(created_at)").
		Order("DATE(created_at)").
		Scan(&rows).Error

	if err != nil {
		return []map[string]interface{}{}
	}

	// 补全 30 天空缺日期
	result := make([]map[string]interface{}, 0, 30)
	rowMap := make(map[string]int)
	for _, r := range rows {
		rowMap[r.Date] = r.Count
	}

	for i := 0; i < 30; i++ {
		date := since.AddDate(0, 0, i).Format("2006-01-02")
		result = append(result, map[string]interface{}{
			"date":  date,
			"count": rowMap[date],
		})
	}

	return result
}

// getVisitorTrend 获取访客趋势
func getVisitorTrend(since time.Time) []map[string]interface{} {
	type VisitorCount struct {
		Date  string `gorm:"column:date"`
		Count int    `gorm:"column:count"`
	}

	var rows []VisitorCount
	db := database.GetDB()
	err := db.Model(&models.Visitor{}).
		Where("created_at >= ?", since).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Group("DATE(created_at)").
		Order("DATE(created_at)").
		Scan(&rows).Error

	if err != nil {
		return []map[string]interface{}{}
	}

	// 补全 30 天空缺日期
	result := make([]map[string]interface{}, 0, 30)
	rowMap := make(map[string]int)
	for _, r := range rows {
		rowMap[r.Date] = r.Count
	}

	for i := 0; i < 30; i++ {
		date := since.AddDate(0, 0, i).Format("2006-01-02")
		result = append(result, map[string]interface{}{
			"date":  date,
			"count": rowMap[date],
		})
	}

	return result
}

// getCategoryDistribution 获取分类分布
func getCategoryDistribution() []map[string]interface{} {
	type CategoryCount struct {
		Name      string `gorm:"column:name"`
		PostCount int    `gorm:"column:post_count"`
	}

	var rows []CategoryCount
	db := database.GetDB()
	err := db.Model(&models.Category{}).
		Where("post_count > 0").
		Select("name, post_count").
		Scan(&rows).Error

	if err != nil {
		return []map[string]interface{}{}
	}

	result := make([]map[string]interface{}, 0, len(rows))
	for _, r := range rows {
		result = append(result, map[string]interface{}{
			"name":  r.Name,
			"value": r.PostCount,
		})
	}

	return result
}

// getBrowserDistribution 获取浏览器分布
func getBrowserDistribution() []map[string]interface{} {
	type BrowserCount struct {
		Browser string `gorm:"column:browser"`
		Count   int    `gorm:"column:count"`
	}

	var rows []BrowserCount
	db := database.GetDB()
	err := db.Model(&models.Visitor{}).
		Select("browser, COUNT(*) as count").
		Group("browser").
		Order("count DESC").
		Scan(&rows).Error

	if err != nil {
		return []map[string]interface{}{}
	}

	result := make([]map[string]interface{}, 0, len(rows))
	for _, r := range rows {
		browserName := r.Browser
		if browserName == "" {
			browserName = "未知"
		}
		result = append(result, map[string]interface{}{
			"name":  browserName,
			"value": r.Count,
		})
	}

	return result
}
