package models

import (
	"time"
)

// SiteConfig 站点配置表
type SiteConfig struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`        // 主键
	Key         string    `gorm:"size:100;uniqueIndex;not null" json:"key"`  // 站点唯一标识符
	Value       string    `gorm:"default:''" json:"value"`                   // 配置值
	Description string    `gorm:"size:200;default:''" json:"description"`    // 描述
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`          // 更新时间
}

// TableName 指定表名
func (SiteConfig) TableName() string {
	return "site_config"
}
