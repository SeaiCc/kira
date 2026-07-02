package schemas

// CategoryCreate 创建分类请求
// 对应 Python: app/schemas/category.py -> CategoryCreate
type CategoryCreate struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Description string `json:"description"`
	Sort        int    `json:"sort"`
}

// CategoryUpdate 更新分类请求
// 对应 Python: app/schemas/category.py -> CategoryUpdate
type CategoryUpdate struct {
	Name        *string `json:"name"`
	Slug        *string `json:"slug"`
	Description *string `json:"description"`
	Sort        *int    `json:"sort"`
}

// TagCreate 创建标签请求
// 对应 Python: app/schemas/tag.py -> TagCreate
type TagCreate struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}

// TagUpdate 更新标签请求
// 对应 Python: app/schemas/tag.py -> TagUpdate
type TagUpdate struct {
	Name *string `json:"name"`
	Slug *string `json:"slug"`
}
