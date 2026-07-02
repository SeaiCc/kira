package service

import (
	"errors"

	"kira-go/internal/database"
	"kira-go/internal/models"
	"kira-go/internal/schemas"

	"gorm.io/gorm"
)

// GetAlbums 获取所有相册列表（按 sort 排序）
func GetAlbums() ([]*models.Album, error) {
	db := database.GetDB()
	var albums []*models.Album
	if err := db.Order("sort").Find(&albums).Error; err != nil {
		return nil, err
	}
	return albums, nil
}

// GetAlbumByID 根据 ID 获取相册
func GetAlbumByID(albumID uint) (*models.Album, error) {
	db := database.GetDB()
	album := &models.Album{}
	if err := db.First(album, albumID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("相册不存在")
		}
		return nil, err
	}
	return album, nil
}

// CreateAlbum 创建相册
func CreateAlbum(data schemas.AlbumCreate) (*models.Album, error) {
	db := database.GetDB()
	album := &models.Album{
		Title:       data.Title,
		Description: data.Description,
		Cover:       data.Cover,
		Sort:        data.Sort,
	}
	if err := db.Create(album).Error; err != nil {
		return nil, err
	}
	return album, nil
}

// UpdateAlbum 更新相册
func UpdateAlbum(albumID uint, data schemas.AlbumUpdate) (*models.Album, error) {
	db := database.GetDB()

	// 1. 查询现有相册
	album := &models.Album{}
	if err := db.First(album, albumID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("相册不存在")
		}
		return nil, err
	}

	// 2. 更新字段（仅更新非 nil 字段）
	if data.Title != nil {
		album.Title = *data.Title
	}
	if data.Description != nil {
		album.Description = *data.Description
	}
	if data.Cover != nil {
		album.Cover = *data.Cover
	}
	if data.Sort != nil {
		album.Sort = *data.Sort
	}

	// 3. 执行更新（Save 会更新所有字段包括零值，并自动更新 UpdatedAt）
	if err := db.Save(album).Error; err != nil {
		return nil, err
	}

	return album, nil
}

// DeleteAlbum 删除相册
func DeleteAlbum(albumID uint) error {
	db := database.GetDB()
	album := &models.Album{}
	if err := db.First(album, albumID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("相册不存在")
		}
		return err
	}

	if err := db.Delete(album).Error; err != nil {
		return err
	}
	return nil
}

// GetPhotos 获取相册中的所有照片（按 sort 排序）
func GetPhotos(albumID uint) ([]*models.Photo, error) {
	db := database.GetDB()
	var photos []*models.Photo
	if err := db.Where("album_id = ?", albumID).Order("sort").Find(&photos).Error; err != nil {
		return nil, err
	}
	return photos, nil
}

// AddPhoto 添加照片到相册
func AddPhoto(data schemas.PhotoCreate) (*models.Photo, error) {
	db := database.GetDB()

	// 开始事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建照片
	photo := &models.Photo{
		AlbumID:     data.AlbumID,
		URL:         data.URL,
		Caption:     data.Caption,
		Orientation: data.Orientation,
		Sort:        data.Sort,
	}
	if err := tx.Create(photo).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 更新相册照片计数
	album := &models.Album{}
	if err := tx.First(album, data.AlbumID).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 统计该相册的照片数量
	var count int64
	if err := tx.Model(&models.Photo{}).Where("album_id = ?", data.AlbumID).Count(&count).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	album.PhotoCount = int(count)

	if err := tx.Save(album).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// 返回照片（不包含 album 关联）
	return photo, nil
}

// DeletePhoto 删除照片
func DeletePhoto(photoID uint) error {
	db := database.GetDB()

	// 开始事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 查询照片
	photo := &models.Photo{}
	if err := tx.First(photo, photoID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("照片不存在")
		}
		return err
	}

	albumID := photo.AlbumID

	// 删除照片
	if err := tx.Delete(photo).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 更新相册照片计数
	album := &models.Album{}
	if err := tx.First(album, albumID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 统计该相册的照片数量
	var count int64
	if err := tx.Model(&models.Photo{}).Where("album_id = ?", albumID).Count(&count).Error; err != nil {
		tx.Rollback()
		return err
	}
	album.PhotoCount = int(count)

	if err := tx.Save(album).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}
