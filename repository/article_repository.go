package repository

import (
	"coderhub/model"
	"coderhub/shared/storage"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"gorm.io/gorm"
)

type ArticleRepository interface {
	CreateArticle(article *model.Articles) error
	GetArticleByID(id int64) (*model.Articles, error)
	GetArticlesByIDs(ids []int64) ([]*model.Articles, error)
	ListRecommendedArticles(type_ string, categoryID int64, page, pageSize int64) ([]int64, error)
	ListArticlesByAuthor(authorID int64, _type string, page, pageSize int64) ([]int64, int64, error)
	BatchGetArticle(ids []int64) ([]*model.ArticlePreviewWithAuthInfo, error)
	UpdateArticle(article *model.Articles) error
	DeleteArticle(id int64) error
	GetUserArticleCount(authorID int64) (int64, error)
	GetUserMicroPostCount(authorID int64) (int64, error)
	GetUserAllArticleIDS(authorID int64) ([]int64, error)
	GetArticlesBySearchKeys(keys string, _type string, page, pageSize int64) ([]int64, error)
}
type ArticleRepositoryImpl struct {
	DB       *gorm.DB
	Redis    storage.RedisDB
	Elastic  *storage.ElasticSearchClient
	minLikes int32
}

func NewArticleRepositoryImpl(db *gorm.DB, rdb storage.RedisDB) *ArticleRepositoryImpl {
	return &ArticleRepositoryImpl{
		DB:    db,
		Redis: rdb,
		//Elastic:  elastic,
		minLikes: 10,
	}
}

type ArticleEsVO struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Content string `json:"content"`
	Tags    string `json:"tags"`
}

func (r *ArticleRepositoryImpl) CreateArticle(article *model.Articles) error {
	if err := r.DB.Create(article).Error; err != nil {
		return err
	}
	//articleEsVO := ArticleEsVO{
	//	ID:      article.ID,
	//	Title:   article.Title,
	//	Summary: article.Summary,
	//	Content: article.Content,
	//	Tags:    article.Tags,
	//}
	//articleJSON, err := json.Marshal(articleEsVO)
	//if err != nil {
	//	fmt.Printf("Failed to marshal article to JSON: %v\n", err)
	//	return err
	//}
	//// 使用 bytes.NewReader 将字节切片转换为 io.Reader
	//err = r.Elastic.CreateIndex(article.Type, article.ID, bytes.NewReader(articleJSON))
	//if err != nil {
	//	return err
	//}
	// 创建后设置缓存
	return r.setCache(article.CacheKeyByID(article.ID), article)
}

func (r *ArticleRepositoryImpl) GetArticleByID(id int64) (*model.Articles, error) {
	var article model.Articles
	key := article.CacheKeyByID(id)

	// 尝试从缓存获取
	if cached, err := r.getCache(key); err == nil {
		return cached, nil
	}

	// 简化数据库查询
	if err := r.DB.First(&article, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("文章不存在: %v", id)
		}
		return nil, err
	}

	// 异步设置缓存，避免影响主流程
	go func() {
		_ = r.setCache(key, &article)
	}()

	return &article, nil
}

func (r *ArticleRepositoryImpl) GetArticlesByIDs(ids []int64) ([]*model.Articles, error) {
	var articles []*model.Articles
	if err := r.DB.Where("id IN ?", ids).Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

func (r *ArticleRepositoryImpl) ListRecommendedArticles(type_ string, categoryID int64, page, pageSize int64) ([]int64, error) {
	var ids []int64

	query := r.DB.Table("articles").Where("type = ?", type_)

	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	if err := query.
		Order("created_at DESC").
		Limit(int(pageSize)).
		Offset(int((page-1)*pageSize)).
		Pluck("id", &ids).Error; err != nil {
		return nil, err
	}

	return ids, nil
}
func (r *ArticleRepositoryImpl) ListArticlesByAuthor(authorID int64, _type string, page, pageSize int64) ([]int64, int64, error) {
	var ids []int64
	var total int64
	if err := r.DB.Table("articles").
		Where("author_id = ? AND type = ? AND deleted_at IS NULL", authorID, _type).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.DB.Table("articles").
		Where("author_id = ? AND type = ? AND deleted_at IS NULL", authorID, _type).
		Order("created_at DESC").
		Limit(int(pageSize)).
		Offset(int((page-1)*pageSize)).
		Pluck("id", &ids).Error; err != nil {
		return nil, 0, err
	}
	return ids, total, nil
}

func (r *ArticleRepositoryImpl) BatchGetArticle(ids []int64) ([]*model.ArticlePreviewWithAuthInfo, error) {
	// 查找文章并且获取文章的创建者
	var articles []*model.ArticlePreviewWithAuthInfo
	err := r.DB.Table("articles AS a").
		Select(`
        a.id AS article_id, 
        a.title, 
        img.url AS cover_image, 
        a.summary, 
        a.created_at AS create_time, 
        u.id AS author_id, 
        u.user_name AS auth_name, 
        u.avatar
    `).
		Joins("JOIN users u ON a.author_id = u.id").
		Joins("JOIN image_relations ir ON a.id = ir.entity_id AND ir.entity_type = ?", "article_cover").
		Joins("JOIN images img ON ir.image_id = img.id").
		Where("a.id IN ?", ids).
		Scan(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (r *ArticleRepositoryImpl) UpdateArticle(article *model.Articles) error {
	if article.ID <= 0 {
		return errors.New("无效的文章ID")
	}

	// 使用事务确保数据一致性
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// 检查文章是否存在
		var count int64
		if err := tx.Model(&model.Articles{}).Select("id").Where("id = ?", article.ID).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return errors.New("文章不存在")
		}

		// 更新数据库
		if err := tx.Model(article).Updates(article).Error; err != nil {
			return err
		}

		// 删除旧缓存并设置新缓存
		key := article.CacheKeyByID(article.ID)
		if err := r.delCache(key); err != nil {
			return err
		}
		return r.setCache(key, article)
	})
}

