package utils

import (
	"errors"
	"strconv"
	"strings"
)

type Validator struct {
	errors []string
}

func NewValidator() *Validator {
	return &Validator{errors: make([]string, 0)}
}

func (v *Validator) Clear() {
	v.errors = make([]string, 0)
}

func (v *Validator) Id(id int64) *Validator {
	if !IdRegex.MatchString(strconv.FormatInt(id, 10)) {
		v.errors = append(v.errors, "ID格式错误：需要为数字")
	}
	return v
}

func (v *Validator) Username(username string) *Validator {
	if !UsernameRegex.MatchString(username) {
		v.errors = append(v.errors, "用户名格式错误：需要3-32位字母、数字或下划线")
	}
	return v
}

func (v *Validator) Password(password string) *Validator {
	if !PasswordRegex.MatchString(password) {
		v.errors = append(v.errors, "密码格式错误：需要6-32位字母或数字")
	}
	return v
}

func (v *Validator) Email(email string) *Validator {
	if !EmailRegex.MatchString(email) {
		v.errors = append(v.errors, "邮箱格式错误：需要符合邮箱格式")
	}
	return v
}

func (v *Validator) Token(token string) *Validator {
	if !TokenRegex.MatchString(token) {
		v.errors = append(v.errors, "token格式错误：需要符合token格式")
	}
	return v
}
func (v *Validator) ConfirmPassword(password string, confirmPassword string) *Validator {
	if password != confirmPassword {
		v.errors = append(v.errors, "确认密码与密码不一致")
	}
	return v
}

func (v *Validator) Title(title string) *Validator {
	if !TitleRegex.MatchString(title) {
		v.errors = append(v.errors, "标题格式错误：需要1-255位")
	}
	return v
}

func (v *Validator) Summary(summary string) *Validator {
	if !SummaryRegex.MatchString(summary) {
		v.errors = append(v.errors, "摘要格式错误：需要1-1024")
	}
	return v
}

func (v *Validator) Content(content string) *Validator {
	if !ContentRegex.MatchString(content) {
		v.errors = append(v.errors, "内容格式错误：需要1-65535位")
	}
	return v
}

func (v *Validator) Tags(tags []string) *Validator {
	for _, tag := range tags {
		if !TagsRegex.MatchString(tag) {
			v.errors = append(v.errors, "标签格式错误：需要1-255位")
		}
	}
	return v
}

func (v *Validator) ArticleType(articleType string) *Validator {
	if articleType != "article" && articleType != "micro_post" {
		v.errors = append(v.errors, "文章类型格式错误：需要为article或micro_post")
	}
	return v
}

func (v *Validator) ArticleStatus(articleStatus []string) *Validator {
	for _, status := range articleStatus {
		if status != "draft" && status != "published" {
			v.errors = append(v.errors, "文章状态格式错误：需要为draft或published")
		}
	}
	return v
}

func (v *Validator) ArticleCoverImage(articleCoverImage string) *Validator {
	if !ArticleCoverImageRegex.MatchString(articleCoverImage) {
		v.errors = append(v.errors, "封面图片格式错误：需要符合URL格式")
	}
	return v
}

func (v *Validator) AuthorID(authorID int64) *Validator {
	if !AuthorIDRegex.MatchString(strconv.FormatInt(authorID, 10)) {
		v.errors = append(v.errors, "作者ID格式错误：需要为数字")
	}
	return v
}

func (v *Validator) ArticleID(articleID int64) *Validator {
	if !ArticleIDRegex.MatchString(strconv.FormatInt(articleID, 10)) {
		v.errors = append(v.errors, "文章ID格式错误：需要为数字")
	}
	return v
}

func (v *Validator) Check() error {
	if len(v.errors) > 0 {
		defer v.Clear()
		return errors.New(strings.Join(v.errors, "; "))
	}
	return nil
}
