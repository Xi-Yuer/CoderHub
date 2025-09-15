package articleservicelogic

import (
	"coderhub/conf"
	"coderhub/model"
	imagerelationservicelogic "coderhub/rpc/coderhub/internal/logic/imagerelationservice"
	"coderhub/shared/utils"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateArticleLogic {
	return &UpdateArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateArticleLogic) UpdateArticle(in *coderhub.UpdateArticleRequest) (*coderhub.UpdateArticleResponse, error) {
	// 记录执行时间
	defer func(start time.Time) {
		l.Logger.Infof("UpdateArticle 执行耗时: %v", time.Since(start))
	}(time.Now())

	// 参数校验
	if err := l.validateArticleUpdate(in); err != nil {
		l.Logger.Errorf(conf.ErrValidationFailed, err)
		return nil, err
	}

	// 获取用户ID
	userID, err := l.getUserID()
	if err != nil {
		return nil, err
	}

	// 检查文章权限
	article, err := l.checkArticlePermission(in.Id, userID)
	if err != nil {
		return nil, err
	}

	// 更新文章信息
	l.updateArticleFields(article, in)

	// 保存更新
	if err := l.svcCtx.ArticleRepository.UpdateArticle(article); err != nil {
		l.Logger.Errorf(conf.ErrUpdateFailed, err)
		return nil, fmt.Errorf(conf.ErrUpdateFailed, err)
	}

	l.Logger.Infof("文章更新成功, ID: %d", in.Id)

	return &coderhub.UpdateArticleResponse{
		Success: true,
	}, nil
}

// validateArticleUpdate 验证文章更新请求
func (l *UpdateArticleLogic) validateArticleUpdate(req *coderhub.UpdateArticleRequest) error {
	// 基础字段验证
	if err := utils.NewValidator().
		ArticleID(req.Id).
		Title(req.Title).
		Summary(req.Summary).
		Content(req.Content).
		Tags(req.Tags).
		Check(); err != nil {
		return fmt.Errorf(conf.ErrValidationFailed, err)
	}

	// 验证图片数量
	if len(req.ImageIds) > conf.MaxImageCount {
		return fmt.Errorf(conf.ErrImageCountExceeded, conf.MaxImageCount)
	}

	return nil
}

// getUserID 获取并转换用户ID
func (l *UpdateArticleLogic) getUserID() (int64, error) {
	userId, err := utils.GetUserMetaData(l.ctx)
	if err != nil {
		l.Logger.Errorf(conf.ErrUserMetaFailed, err)
		return 0, fmt.Errorf(conf.ErrUserMetaFailed, err)
	}

	userIdInt, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		l.Logger.Errorf(conf.ErrUserIDConversion, err)
		return 0, fmt.Errorf(conf.ErrUserIDConversion, err)
	}

	return userIdInt, nil
}

// checkArticlePermission 检查文章权限
func (l *UpdateArticleLogic) checkArticlePermission(articleID int64, userID int64) (*model.Articles, error) {
	article, err := l.svcCtx.ArticleRepository.GetArticleByID(articleID)
	if err != nil {
		l.Logger.Errorf("获取文章失败: %v", err)
		return nil, err
	}

	if article == nil {
		return nil, errors.New(conf.ErrArticleNotFound)
	}

	if userID != article.AuthorID {
		return nil, errors.New(conf.ErrNoPermission)
	}

	return article, nil
}

// updateArticleFields 更新文章字段
func (l *UpdateArticleLogic) updateArticleFields(article *model.Articles, in *coderhub.UpdateArticleRequest) {
	article.Title = in.Title
	article.Content = in.Content
	article.Summary = in.Summary
	article.Tags = strings.Join(in.Tags, ",")
	article.Status = in.Status
	if in.ImageIds != nil {
		// 更新到的时候有新的图片，则删除旧的图片
		batchGetImageService := imagerelationservicelogic.NewDeleteByEntityIDLogic(l.ctx, l.svcCtx)
		_, err := batchGetImageService.DeleteByEntityID(&coderhub.DeleteByEntityIDRequest{
			EntityId:   article.ID,
			EntityType: model.ImageRelationArticleContent,
		})
		if err != nil {
			l.Logger.Errorf("删除旧的图片关联失败: %v", err)
		}
		// 批量创建新的图片关联
		relations := make([]*coderhub.CreateRelationRequest, len(in.ImageIds))
		for i, imageId := range in.ImageIds {
			relations[i] = &coderhub.CreateRelationRequest{
				ImageId:    imageId,
				EntityId:   article.ID,
				EntityType: model.ImageRelationArticleContent,
			}
		}

		batchCreateImageService := imagerelationservicelogic.NewBatchCreateRelationLogic(l.ctx, l.svcCtx)
		if _, err := batchCreateImageService.BatchCreateRelation(&coderhub.BatchCreateRelationRequest{
			Relations: relations,
		}); err != nil {
			return
		}
	}
}