func (r *ArticleRepositoryImpl) DeleteArticle(id int64) error {
	article, err := r.GetArticleByID(id)
	if err != nil {
		return err
	}
	err = r.delCache(article.CacheKeyByID(id))
	if err != nil {
		return err
	}
	// 删除ES索引
	//err = r.Elastic.DeleteByID(model.ArticleType, id)
	return r.DB.Delete(&model.Articles{}, id).Error
}

func (r *ArticleRepositoryImpl) getCache(key string) (*model.Articles, error) {
	var article model.Articles
	data, err := r.Redis.Get(key)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(data), &article); err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *ArticleRepositoryImpl) setCache(key string, article *model.Articles) error {
	data, err := json.Marshal(article)
	if err != nil {
		return err
	}
	return r.Redis.Set(key, string(data))
}

func (r *ArticleRepositoryImpl) delCache(key string) error {
	return r.Redis.Del(key)
}

func (r *ArticleRepositoryImpl) GetUserArticleCount(authorID int64) (int64, error) {
	var count int64
	if err := r.DB.Model(&model.Articles{}).Where("author_id = ? AND type = ? AND deleted_at IS NULL", authorID, model.ArticleType).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *ArticleRepositoryImpl) GetUserMicroPostCount(authorID int64) (int64, error) {
	var count int64
	if err := r.DB.Model(&model.Articles{}).Where("author_id = ? AND type = ? AND deleted_at IS NULL", authorID, model.MicroPostType).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
func (r *ArticleRepositoryImpl) GetUserAllArticleIDS(authorID int64) ([]int64, error) {
	var ids []int64
	if err := r.DB.Model(&model.Articles{}).Where("author_id = ?", authorID).Pluck("id", &ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}
func (r *ArticleRepositoryImpl) GetArticlesBySearchKeys(keys string, _type string, page, pageSize int64) ([]int64, error) {
	//fields, err := r.Elastic.SearchByFields(_type, keys, []string{"title", "summary", "content", "tags"}, page, pageSize)
	//if err != nil {
	//	return nil, err
	//}
	//return fields, nil

	// 服务降级使用
	var ids []int64

	// 解码查询关键词
	keys, _ = url.QueryUnescape(keys)

	// 将用户输入的关键词转为 Boolean 模式格式，如：Go 微服务 => +Go +微服务
	keywords := strings.Fields(keys)
	for i, kw := range keywords {
		keywords[i] = "+" + kw
	}
	keys = strings.Join(keywords, " ")

	// 执行全文索引查询，使用 BOOLEAN MODE
	// sql := `
	// 	SELECT id FROM articles
	// 	WHERE type = ?
	// 	  AND MATCH(title, summary, content, tags) AGAINST(? IN BOOLEAN MODE)
	// 	LIMIT ? OFFSET ?
	// `

	// if err := r.DB.Raw(sql, _type, keys, pageSize, (page-1)*pageSize).Scan(&ids).Error; err != nil {
	// 	return nil, err
	// }
	r.DB.Pluck("id", &ids).Model(&model.Articles{}).Where("type =? AND title LIKE ? OR summary LIKE ? OR content LIKE ? OR tags LIKE ?)", _type, "%"+keys, "%"+keys, "%"+keys, "%"+keys).Limit(int(pageSize)).Offset(int((page - 1) * pageSize))

	return ids, nil
}
