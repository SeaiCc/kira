package service

import (
	"errors"

	"kira-go/internal/database"
	"kira-go/internal/models"
	"kira-go/internal/schemas"

	"gorm.io/gorm"
)

// GetCategories 获取所有分类（按 sort 排序）
func GetCategories() ([]*models.Category, error) {
	db := database.GetDB()

	var categories []*models.Category
	if err := db.Order("sort").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// GetCategoryByID 根据 ID 获取分类
func GetCategoryByID(catID uint) (*models.Category, error) {
	db := database.GetDB()

	category := &models.Category{}
	if err := db.First(category, catID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("分类不存在")
		}
		return nil, err
	}
	return category, nil
}

// CreateCategory 创建分类
func CreateCategory(data schemas.CategoryCreate) (*models.Category, error) {
	db := database.GetDB()

	category := &models.Category{
		Name:        data.Name,
		Slug:        data.Slug,
		Description: data.Description,
		Sort:        data.Sort,
	}

	if err := db.Create(category).Error; err != nil {
		return nil, err
	}

	return category, nil
}

// UpdateCategory 更新分类
func UpdateCategory(catID uint, data schemas.CategoryUpdate) (*models.Category, error) {
	db := database.GetDB()

	// 查询现有分类
	category := &models.Category{}
	if err := db.First(category, catID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("分类不存在")
		}
		return nil, err
	}

	// 直接赋值更新（对应 Python: for k, v in data.model_dump(exclude_unset=True).items(): setattr(cat, k, v)）
	if data.Name != nil {
		category.Name = *data.Name
	}
	if data.Slug != nil {
		category.Slug = *data.Slug
	}
	if data.Description != nil {
		category.Description = *data.Description
	}
	if data.Sort != nil {
		category.Sort = *data.Sort
	}

	// 更新（Save 会更新所有字段包括零值，并自动更新 UpdatedAt）
	if err := db.Save(category).Error; err != nil {
		return nil, err
	}

	return category, nil
}

// DeleteCategory 删除分类
func DeleteCategory(catID uint) error {
	db := database.GetDB()

	category := &models.Category{}
	if err := db.First(category, catID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("分类不存在")
		}
		return err
	}

	if err := db.Delete(category).Error; err != nil {
		return err
	}

	return nil
}
