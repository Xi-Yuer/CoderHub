package model

import "gorm.io/gorm"

type SandpackProjects struct {
	gorm.Model
	Name        string `gorm:"column:name;not null"`                                                                                                                                                                                                                                                                          // 项目名称
	Description string `gorm:"column:description;default:'';"`                                                                                                                                                                                                                                                                // 项目描述
	Template    string `gorm:"column:template;type:enum('static','angular','vanilla','react','react-ts','solid','svelte','vue','vue-ts','node','nextjs','test-ts','vanilla-ts','vite','vite-react','vite-react-ts','vite-preact','vite-preact-ts','vite-vue','vite-vue-ts','vite-svelte','vite-svelte-ts','astro');not null"` // 模板类型
	UserID      int64  `gorm:"column:user_id;not null"`
	ArticleID   int64  `gorm:"column:article_id;not null"`
	_           struct {
		UserID    int64
		ArticleID int64
	} `gorm:"uniqueIndex:idx_user_id_article_id"` // 唯一索引
}

var ValidSandpackTemplates = []string{
	"static", "angular", "vanilla", "react", "react-ts",
	"solid", "svelte", "vue", "vue-ts", "node", "nextjs",
	"test-ts", "vanilla-ts", "vite", "vite-react", "vite-react-ts",
	"vite-preact", "vite-preact-ts", "vite-vue", "vite-vue-ts",
	"vite-svelte", "vite-svelte-ts", "astro",
}
