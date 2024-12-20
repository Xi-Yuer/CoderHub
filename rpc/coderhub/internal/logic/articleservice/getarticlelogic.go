package articleservicelogic

import (
	"coderhub/model"
	imagerelationservicelogic "coderhub/rpc/coderhub/internal/logic/imagerelationservice"
	"coderhub/shared/utils"
	"context"
	"errors"
	"fmt"
	"strings"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticleLogic {
	return &GetArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetArticleLogic) GetArticle(in *coderhub.GetArticleRequest) (*coderhub.GetArticleResponse, error) {
	// 参数校验
	if err := utils.NewValidator().ArticleID(in.Id).Check(); err != nil {
		return nil, errors.New(err.Error())
	}

	// 从数据库获取文章
	article, err := l.svcCtx.ArticleRepository.GetArticleByID(in.Id)
	if err != nil {
		l.Logger.Errorf("获取文章失败: %v", err)
		return nil, fmt.Errorf("获取文章失败: %v", err)
	}

	if article == nil {
		return nil, fmt.Errorf("文章不存在")
	}

	// 获取图片关联（文章配图）
	batchGetImageService := imagerelationservicelogic.NewBatchGetImagesByEntityLogic(l.ctx, l.svcCtx)
	images, err := batchGetImageService.BatchGetImagesByEntity(&coderhub.BatchGetImagesByEntityRequest{
		EntityIds:  []int64{article.ID},
		EntityType: model.ImageRelationArticleContent,
	})
	if err != nil {
		l.Logger.Errorf("获取图片关联失败: %v", err)
		return nil, fmt.Errorf("获取图片关联失败: %v", err)
	}
	l.Logger.Info("RPC: 获取文章配图成功, 配图数量:", len(images.Relations))
	// 获取图片关联（文章封面）
	coverImages, err := batchGetImageService.BatchGetImagesByEntity(&coderhub.BatchGetImagesByEntityRequest{
		EntityIds:  []int64{article.ID},
		EntityType: model.ImageRelationArticleCover,
	})
	if err != nil {
		l.Logger.Errorf("获取图片关联失败: %v", err)
		return nil, fmt.Errorf("获取图片关联失败: %v", err)
	}
	l.Logger.Info("RPC: 获取文章封面成功, 封面数量:", len(coverImages.Relations))
	// 获取文章配图
	articleImages := make([]*coderhub.Image, 0)
	if len(images.Relations) > 0 {
		for _, image := range images.Relations {
			articleImages = append(articleImages, &coderhub.Image{
				ImageId:      image.ImageId,
				Url:          image.Url,
				ThumbnailUrl: image.ThumbnailUrl,
			})
		}
	}
	var coverImage *coderhub.Image
	if len(coverImages.Relations) > 0 {
		coverImage = &coderhub.Image{
			ImageId:      coverImages.Relations[0].ImageId,
			Url:          coverImages.Relations[0].Url,
			ThumbnailUrl: coverImages.Relations[0].ThumbnailUrl,
		}
	}

	// 获取文章点赞数
	likeCount, err := l.svcCtx.ArticlesRelationLikeRepository.List(l.ctx, article.ID)
	if err != nil {
		l.Logger.Errorf("获取文章点赞数失败: %v", err)
		return nil, fmt.Errorf("获取文章点赞数失败: %v", err)
	}

	// 获取文章浏览量
	articlePV, err := l.svcCtx.ArticlePVRepository.GetArticlePVByArticleID(article.ID)
	if err != nil {
		l.Logger.Errorf("获取文章浏览量失败: %v", err)
		return nil, fmt.Errorf("获取文章浏览量失败: %v", err)
	}

	// 获取文章评论数
	commentCount, err := l.svcCtx.CommentRepository.CountByArticleID(l.ctx, article.ID)
	if err != nil {
		l.Logger.Errorf("获取文章评论数失败: %v", err)
		return nil, fmt.Errorf("获取文章评论数失败: %v", err)
	}

	// 文章浏览量+1
	err = l.svcCtx.ArticlePVRepository.CreateArticlePV(&model.ArticlePV{
		ArticleID: article.ID,
		Count:     1,
	})
	if err != nil {
		l.Logger.Errorf("文章浏览量+1失败: %v", err)
	}

	// 转换为响应格式
	tags := make([]string, 0)
	if article.Tags != "" {
		tags = strings.Split(article.Tags, ",")
	} else {
		// 如果Tags为空，则返回空数组
		tags = []string{}
	}
	response := &coderhub.GetArticleResponse{
		Article: &coderhub.Article{
			Id:           article.ID,
			Type:         article.Type,
			Title:        article.Title,
			Content:      article.Content,
			Summary:      article.Summary,
			Images:       articleImages,
			CoverImage:   coverImage,
			AuthorId:     article.AuthorID,
			Tags:         tags,
			ViewCount:    articlePV.Count,
			LikeCount:    likeCount,
			CommentCount: commentCount,
			Status:       article.Status,
			CreatedAt:    article.CreatedAt.Unix(),
			UpdatedAt:    article.UpdatedAt.Unix(),
		},
	}
	return response, nil
}
