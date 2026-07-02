package service

import (
	"errors"
	"kira-go/internal/database"
	"kira-go/internal/models"
	"kira-go/internal/schemas"

	"encoding/json"

	"gorm.io/gorm"
)

// projectDict 内部数据字典，对应 Python 的 _to_dict 返回的 dict
// Service 层不依赖 schemas 包，只在 API 层转换为响应
type projectDict struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	Slug        string   `json:"slug"`
	Description string   `json:"description"`
	LongDesc    string   `json:"long_description"`
	CoverImage  string   `json:"cover_image"`
	TechStack   []string `json:"tech_stack"`
	LinkGithub  string   `json:"link_github"`
	LinkGitee   string   `json:"link_gitee"`
	LinkLive    string   `json:"link_live"`
	LinkDocs    string   `json:"link_docs"`
	Status      string   `json:"status"`
	StatusLabel string   `json:"status_label"`
	IsFeatured  bool     `json:"is_featured"`
	Sort        int      `json:"sort"`
	CreatedAt   string   `json:"created_at"`
}

// GetProjects 获取所有项目列表（按 sort 排序）
func GetProjects() ([]projectDict, error) {
	db := database.GetDB()
	var projects []*models.Project

	if err := db.Order("sort").Find(&projects).Error; err != nil {
		return nil, err
	}

	result := make([]projectDict, 0, len(projects))
	for _, p := range projects {
		result = append(result, projectToDict(p))
	}
	return result, nil
}

// GetProjectBySlug 根据 slug 获取项目
func GetProjectBySlug(slug string) (*projectDict, error) {
	db := database.GetDB()

	p := &models.Project{}
	if err := db.Where("slug = ?", slug).First(p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("项目不存在")
		}
		return nil, err
	}

	result := projectToDict(p)
	return &result, nil
}

// CreateProject 创建项目
func CreateProject(data schemas.ProjectCreate) (*projectDict, error) {
	db := database.GetDB()

	// 序列化 tech_stack 为 JSON 字符串
	techStackJSON, err := json.Marshal(data.TechStack)
	if err != nil {
		return nil, err
	}

	p := &models.Project{
		Name:        data.Name,
		Slug:        data.Slug,
		Description: data.Description,
		LongDesc:    data.LongDesc,
		CoverImage:  data.CoverImage,
		TechStack:   string(techStackJSON),
		LinkGithub:  data.LinkGithub,
		LinkGitee:   data.LinkGitee,
		LinkLive:    data.LinkLive,
		LinkDocs:    data.LinkDocs,
		Status:      data.Status,
		StatusLabel: data.StatusLabel,
		IsFeatured:  data.IsFeatured,
		Sort:        data.Sort,
	}

	if err := db.Create(p).Error; err != nil {
		return nil, err
	}

	result := projectToDict(p)
	return &result, nil
}

// UpdateProject 更新项目
func UpdateProject(projectID uint, data schemas.ProjectUpdate) (*projectDict, error) {
	db := database.GetDB()

	// 查询现有项目
	p := &models.Project{}
	if err := db.First(p, projectID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("项目不存在")
		}
		return nil, err
	}

	// 直接赋值更新（仅更新非空值，对应 Python: for k, v in data.model_dump(exclude_unset=True).items(): setattr(p, k, v)）
	if data.Name != nil {
		p.Name = *data.Name
	}
	if data.Slug != nil {
		p.Slug = *data.Slug
	}
	if data.Description != nil {
		p.Description = *data.Description
	}
	if data.LongDesc != nil {
		p.LongDesc = *data.LongDesc
	}
	if data.CoverImage != nil {
		p.CoverImage = *data.CoverImage
	}
	if data.LinkGithub != nil {
		p.LinkGithub = *data.LinkGithub
	}
	if data.LinkGitee != nil {
		p.LinkGitee = *data.LinkGitee
	}
	if data.LinkLive != nil {
		p.LinkLive = *data.LinkLive
	}
	if data.LinkDocs != nil {
		p.LinkDocs = *data.LinkDocs
	}
	if data.Status != nil {
		p.Status = *data.Status
	}
	if data.StatusLabel != nil {
		p.StatusLabel = *data.StatusLabel
	}
	if data.IsFeatured != nil {
		p.IsFeatured = *data.IsFeatured
	}
	if data.Sort != nil {
		p.Sort = *data.Sort
	}
	// 特殊处理 tech_stack：需要 JSON 序列化
	if data.TechStack != nil {
		techStackJSON, err := json.Marshal(data.TechStack)
		if err != nil {
			return nil, err
		}
		p.TechStack = string(techStackJSON)
	}

	// 保存（GORM Save 自动更新 UpdatedAt）
	if err := db.Save(p).Error; err != nil {
		return nil, err
	}

	result := projectToDict(p)
	return &result, nil
}

// DeleteProject 删除项目
func DeleteProject(projectID uint) error {
	db := database.GetDB()

	p := &models.Project{}
	if err := db.First(p, projectID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("项目不存在")
		}
		return err
	}

	if err := db.Delete(p).Error; err != nil {
		return err
	}
	return nil
}

// projectToDict 将 Project 对象转为字典
func projectToDict(p *models.Project) projectDict {
	// 反序列化 tech_stack
	var techStack []string
	if p.TechStack != "" {
		json.Unmarshal([]byte(p.TechStack), &techStack)
	}

	return projectDict{
		ID:          p.ID,
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		LongDesc:    p.LongDesc,
		CoverImage:  p.CoverImage,
		TechStack:   techStack,
		LinkGithub:  p.LinkGithub,
		LinkGitee:   p.LinkGitee,
		LinkLive:    p.LinkLive,
		LinkDocs:    p.LinkDocs,
		Status:      p.Status,
		StatusLabel: p.StatusLabel,
		IsFeatured:  p.IsFeatured,
		Sort:        p.Sort,
		CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
