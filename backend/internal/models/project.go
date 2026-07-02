package models

import (
	"time"
)

// Project 项目表
type Project struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`                     // 主键
	Name         string    `gorm:"size:100;not null" json:"name"`                         // 项目名称
	Slug         string    `gorm:"size:100;uniqueIndex;not null" json:"slug"`             // URL 友好标识
	Description  string    `gorm:"size:500;default:''" json:"description"`                // 简短描述
	LongDesc     string    `gorm:"type:text;default:''" json:"long_description"`          // 详细描述
	CoverImage   string    `gorm:"size:500;default:''" json:"cover_image"`                // 封面图片
	TechStack    string    `gorm:"default:'[]'" json:"tech_stack"`                       // 技术栈 JSON
	LinkGithub   string    `gorm:"size:300;default:''" json:"link_github"`                // GitHub 链接
	LinkGitee    string    `gorm:"size:300;default:''" json:"link_gitee"`                 // Gitee 链接
	LinkLive     string    `gorm:"size:300;default:''" json:"link_live"`                  // 在线演示链接
	LinkDocs     string    `gorm:"size:300;default:''" json:"link_docs"`                  // 文档链接
	Status       string    `gorm:"size:20;default:'developing'" json:"status"`             // 状态：developing | completed | archived | planned | maintenance
	StatusLabel  string    `gorm:"size:20;default:''" json:"status_label"`                 // 状态中文标签
	IsFeatured   bool      `gorm:"default:false" json:"is_featured"`                     // 是否精选
	Sort         int       `gorm:"default:0" json:"sort"`                               // 排序
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`                    // 创建时间
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`                    // 更新时间
}

// TableName 指定表名
func (Project) TableName() string {
	return "project"
}
