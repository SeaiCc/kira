package models

import (
	"time"
)

// Post 文章表
type Post struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`              // 主键
	Title       string     `gorm:"size:200;not null" json:"title"`                  // 标题
	Slug        string     `gorm:"size:200;uniqueIndex;not null;index" json:"slug"` // URL 友好标识
	Description string     `gorm:"size:500;default:''" json:"description"`          // 描述
	Content     string     `gorm:"type:text" json:"content"`                        // 内容
	Cover       string     `gorm:"size:500;default:''" json:"cover"`                // 封面图片 URL
	CategoryID  *uint      `gorm:"index" json:"category_id"`                        // 分类外键
	Status      string     `gorm:"size:20;default:'draft';index" json:"status"`     // 状态：draft|published
	IsPinned    bool       `gorm:"default:false" json:"is_pinned"`                  // 是否置顶
	Views       int        `gorm:"default:0" json:"views"`                          // 浏览次数
	Likes       int        `gorm:"default:0" json:"likes"`                          // 点赞数
	WordCount   int        `gorm:"default:0" json:"word_count"`                     // 字数
	ReadingTime int        `gorm:"default:0" json:"reading_time"`                   // 预计阅读时间 (分钟)
	PublishedAt *time.Time `gorm:"index" json:"published_at"`                       // 发布时间
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`                // 创建时间
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`                // 更新时间
	// 关联字段（用于 Preload）
	Category    *Category  `gorm:"foreignKey:CategoryID"` // 分类
	Tags        []Tag      `gorm:"many2many:post_tag;"`  // 标签（多对多）
}

// TableName 指定表名
func (Post) TableName() string {
	return "post"
}

// Category 分类表
type Category struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`       // 主键
	Name        string    `gorm:"size:50;uniqueIndex;not null" json:"name"` // 人类可读名称
	Slug        string    `gorm:"size:50;uniqueIndex;not null" json:"slug"` // URL 友好标识
	Description string    `gorm:"size:200;default:''" json:"description"`   // 描述
	Sort        int       `gorm:"default:0" json:"sort"`                    // 排序
	PostCount   int       `gorm:"default:0" json:"post_count"`              // 文章数量
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`         // 创建时间
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`         // 更新时间
}

// TableName 指定表名
func (Category) TableName() string {
	return "category"
}

// Tag 标签表
type Tag struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`       // 主键
	Name      string `gorm:"size:50;uniqueIndex;not null" json:"name"` // 人类可读名称
	Slug      string `gorm:"size:50;uniqueIndex;not null" json:"slug"` // URL 友好标识
	PostCount int    `gorm:"default:0" json:"post_count"`              // 文章数量
}

// TableName 指定表名
func (Tag) TableName() string {
	return "tag"
}

// PostTag 文章与标签的多对多关联表
type PostTag struct {
	PostID uint `gorm:"primaryKey;autoIncrement:false" json:"post_id"` // 文章外键
	TagID  uint `gorm:"primaryKey;autoIncrement:false" json:"tag_id"`  // 标签外键
}

// TableName 指定表名
func (PostTag) TableName() string {
	return "post_tag"
}
