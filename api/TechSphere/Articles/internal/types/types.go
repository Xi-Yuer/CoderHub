// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.3

package types

type Article struct {
	Id           int64    `json:"id" form:"id"`                     // 主键 ID
	Type         string   `json:"type" form:"type"`                 // 内容类型：长文或短文
	Title        string   `json:"title" form:"title"`               // 标题
	Content      string   `json:"content" form:"content"`           // 内容
	Summary      string   `json:"summary" form:"summary"`           // 摘要
	ImageUrls    []string `json:"imageUrls" form:"imageUrls"`       // 图片 URL 列表
	CoverImage   string   `json:"coverImage" form:"coverImage"`     // 封面图片 URL
	AuthorId     int64    `json:"authorId" form:"authorId"`         // 作者 ID
	Tags         []string `json:"tags" form:"tags"`                 // 标签列表
	ViewCount    int64    `json:"viewCount" form:"viewCount"`       // 阅读次数
	LikeCount    int64    `json:"likeCount" form:"likeCount"`       // 点赞次数
	CommentCount int64    `json:"commentCount" form:"commentCount"` // 评论数
	Status       string   `json:"status" form:"status"`             // 文章状态
	CreatedAt    int64    `json:"createdAt" form:"createdAt"`       // 创建时间
	UpdatedAt    int64    `json:"updatedAt" form:"updatedAt"`       // 更新时间
}

type CreateArticleReq struct {
	Type         string   `json:"type,options=draft|published"`   // 内容类型
	Title        string   `json:"title"`                          // 标题
	Content      string   `json:"content"`                        // 内容
	Summary      string   `json:"summary"`                        // 摘要
	ImageIds     []string `json:"imageIds"`                       // 图片 URL 列表
	CoverImageID string   `json:"coverImageID"`                   // 封面图片 URL
	Tags         []string `json:"tags"`                           // 标签列表
	Status       string   `json:"status,options=draft|published"` // 文章状态
}

type CreateArticleResp struct {
	Response
	Data int64 `json:"data"` // 文章详情
}

type DeleteArticleReq struct {
	Id int64 `path:"id"` // 文章 ID
}

type DeleteArticleResp struct {
	Response
	Data bool `json:"data"` // 是否删除成功
}

type GetArticleReq struct {
	Id int64 `path:"id"` // 文章 ID
}

type GetArticleResp struct {
	Response
	Data *Article `json:"data"` // 文章详情
}

type HealthResp struct {
	Response
}

type Response struct {
	Code    int32  `json:"code"`    // 状态码
	Message string `json:"message"` // 提示信息
}

type UpdateArticleReq struct {
	Id           int64    `path:"id"`           // 文章 ID
	Title        string   `json:"title"`        // 标题
	Content      string   `json:"content"`      // 内容
	Summary      string   `json:"summary"`      // 摘要
	ImageIds     []string `json:"imageIds"`     // 图片 URL 列表
	CoverImageID string   `json:"coverImageID"` // 封面图片 URL
	Tags         []string `json:"tags"`         // 标签列表
	Status       string   `json:"status"`       // 文章状态
}

type UpdateArticleResp struct {
	Response
	Data bool `json:"data"` // 是否更新成功
}
