package schemas

import (
	"time"
)

// PostCreate 创建文章请求
// 对应 Python: app/schemas/post.py -> PostCreate
type PostCreate struct {
	Title       string   `json:"title" binding:"required"`
	Slug        string   `json:"slug" binding:"required"`
	Description string   `json:"description"`
	Content     string   `json:"content"`
	Cover       string   `json:"cover"`
	CategoryID  *uint    `json:"category_id"`
	Tags        []string `json:"tags"` // 标签名称列表
	Status      string   `json:"status"`
	IsPinned    bool     `json:"is_pinned"`
	ReadingTime int      `json:"reading_time"`
	WordCount   int      `json:"word_count"`
}

// PostUpdate 更新文章请求
// 对应 Python: app/schemas/post.py -> PostUpdate
type PostUpdate struct {
	Title       *string   `json:"title"`
	Slug        *string   `json:"slug"`
	Description *string   `json:"description"`
	Content     *string   `json:"content"`
	Cover       *string   `json:"cover"`
	CategoryID  *uint     `json:"category_id"`
	Tags        *[]string `json:"tags"` // 标签名称列表
	Status      *string   `json:"status"`
	IsPinned    *bool     `json:"is_pinned"`
	ReadingTime *int      `json:"reading_time"`
	WordCount   *int      `json:"word_count"`
}

// PostOut 文章输出 DTO
// 对应 Python: app/schemas/post.py -> PostOut
// 特殊字段映射：
// - category: 从 CategoryID 或 Category 对象转换为字符串名称
// - tags: 从 []Tag 数组转换为 []string 标签名称列表
type PostOut struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Description string     `json:"description"`
	Cover       string     `json:"cover"`
	Category    string     `json:"category"` // 分类名称（不是 ID）
	Tags        []string   `json:"tags"`     // 标签名称列表（不是对象）
	Status      string     `json:"status"`
	IsPinned    bool       `json:"is_pinned"`
	Views       int        `json:"views"`
	Likes       int        `json:"likes"`
	WordCount   int        `json:"word_count"`
	ReadingTime int        `json:"reading_time"`
	PublishedAt *time.Time `json:"published_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// PostDetail 文章详情输出 DTO
// 对应 Python: app/schemas/post.py -> PostDetail (继承 PostOut + content)
type PostDetail struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Description string     `json:"description"`
	Cover       string     `json:"cover"`
	Category    string     `json:"category"`
	Tags        []string   `json:"tags"`
	Status      string     `json:"status"`
	IsPinned    bool       `json:"is_pinned"`
	Views       int        `json:"views"`
	Likes       int        `json:"likes"`
	WordCount   int        `json:"word_count"`
	ReadingTime int        `json:"reading_time"`
	PublishedAt *time.Time `json:"published_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Content     string     `json:"content"` // 详情才包含内容
}
