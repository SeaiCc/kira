package models

import (
	"time"
)

// Album 相册表
type Album struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`          // 主键
	Title       string    `gorm:"size:100;not null" json:"title"`              // 标题
	Description string    `gorm:"size:500;default:''" json:"description"`      // 描述
	Cover       string    `gorm:"size:500;default:''" json:"cover"`            // 封面图片 URL
	PhotoCount  int       `gorm:"default:0" json:"photo_count"`                // 照片数量
	Sort        int       `gorm:"default:0" json:"sort"`                       // 排序
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`            // 创建时间
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`            // 更新时间

	// 关联
	Photos []Photo `gorm:"foreignKey:AlbumID" json:"photos,omitempty" preload:"Photos"` // 照片一对多关联
}

// TableName 指定表名
func (Album) TableName() string {
	return "album"
}

// Photo 照片表
type Photo struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`           // 主键
	AlbumID     uint      `gorm:"index;not null" json:"album_id"`               // 相册外键
	URL         string    `gorm:"size:500;not null" json:"url"`                 // 图片 URL
	Caption     string    `gorm:"size:200;default:''" json:"caption"`            // 照片标题
	Orientation string    `gorm:"size:20;default:'landscape'" json:"orientation"` // 照片朝向：landscape | portrait | square
	Sort        int       `gorm:"default:0" json:"sort"`                       // 排序
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`            // 创建时间

	// 关联
	Album *Album `gorm:"foreignKey:AlbumID" json:"album,omitempty" preload:"Album"` // 相册关联
}

// TableName 指定表名
func (Photo) TableName() string {
	return "photo"
}
