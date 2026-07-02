package schemas

// ProjectCreate 创建项目请求
// 对应 Python: app/schemas/project.py -> ProjectCreate
type ProjectCreate struct {
	Name          string   `json:"name" binding:"required"`
	Slug          string   `json:"slug" binding:"required"`
	Description   string   `json:"description"`
	LongDesc      string   `json:"long_description"`
	CoverImage    string   `json:"cover_image"`
	TechStack     []string `json:"tech_stack"`
	LinkGithub    string   `json:"link_github"`
	LinkGitee     string   `json:"link_gitee"`
	LinkLive      string   `json:"link_live"`
	LinkDocs      string   `json:"link_docs"`
	Status        string   `json:"status"`
	StatusLabel   string   `json:"status_label"`
	IsFeatured    bool     `json:"is_featured"`
	Sort          int      `json:"sort"`
}

// ProjectUpdate 更新项目请求
// 对应 Python: app/schemas/project.py -> ProjectUpdate
type ProjectUpdate struct {
	Name          *string  `json:"name"`
	Slug          *string  `json:"slug"`
	Description   *string  `json:"description"`
	LongDesc      *string  `json:"long_description"`
	CoverImage    *string  `json:"cover_image"`
	TechStack     *[]string `json:"tech_stack"`
	LinkGithub    *string  `json:"link_github"`
	LinkGitee     *string  `json:"link_gitee"`
	LinkLive      *string  `json:"link_live"`
	LinkDocs      *string  `json:"link_docs"`
	Status        *string  `json:"status"`
	StatusLabel   *string  `json:"status_label"`
	IsFeatured    *bool    `json:"is_featured"`
	Sort          *int     `json:"sort"`
}

// SiteConfigUpdate 更新站点配置请求
// 对应 Python: app/schemas/site_config.py -> SiteConfigUpdate
type SiteConfigUpdate struct {
	Value       string `json:"value" binding:"required"`
	Description string `json:"description"`
}
