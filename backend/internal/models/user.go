package models

import (
	"time"
)

// User 用户模型，对应数据库表 user
type User struct {
	ID int `gorm:"primaryKey;autoIncrement" json:"id"` // 主键 ID
	// GORM uniqueIndex;index 语义有问题，且GORM无智能合并机制 会创建两个相同名索引idx_user_username
	// GORM一般只会间接触发数据库引擎的约束相关操作，唯一索引 idx_user_username 数据库约束 user_username_key
	Username       string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"` // 用户名，唯一索引
	HashedPassword string    `gorm:"type:varchar(128);not null" json:"-"`                   // 哈希密码（不序列化）
	Nickname       string    `gorm:"type:varchar(50);default:''" json:"nickname"`           // 昵称
	Avatar         string    `gorm:"type:varchar(500);default:''" json:"avatar"`            // 头像图片地址
	Email          string    `gorm:"type:varchar(100);default:''" json:"email"`             // 邮箱
	Bio            string    `gorm:"type:varchar(500);default:''" json:"bio"`               // 个人简介
	IsAdmin        bool      `gorm:"default:false" json:"is_admin"`                         // 是否管理员
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`                      // 创建时间
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`                      // 更新时间
}

// TableName 指定表名
func (User) TableName() string {
	return "user"
}

// GitHubUser GitHub 用户模型，对应数据库表 github_user
type GitHubUser struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`         // 主键 ID
	GitHubID  int       `gorm:"uniqueIndex;index" json:"github_id"`         // GitHub 用户 ID，唯一索引
	Login     string    `gorm:"type:varchar(100)" json:"login"`             // GitHub 用户名
	Avatar    string    `gorm:"type:varchar(500);default:''" json:"avatar"` // GitHub 头像
	Bio       string    `gorm:"type:varchar(500);default:''" json:"bio"`    // GitHub 简介
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`           // 创建时间
}

// TableName 指定表名
func (GitHubUser) TableName() string {
	return "github_user"
}
