package service

import (
	"errors"
	"kira-go/internal/database"
	"kira-go/internal/models"
	"kira-go/internal/schemas"

	"encoding/json"

	"gorm.io/gorm"
)

// GetAllConfig 返回所有配置的 key-value 字典
func GetAllConfig() (map[string]interface{}, error) {
	db := database.GetDB()
	var configs []*models.SiteConfig

	if err := db.Find(&configs).Error; err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	for _, c := range configs {
		var value interface{}
		// 尝试解析为 JSON，如果失败则使用原始字符串
		if err := json.Unmarshal([]byte(c.Value), &value); err != nil {
			value = c.Value
		}
		result[c.Key] = value
	}

	return result, nil
}

// GetAllConfigList 返回所有配置的完整列表（含 id、description、updated_at）
func GetAllConfigList() ([]map[string]interface{}, error) {
	db := database.GetDB()
	var configs []*models.SiteConfig

	if err := db.Order("id").Find(&configs).Error; err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0, len(configs))
	for _, c := range configs {
		result = append(result, map[string]interface{}{
			"id":          c.ID,
			"key":         c.Key,
			"value":       c.Value,
			"description": c.Description,
			"updated_at":  c.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	return result, nil
}

// CreateConfig 新建配置项（保证 key 唯一）
func CreateConfig(key, value, description string) (*models.SiteConfig, error) {
	db := database.GetDB()

	// 检查 key 是否已存在
	existing := &models.SiteConfig{}
	if err := db.Where("key = ?", key).First(existing).Error; err == nil {
		return nil, errors.New("配置已存在")
	}

	row := &models.SiteConfig{
		Key:         key,
		Value:       value,
		Description: description,
	}

	if err := db.Create(row).Error; err != nil {
		return nil, err
	}

	return row, nil
}

// DeleteConfig 删除配置项
func DeleteConfig(key string) error {
	db := database.GetDB()

	row := &models.SiteConfig{}
	if err := db.Where("key = ?", key).First(row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("配置不存在")
		}
		return err
	}

	if err := db.Delete(row).Error; err != nil {
		return err
	}

	return nil
}

// GetConfig 获取单个配置值（解析 JSON）
func GetConfig(key string) (interface{}, error) {
	db := database.GetDB()

	row := &models.SiteConfig{}
	if err := db.Where("key = ?", key).First(row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("配置不存在")
		}
		return nil, err
	}

	var value interface{}
	if err := json.Unmarshal([]byte(row.Value), &value); err != nil {
		value = row.Value
	}

	return value, nil
}

// UpdateConfig 更新配置（如果不存在则创建）
func UpdateConfig(key string, data schemas.SiteConfigUpdate) (*models.SiteConfig, error) {
	db := database.GetDB()

	row := &models.SiteConfig{}
	if err := db.Where("key = ?", key).First(row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建新配置（GORM 自动设置 UpdatedAt）
			row = &models.SiteConfig{
				Key:         key,
				Value:       data.Value,
				Description: data.Description,
			}
			if err := db.Create(row).Error; err != nil {
				return nil, err
			}
			return row, nil
		}
		return nil, err
	}

	// 更新现有配置（GORM 自动更新 UpdatedAt）
	row.Value = data.Value
	if data.Description != "" {
		row.Description = data.Description
	}

	if err := db.Save(row).Error; err != nil {
		return nil, err
	}

	return row, nil
}

// BatchUpdateConfig 批量更新配置
func BatchUpdateConfig(configs map[string]interface{}) (map[string]interface{}, error) {
	db := database.GetDB()

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for key, value := range configs {
		// 序列化值为 JSON 字符串
		valueJSON, err := json.Marshal(value)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		row := &models.SiteConfig{}
		if err := tx.Where("key = ?", key).First(row).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 创建新配置（GORM 自动设置 UpdatedAt）
				row = &models.SiteConfig{
					Key:   key,
					Value: string(valueJSON),
				}
				if err := tx.Create(row).Error; err != nil {
					tx.Rollback()
					return nil, err
				}
			} else {
				tx.Rollback()
				return nil, err
			}
		} else {
			// 更新现有配置（GORM 自动更新 UpdatedAt）
			row.Value = string(valueJSON)
			if err := tx.Save(row).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// 返回所有配置
	return GetAllConfig()
}
