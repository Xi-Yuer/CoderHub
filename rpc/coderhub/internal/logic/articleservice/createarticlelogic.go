package articleservicelogic

import (
	"coderhub/conf"
	"coderhub/model"
	imagerelationservicelogic "coderhub/rpc/coderhub/internal/logic/imagerelationservice"
	"coderhub/shared/utils"
	"context"
	"fmt"
	"strings"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

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

func (l *CreateArticleLogic) CreateArticle(in *coderhub.CreateArticleRequest) (*coderhub.CreateArticleResponse, error) {
	// 验证请求参数
	if err := l.validateCreateArticleRequest(in); err != nil {
		return nil, fmt.Errorf("请求参数验证失败: %w", err)
	}

	// 生成文章ID
	articleID := utils.GenID()

	// 文章 Tags，因为Tags 可能是空的，所以需要判断
	tags := ""
	if len(in.Tags) > 0 {
		tags = strings.Join(in.Tags, ",")
	}

	// 创建文章模型
	article := &model.Articles{
		ID:       articleID,
		Type:     in.Type,
		Title:    in.Title,
		Content:  in.Content,
		Summary:  in.Summary,
		AuthorID: in.AuthorId,
		Tags:     tags,
		Status:   in.Status,
	}

	// 创建封面图片关联
	coverImageRelation := &model.ImageRelation{
		ImageID:    in.CoverImageId,
		EntityID:   articleID,
		EntityType: model.ImageRelationArticleCover,
		Sort:       0,
	}

	// 创建正文配图关联
	imageRelations := make([]*model.ImageRelation, len(in.ImageIds))
	for i, imageId := range in.ImageIds {
		imageRelations[i] = &model.ImageRelation{
			ImageID:    imageId,
			EntityID:   articleID,
			EntityType: model.ImageRelationArticleContent,
			Sort:       int32(i),
		}
	}

	// 合并封面和配图关联
	allImageRelations := append([]*model.ImageRelation{coverImageRelation}, imageRelations...)

	// 转换为请求格式
	imageRelationReq := make([]*coderhub.CreateRelationRequest, len(allImageRelations))
	for i, imageRelation := range allImageRelations {
		imageRelationReq[i] = &coderhub.CreateRelationRequest{
			ImageId:    imageRelation.ImageID,
			EntityId:   imageRelation.EntityID,
			EntityType: imageRelation.EntityType,
			Sort:       imageRelation.Sort,
		}
	}

	// 保存图片关联
	batchCreateRelationship := imagerelationservicelogic.NewBatchCreateRelationLogic(l.ctx, l.svcCtx)
	batchDeleteRelationShip := imagerelationservicelogic.NewBatchDeleteRelationLogic(l.ctx, l.svcCtx)
	if _, err := batchCreateRelationship.BatchCreateRelation(&coderhub.BatchCreateRelationRequest{
		Relations: imageRelationReq,
	}); err != nil {
		return nil, fmt.Errorf("保存图片关联失败: %w", err)
	}
	l.Logger.Info("RPC: 保存图片关联成功, 图片关联数量:", len(imageRelationReq))

	// 保存文章
	if err := l.svcCtx.ArticleRepository.CreateArticle(article); err != nil {
		// 事务回滚
		_, _ = batchDeleteRelationShip.BatchDeleteRelation(&coderhub.BatchDeleteRelationRequest{
			Ids: []int64{articleID},
		})
		return nil, fmt.Errorf("保存文章失败: %w", err)
	}

	return &coderhub.CreateArticleResponse{
		Id: articleID,
	}, nil
}

// validateCreateArticleRequest 验证文章创建请求
func (l *CreateArticleLogic) validateCreateArticleRequest(req *coderhub.CreateArticleRequest) error {
	// 验证基本字段
	if err := utils.NewValidator().
		Title(req.Title).
		Summary(req.Summary).
		Content(req.Content).
		Tags(req.Tags).
		Check(); err != nil {
		return fmt.Errorf("字段验证失败: %w", err)
	}

	// 验证图片数量
	if len(req.ImageIds) > conf.MaxImageCount {
		return fmt.Errorf("图片数量不能超过%d张", conf.MaxImageCount)
	}

	return nil
}
