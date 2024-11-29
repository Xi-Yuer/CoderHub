package logic

import (
	"coderhub/model"
	"coderhub/rpc/ImageRelation/imagerelationservice"
	"coderhub/rpc/TechSphere/Articles/articles"
	"coderhub/rpc/TechSphere/Articles/internal/svc"
	"coderhub/shared/SnowFlake"
	"coderhub/shared/Validator"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateArticleLogic {
	return &CreateArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// validateCreateArticleRequest 验证文章创建请求
func (l *CreateArticleLogic) validateCreateArticleRequest(req *articles.CreateArticleRequest) error {
	// 验证基本字段
	if err := Validator.New().
		Title(req.Title).
		Summary(req.Summary).
		Content(req.Content).
		Tags(req.Tags).
		Check(); err != nil {
		return fmt.Errorf("字段验证失败: %w", err)
	}

	// 验证图片数量
	if len(req.ImageIds) > maxImageCount {
		return fmt.Errorf("图片数量不能超过%d张", maxImageCount)
	}

	return nil
}

// CreateArticle 创建文章
func (l *CreateArticleLogic) CreateArticle(in *articles.CreateArticleRequest) (*articles.CreateArticleResponse, error) {
	// 验证请求参数
	if err := l.validateCreateArticleRequest(in); err != nil {
		return nil, fmt.Errorf("请求参数验证失败: %w", err)
	}

	// 生成文章ID
	articleID := SnowFlake.GenID()
	// 获取封面图片
	coverImageId, err := strconv.ParseInt(in.CoverImageId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("封面图片ID转换失败: %w", err)
	}

	// 创建文章模型
	article := &model.Articles{
		ID:       articleID,
		Type:     in.Type,
		Title:    in.Title,
		Content:  in.Content,
		Summary:  in.Summary,
		AuthorID: in.AuthorId,
		Tags:     strings.Join(in.Tags, ","),
		Status:   in.Status,
	}

	// 获取封面图片并创建封面图片关联
	coverImageRelation := &model.ImageRelation{
		ImageID:    coverImageId,
		EntityID:   articleID,
		EntityType: model.ImageRelationArticleCover,
		Sort:       0,
	}

	// 处理正文配图关联
	imageRelations := make([]*model.ImageRelation, len(in.ImageIds))
	for i, imageId := range in.ImageIds {
		imageIdInt, err := strconv.ParseInt(imageId, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("图片ID转换失败: %w", err)
		}
		imageRelations[i] = &model.ImageRelation{
			ImageID:    imageIdInt,
			EntityID:   articleID,
			EntityType: model.ImageRelationArticleContent,
			Sort:       int32(i),
		}
	}

	// 合并封面和配图关联
	allImageRelations := append([]*model.ImageRelation{coverImageRelation}, imageRelations...)

	// 转换为请求格式
	imageRelationReq := make([]*imagerelationservice.CreateRelationRequest, len(allImageRelations))
	for i, imageRelation := range allImageRelations {
		imageRelationReq[i] = &imagerelationservice.CreateRelationRequest{
			ImageId:    imageRelation.ImageID,
			EntityId:   imageRelation.EntityID,
			EntityType: imageRelation.EntityType,
			Sort:       imageRelation.Sort,
		}
	}

	// 保存图片关联
	if _, err := l.svcCtx.ImageRelationService.BatchCreateRelation(l.ctx, &imagerelationservice.BatchCreateRelationRequest{
		Relations: imageRelationReq,
	}); err != nil {
		return nil, fmt.Errorf("保存图片关联失败: %w", err)
	}

	// 保存文章
	if err := l.svcCtx.ArticleRepository.CreateArticle(article); err != nil {
		// 事务回滚
		_, _ = l.svcCtx.ImageRelationService.BatchDeleteRelation(l.ctx, &imagerelationservice.BatchDeleteRelationRequest{
			Ids: []int64{articleID},
		})
		return nil, fmt.Errorf("保存文章失败: %w", err)
	}

	return &articles.CreateArticleResponse{
		Id: articleID,
	}, nil
}
