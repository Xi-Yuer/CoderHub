package repository

import (
	"coderhub/model"
	"coderhub/shared/storage"
	"errors"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type ArticlePVRepository interface {
	CreateArticlePV(articlePV *model.ArticlePV) error
	GetArticlePVByArticleID(articleID int64) (*model.ArticlePV, error)
	GetArticlePVsByArticleIDs(articleIDs []int64) ([]*model.ArticlePV, error)
	GetAllArticlePVByArticleIDs(articleIDs []int64) (int64, error)
	SyncIncrementalPVToDB() error
}

type ArticlePVRepositoryImpl struct {
	DB            *gorm.DB
	Redis         storage.RedisDB
	lastSyncTime  time.Time
	syncThreshold int64
}

func NewArticlePVRepositoryImpl(db *gorm.DB, rdb storage.RedisDB) *ArticlePVRepositoryImpl {
	return &ArticlePVRepositoryImpl{
		DB:            db,
		Redis:         rdb,
		lastSyncTime:  time.Now(),
		syncThreshold: 100,
	}
}

func (r *ArticlePVRepositoryImpl) CreateArticlePV(articlePV *model.ArticlePV) error {
	articleIDStr := fmt.Sprintf("%d", articlePV.ArticleID)

	// 增加PV计数
	newCount, err := r.Redis.HIncrBy("article:pv", articleIDStr, 1)
	if err != nil {
		return err
	}

	// 将文章ID添加到待同步集合
	err = r.Redis.SAdd("article:pv:need_sync", articleIDStr)
	if err != nil {
		return err
	}

	// 检查是否需要触发同步
	shouldSync := false

	// 尝试获取分布式锁
	lockKey := "article:pv:sync:lock"
	// 设置锁的过期时间为30秒
	locked, err := r.Redis.SetNX(lockKey, "1", 30*time.Second)
	if err != nil {
		return err
	}

	if locked {
		if newCount%r.syncThreshold == 0 {
			shouldSync = true
		}

		// 获取上次同步时间
		lastSyncStr, err := r.Redis.Get("article:pv:last_sync_time")
		if err == nil && lastSyncStr != "" {
			lastSyncTime, err := time.Parse(time.RFC3339, lastSyncStr)
			if err == nil && time.Since(lastSyncTime) > 5*time.Minute {
				shouldSync = true
			}
		}

		if shouldSync {
			go func() {
				// 在goroutine中处理锁的续期
				stopChan := make(chan struct{})
				go func() {
					ticker := time.NewTicker(10 * time.Second)
					defer ticker.Stop()
					for {
						select {
						case <-ticker.C:
							// 每10秒续期一次，将过期时间重置为30秒
							_ = r.Redis.Expire(lockKey, 30*time.Second)
						case <-stopChan:
							return
						}
					}
				}()

				// 执行同步
				err := r.SyncIncrementalPVToDB()
				// 同步完成后，关闭续期goroutine并释放锁
				close(stopChan)
				_ = r.Redis.Del(lockKey)

				if err == nil {
					_ = r.Redis.Set("article:pv:last_sync_time", time.Now().Format(time.RFC3339))
				}
			}()
		}
	}

	return nil
}

func (r *ArticlePVRepositoryImpl) GetArticlePVByArticleID(articleID int64) (*model.ArticlePV, error) {
	count, err := r.Redis.HGet("article:pv", fmt.Sprintf("%d", articleID))
	if err == nil && count != "" {
		countInt, err1 := strconv.ParseInt(count, 10, 64)
		if err1 == nil {
			return &model.ArticlePV{ArticleID: articleID, Count: countInt}, nil
		}
	}

	var articlePV model.ArticlePV
	err = r.DB.Where("article_id = ?", articleID).First(&articlePV).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &model.ArticlePV{ArticleID: articleID, Count: 0}, nil
		}
		return nil, err
	}
	return &articlePV, nil
}

func (r *ArticlePVRepositoryImpl) GetArticlePVsByArticleIDs(articleIDs []int64) ([]*model.ArticlePV, error) {
	// Step 1: 批量从 Redis 获取 PV 数据
	keys := make([]string, len(articleIDs))
	for i, articleID := range articleIDs {
		keys[i] = fmt.Sprintf("%d", articleID)
	}
	// 批量获取
	counts, err := r.Redis.HMGet("article:pv", keys...)
	if err != nil {
		return nil, err
	}

	// Step 2: 解析 Redis 数据并构建结果
	result := make([]*model.ArticlePV, len(articleIDs))
	var remainingIDs []int64
	for i, articleID := range articleIDs {
		// 检查 Redis 是否有该文章的 PV 数据
		if count, ok := counts[i].(string); ok && count != "" {
			countInt, err := strconv.ParseInt(count, 10, 64)
			if err == nil {
				result[i] = &model.ArticlePV{ArticleID: articleID, Count: countInt}
			} else {
				remainingIDs = append(remainingIDs, articleID)
			}
		} else {
			remainingIDs = append(remainingIDs, articleID)
		}
	}

	// Step 3: 对缺失的部分从数据库查询
	if len(remainingIDs) > 0 {
		var articlePVs []model.ArticlePV
		err := r.DB.Where("article_id IN ?", remainingIDs).Find(&articlePVs).Error
		if err != nil {
			return nil, err
		}

		// 将查询结果填充到最终结果
		for _, articlePV := range articlePVs {
			for i, articleID := range articleIDs {
				if articlePV.ArticleID == articleID {
					result[i] = &articlePV
					// 同时将数据库查询到的数据缓存到 Redis
					err := r.Redis.HSet("article:pv", fmt.Sprintf("%d", articleID), articlePV.Count)
					if err != nil {
						return nil, err
					}
				}
			}
		}
	}

	// 填充默认值，若数据库中没有数据
	for i, articleID := range articleIDs {
		if result[i] == nil {
			result[i] = &model.ArticlePV{ArticleID: articleID, Count: 0}
		}
	}

	return result, nil
}

func (r *ArticlePVRepositoryImpl) SyncIncrementalPVToDB() error {
	// 获取需要同步的文章ID集合
	needSyncSet, err := r.Redis.SMembers("article:pv:need_sync")
	if err != nil {
		return err
	}

	if len(needSyncSet) == 0 {
		return nil
	}

	// 批量获取这些文章的PV数据
	for _, articleIDStr := range needSyncSet {
		articleID, err := strconv.ParseInt(articleIDStr, 10, 64)
		if err != nil {
			continue
		}

		// 获取PV数据
		countStr, err := r.Redis.HGet("article:pv", articleIDStr)
		if err != nil || countStr == "" {
			continue
		}

		count, err := strconv.ParseInt(countStr, 10, 64)
		if err != nil {
			continue
		}

		// 更新数据库
		err = r.DB.Model(&model.ArticlePV{}).Where("article_id = ?", articleID).Update("count", count).Error
		if err != nil {
			continue
		}

		// 同步成功后从待同步集合中删除
		_ = r.Redis.SRem("article:pv:need_sync", articleIDStr)
	}

	return nil
}

// GetAllArticlePVByArticleIDs 获取所有文章的PV数据
func (r *ArticlePVRepositoryImpl) GetAllArticlePVByArticleIDs(articleIDs []int64) (int64, error) {
	var totalPV int64
	for _, articleID := range articleIDs {
		countStr, err := r.Redis.HGet("article:pv", fmt.Sprintf("%d", articleID))
		if err != nil {
			return 0, err
		}

		count, err := strconv.ParseInt(countStr, 10, 64)
		if err != nil {
			return 0, err
		}

		totalPV += count
	}
	return totalPV, nil
}
