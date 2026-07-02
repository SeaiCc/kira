package service

import (
	"errors"

	"kira-go/internal/database"
	"kira-go/internal/models"
	"kira-go/internal/schemas"

	"gorm.io/gorm"
)

// GetTags 获取所有标签（按文章数降序）
func GetTags() ([]*models.Tag, error) {
	db := database.GetDB()

	var tags []*models.Tag
	if err := db.Order("post_count DESC").Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// CreateTag 创建标签
func CreateTag(data schemas.TagCreate) (*models.Tag, error) {
	db := database.GetDB()

	tag := &models.Tag{
		Name: data.Name,
		Slug: data.Slug,
	}

	if err := db.Create(tag).Error; err != nil {
		return nil, err
	}

	return tag, nil
}

// UpdateTag 更新标签
func UpdateTag(tagID uint, data schemas.TagUpdate) (*models.Tag, error) {
	db := database.GetDB()

	// 查询现有标签
	tag := &models.Tag{}
	if err := db.First(tag, tagID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("标签不存在")
		}
		return nil, err
	}

	// 直接赋值更新（对应 Python: for k, v in data.model_dump(exclude_unset=True).items(): setattr(tag, k, v)）
	if data.Name != nil {
		tag.Name = *data.Name
	}
	if data.Slug != nil {
		tag.Slug = *data.Slug
	}

	// 更新（Save 会更新所有字段包括零值，并自动更新 UpdatedAt）
	if err := db.Save(tag).Error; err != nil {
		return nil, err
	}

	return tag, nil
}

// DeleteTag 删除标签
func DeleteTag(tagID uint) error {
	db := database.GetDB()

	tag := &models.Tag{}
	if err := db.First(tag, tagID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("标签不存在")
		}
		return err
	}

	if err := db.Delete(tag).Error; err != nil {
		return err
	}

	return nil
}
