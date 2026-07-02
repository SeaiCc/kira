package schemas

// AlbumCreate 创建相册请求
type AlbumCreate struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Cover       string `json:"cover"`
	Sort        int    `json:"sort"`
}

// AlbumUpdate 更新相册请求
type AlbumUpdate struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Cover       *string `json:"cover"`
	Sort        *int    `json:"sort"`
}

// PhotoCreate 创建照片请求
type PhotoCreate struct {
	AlbumID     uint   `json:"album_id" binding:"required"`
	URL         string `json:"url" binding:"required"`
	Caption     string `json:"caption"`
	Orientation string `json:"orientation"` // landscape | portrait | square
	Sort        int    `json:"sort"`
}
