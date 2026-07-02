package service

import (
	"errors"
	"strings"
	"time"

	"kira-go/internal/database"
	"kira-go/internal/models"
	"kira-go/internal/schemas"

	"gorm.io/gorm"
)

// postDict 内部数据字典，对应 Python 的 _post_to_dict 返回的 dict
// Service 层不依赖 schemas 包，只在 API 层转换为 schemas.PostOut/schemas.PostDetail
type postDict struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Description string     `json:"description"`
	Cover       string     `json:"cover"`
	Category    string     `json:"category"`
	Tags        []string   `json:"tags"`
	Status      string     `json:"status"`
	IsPinned    bool       `json:"is_pinned"`
	Views       int        `json:"views"`
	Likes       int        `json:"likes"`
	WordCount   int        `json:"word_count"`
	ReadingTime int        `json:"reading_time"`
	PublishedAt *time.Time `json:"published_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Content     string     `json:"content,omitempty"` // 仅在详情时包含
}

// GetPosts 获取文章列表（分页、状态筛选、分类筛选、标签筛选）
func GetPosts(status, category, tag string, page, size int) ([]postDict, error) {
	db := database.GetDB()
	query := db.Model(&models.Post{}).Preload("Category").Preload("Tags")

	// status 筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 根据分类 slug 筛选（使用 Join 避免子查询）
	if category != "" {
		query = query.Joins("LEFT JOIN category ON category.id = post.category_id").
			Where("category.slug = ?", category)
	}

	// 根据标签 slug 筛选（使用 Join）
	if tag != "" {
		query = query.Joins("INNER JOIN post_tag ON post_tag.post_id = post.id").
			Joins("INNER JOIN tag ON tag.id = post_tag.tag_id").
			Where("tag.slug = ?", tag)
	}

	// 排序：置顶优先，然后按创建时间降序
	query = query.Order("is_pinned DESC, created_at DESC")

	// 分页
	offset := (page - 1) * size
	var posts []*models.Post // Post结构体较大使用指针切片
	if err := query.Offset(offset).Limit(size).Find(&posts).Error; err != nil {
		return nil, err
	}

	// 转换为 dict（不包含 content）
	result := make([]postDict, 0, len(posts))
	for _, post := range posts {
		result = append(result, postToDict(post))
	}
	return result, nil
}

// GetPostBySlug 根据 slug 获取文章（增加浏览量）
func GetPostBySlug(slug string) (*postDict, error) {
	db := database.GetDB()

	post := &models.Post{}
	if err := db.Preload("Category").Preload("Tags").
		Where("slug = ?", slug).
		First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在")
		}
		return nil, err
	}

	// 2. 原子更新 views（单条 SQL，无需事务）
	if err := db.Model(post).UpdateColumn("views", db.Raw("views + 1")).Error; err != nil {
		return nil, err
	}
	post.Views++ // 同步内存值

	result := postToDict(post)
	return &result, nil
}

// GetPostByID 根据 ID 获取文章
func GetPostByID(postID uint) (*postDict, error) {
	db := database.GetDB()

	post := &models.Post{}
	if err := db.Preload("Category").Preload("Tags").First(post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在")
		}
		return nil, err
	}

	result := postToDict(post)
	return &result, nil
}

// CreatePost 创建文章
func CreatePost(data schemas.PostCreate) (*postDict, error) {
	db := database.GetDB()

	// 开始事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建 Post 对象（排除 tags）
	post := &models.Post{
		Title:       data.Title,
		Slug:        data.Slug,
		Description: data.Description,
		Content:     data.Content,
		Cover:       data.Cover,
		CategoryID:  data.CategoryID,
		Status:      data.Status,
		IsPinned:    data.IsPinned,
		WordCount:   data.WordCount,
		ReadingTime: data.ReadingTime,
	}

	// 自动计算字数和阅读时间（仅当未手动指定时）
	if post.Content != "" {
		if post.WordCount == 0 {
			post.WordCount = len(post.Content)
		}
		if post.ReadingTime == 0 {
			post.ReadingTime = max(1, post.WordCount/300) // 300 字符/分钟
		}
	}

	// 如果状态为 published 且未设置发布时间，设为当前时间
	if post.Status == "published" && post.PublishedAt == nil {
		now := time.Now()
		post.PublishedAt = &now
	}

	// 插入 Post（执行 SQL，获取自增 ID）
	if err := tx.Create(post).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	// 同步标签（删除旧的，创建新的）
	if len(data.Tags) > 0 {
		if err := syncTags(tx, post.ID, data.Tags); err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	// 更新分类计数
	if post.CategoryID != nil {
		if err := updateCategoryCount(tx, *post.CategoryID); err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	// 重新加载 Post（带关联）
	if err := db.Preload("Category").Preload("Tags").First(post, post.ID).Error; err != nil {
		return nil, err
	}
	// 返回 dict
	result := postToDict(post)
	return &result, nil
}

// UpdatePost 更新文章
func UpdatePost(postID uint, data schemas.PostUpdate) (*postDict, error) {
	db := database.GetDB()

	// 开始事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 查询现有 Post
	post := &models.Post{}
	if err := tx.First(post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在")
		}
		return nil, err
	}

	oldCategoryID := post.CategoryID

	// 2. 直接赋值更新 post 对象（对应 Python: for k, v in update_data.items(): setattr(post, k, v)）
	if data.Title != nil {
		post.Title = *data.Title
	}
	if data.Slug != nil {
		post.Slug = *data.Slug
	}
	if data.Description != nil {
		post.Description = *data.Description
	}
	if data.Cover != nil {
		post.Cover = *data.Cover
	}
	if data.CategoryID != nil {
		post.CategoryID = data.CategoryID
	}
	if data.IsPinned != nil {
		post.IsPinned = *data.IsPinned
	}
	if data.Status != nil {
		post.Status = *data.Status
	}
	if data.Content != nil {
		post.Content = *data.Content
	}
	if data.WordCount != nil {
		post.WordCount = *data.WordCount
	}
	if data.ReadingTime != nil {
		post.ReadingTime = *data.ReadingTime
	}

	// 3. 自动计算（对应 Python: if post.content: ...）
	// 如果 post 有内容（此时 post.Content 已经是新值）
	if post.Content != "" {
		// 如果用户没有手动指定 word_count
		if data.WordCount == nil {
			post.WordCount = len(post.Content)
		}
		// 如果用户没有手动指定 reading_time
		if data.ReadingTime == nil {
			post.ReadingTime = max(1, post.WordCount/300)
		}
	}

	// 4. 处理 published_at（对应 Python: if post.status == "published" ...）
	if post.Status == "published" && post.PublishedAt == nil {
		now := time.Now()
		post.PublishedAt = &now
	}

	// 5. 一次性更新（GORM 会自动跳过 nil 指针）
	if err := tx.Model(post).Updates(post).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 6. 同步标签（如果提供了 tags）
	if data.Tags != nil {
		if err := syncTags(tx, post.ID, *data.Tags); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 7. 更新分类计数（旧分类和新分类都更新）
	if oldCategoryID != nil {
		if err := updateCategoryCount(tx, *oldCategoryID); err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	if post.CategoryID != nil {
		if err := updateCategoryCount(tx, *post.CategoryID); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 8. 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// 9. 重新加载 Post 及其关联
	if err := db.Preload("Category").Preload("Tags").First(&post, postID).Error; err != nil {
		return nil, err
	}

	// 10. 返回 dict
	result := postToDict(post)
	return &result, nil
}

// DeletePost 删除文章
func DeletePost(postID uint) error {
	db := database.GetDB()

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 查询 Post
	post := &models.Post{}
	if err := tx.First(post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文章不存在")
		}
		return err
	}

	catID := post.CategoryID

	// 删除 Post
	if err := tx.Delete(post).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 更新分类计数
	if catID != nil {
		if err := updateCategoryCount(tx, *catID); err != nil {
			tx.Rollback()
			return err
		}
	}

	// 更新所有标签计数
	if err := updateTagCounts(tx); err != nil {
		tx.Rollback()
		return err
	}

	// 提交
	return tx.Commit().Error
}

// CountPosts 统计文章数量
func CountPosts(status string) (int64, error) {
	db := database.GetDB()

	var count int64
	query := db.Model(&models.Post{}).Where("1=1")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ToggleLike 点赞/取消点赞
func ToggleLike(postID uint, unlike bool) (int, error) {
	db := database.GetDB()

	post := &models.Post{}
	if err := db.First(post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("文章不存在")
		}
		return 0, err
	}

	// 更新点赞数
	if unlike {
		post.Likes = max(0, post.Likes-1)
	} else {
		post.Likes += 1
	}

	if err := db.Save(post).Error; err != nil {
		return 0, err
	}

	return post.Likes, nil
}

// ========================== 辅助方法 ==========================

// syncTags 同步文章标签：删除旧关联，创建新标签（如不存在），建立关联。
func syncTags(tx *gorm.DB, postID uint, tagNames []string) error {
	// 删除旧关联
	if err := tx.Where("post_id = ?", postID).Delete(&models.PostTag{}).Error; err != nil {
		return err
	}
	for _, name := range tagNames {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		// slug: 小写，空格替换为连字符
		slug := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
		// 查找或创建标签
		tag := models.Tag{}
		if err := tx.Where("name = ?", name).First(tag).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 创建新标签
				tag = models.Tag{Name: name, Slug: slug}
				if err := tx.Create(tag).Error; err != nil {
					return err
				}
			} else {
				// 出现找不到之外的其他err
				return err
			}
		}
		// 更新PostTag表项
		if err := tx.Create(&models.PostTag{PostID: postID, TagID: tag.ID}).Error; err != nil {
			return err
		}
	}

	// 全量更新标签计数
	return updateTagCounts(tx)
}

// updateTagCounts 更新所有的 tag 计数，一般在 syncTags 更新完一批 tag 之后执行。
func updateTagCounts(tx *gorm.DB) error {
	// 查询所有标签
	var tags []*models.Tag // 同一风格使用指针切片
	if err := tx.Find(&tags).Error; err != nil {
		return err
	}

	for _, tag := range tags {
		// 统计该标签的文章数量
		var count int64
		if err := tx.Model(&models.PostTag{}).Where("tag_id = ?", tag.ID).Count(&count).Error; err != nil {
			return err
		}
		tag.PostCount = int(count)
		if err := tx.Save(tag).Error; err != nil {
			return err
		}
	}
	return nil
}

// updateCategoryCount 更新 category_id 的 post_count
func updateCategoryCount(tx *gorm.DB, categoryID uint) error {
	if categoryID == 0 {
		return nil
	}

	// 统计该分类的文章数量
	var count int64
	if err := tx.Model(&models.Post{}).Where("category_id = ?", categoryID).Count(&count).Error; err != nil {
		return err
	}

	// 更新分类
	return tx.Model(&models.Category{}).Where("id = ?", categoryID).Update("post_count", count).Error
}

// postToDict 将 Post 对象转为带 category 和 tags 的字典
// 注意：Post 需已 Preload Category 和 Tags
func postToDict(post *models.Post) postDict {
	// 分类名称
	catName := ""
	if post.Category != nil {
		catName = post.Category.Name
	}

	// 标签名称列表（已 Preload）
	tagNames := make([]string, 0, len(post.Tags))
	for _, t := range post.Tags {
		tagNames = append(tagNames, t.Name)
	}

	return postDict{
		ID:          post.ID,
		Title:       post.Title,
		Slug:        post.Slug,
		Description: post.Description,
		Cover:       post.Cover,
		Category:    catName,
		Tags:        tagNames,
		Status:      post.Status,
		IsPinned:    post.IsPinned,
		Views:       post.Views,
		Likes:       post.Likes,
		WordCount:   post.WordCount,
		ReadingTime: post.ReadingTime,
		PublishedAt: post.PublishedAt,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		Content:     post.Content,
	}
}
