package utils

import "regexp"

var (
	// IdRegex 主键ID不能为空
	IdRegex = regexp.MustCompile(`^\d+$`)
	// UsernameRegex 用户名不能为空，长度为3-32，只能包含字母、数字和下划线
	UsernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,32}$`)
	// PasswordRegex 密码不能为空，长度为6-32，只能包含字母、数字
	PasswordRegex = regexp.MustCompile(`^[A-Za-z\d]{6,32}$`)
	// EmailRegex 邮箱不能为空，必须符合邮箱格式
	EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	// TitleRegex 标题不能为空，长度为1-100，只能包含字母、数字和空格
	TitleRegex = regexp.MustCompile(`.+`)
	// SummaryRegex 摘要不能为空，长度为1-1024
	SummaryRegex = regexp.MustCompile(`.+`)
	// ContentRegex 内容不能为空，长度为1-65535，只能包含字母、数字和空格
	ContentRegex = regexp.MustCompile(`.+`)
	// TagsRegex 标签不能为空，长度为1-10
	TagsRegex = regexp.MustCompile(`.+`)
	// ArticleIDRegex 文章ID不能为空
	ArticleIDRegex = regexp.MustCompile(`^\d+$`)
	// AuthorIDRegex 作者ID不能为空
	AuthorIDRegex = regexp.MustCompile(`^\d+$`)
	// ArticleCoverImageRegex 封面图片不能为空，必须符合URL格式
	ArticleCoverImageRegex = regexp.MustCompile(`^https?://\S+$`)
	// TokenRegex token不能为空
	TokenRegex = regexp.MustCompile(`^\d+$`)
)
