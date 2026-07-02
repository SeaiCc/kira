package models

import (
	"time"
)

// Visitor 访客记录表
type Visitor struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`              // 主键
	IP         string    `gorm:"size:45;not null" json:"ip"`                      // IP 地址
	Path       string    `gorm:"size:500;default:''" json:"path"`                  // 访问路径
	UserAgent  string    `gorm:"default:''" json:"user_agent"`                     // User-Agent
	City       string    `gorm:"size:100;default:''" json:"city"`                  // 城市
	Region     string    `gorm:"size:100;default:''" json:"region"`                // 省份/区域
	Country    string    `gorm:"size:100;default:''" json:"country"`               // 国家
	District   string    `gorm:"size:100;default:''" json:"district"`              // 区/县
	Org        string    `gorm:"size:200;default:''" json:"org"`                   // 网络运营商/组织
	ASN        string    `gorm:"size:50;default:''" json:"asn"`                    // 自治系统编号
	IsMobile   bool      `gorm:"default:false" json:"is_mobile"`                  // 是否移动端
	IsProxy    bool      `gorm:"default:false" json:"is_proxy"`                   // 是否代理/VPN
	IsHosting  bool      `gorm:"default:false" json:"is_hosting"`                 // 是否机房 IP
	Browser    string    `gorm:"size:50;default:''" json:"browser"`                // 浏览器
	OS         string    `gorm:"size:50;default:''" json:"os"`                     // 操作系统
	DeviceType string    `gorm:"size:20;default:''" json:"device_type"`            // 设备类型：手机 | 电脑
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`                // 创建时间
}

// TableName 指定表名
func (Visitor) TableName() string {
	return "visitor"
}
