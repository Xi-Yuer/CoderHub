package articleservicelogic

import (
	"coderhub/model"
	imagerelationservicelogic "coderhub/rpc/coderhub/internal/logic/imagerelationservice"
	"context"
	"fmt"
	"strings"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListArticlesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListArticlesLogic {
	return &ListArticlesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListArticlesLogic) ListArticles(in *coderhub.GetArticlesRequest) (*coderhub.GetArticlesResponse, error) {
	// 参数校验
	if len(in.Ids) == 0 {
		return nil, fmt.Errorf("文章 ID 列表不能为空")
	}

	// 获取文章列表
	articles, err := l.svcCtx.ArticleRepository.GetArticlesByIDs(in.Ids)
	if err != nil {
		l.Logger.Errorf("批量获取文章失败: %v", err)
		return nil, fmt.Errorf("获取文章失败: %v", err)
	}

	if len(articles) == 0 {
		return nil, fmt.Errorf("文章不存在")
	}

	// 获取文章配图和封面图
	batchGetImageService := imagerelationservicelogic.NewBatchGetImagesByEntityLogic(l.ctx, l.svcCtx)
	contentImages, err := batchGetImageService.BatchGetImagesByEntity(&coderhub.BatchGetImagesByEntityRequest{
		EntityIds:  in.Ids,
		EntityType: model.ImageRelationArticleContent,
	})
	if err != nil {
		l.Logger.Errorf("获取文章配图失败: %v", err)
		return nil, fmt.Errorf("获取文章配图失败: %v", err)
	}
	coverImages, err := batchGetImageService.BatchGetImagesByEntity(&coderhub.BatchGetImagesByEntityRequest{
		EntityIds:  in.Ids,
		EntityType: model.ImageRelationArticleCover,
	})
	if err != nil {
		l.Logger.Errorf("获取文章封面失败: %v", err)
		return nil, fmt.Errorf("获取文章封面失败: %v", err)
	}

	// 获取点赞数、浏览量和评论数
	likeCounts, err := l.svcCtx.ArticlesRelationLikeRepository.BatchList(l.ctx, in.Ids)
	if err != nil {
		l.Logger.Errorf("批量获取文章点赞数失败: %v", err)
		return nil, fmt.Errorf("获取点赞数失败: %v", err)
	}

	articlePVs, err := l.svcCtx.ArticlePVRepository.GetArticlePVsByArticleIDs(in.Ids)
	if err != nil {
		l.Logger.Errorf("批量获取文章浏览量失败: %v", err)
		return nil, fmt.Errorf("获取浏览量失败: %v", err)
	}
	articlePVsMap := make(map[int64]int64, len(articlePVs))
	for _, articlePV := range articlePVs {
		articlePVsMap[articlePV.ArticleID] = articlePV.Count
	}

	commentCounts, err := l.svcCtx.CommentRepository.BatchCountByArticleIDs(l.ctx, in.Ids)
	if err != nil {
		l.Logger.Errorf("批量获取文章评论数失败: %v", err)
		return nil, fmt.Errorf("获取评论数失败: %v", err)
	}

	// 获取作者信息
	authorIDs := make([]int64, 0)
	for _, article := range articles {
		authorIDs = append(authorIDs, article.AuthorID)
	}
	authors, err := l.svcCtx.UserRepository.BatchGetUserByID(authorIDs)
	if err != nil {
		l.Logger.Errorf("获取作者信息失败: %v", err)
		return nil, fmt.Errorf("获取作者信息失败: %v", err)
	}

	// 获取文章是否被用户点赞
	isUserLiked, err := l.svcCtx.ArticlesRelationLikeRepository.BatchArticlesHasBeenUserLiked(l.ctx, in.Ids, in.UserId)

	// 获取文章是否被用户收藏
	isUserFavorite, err := l.svcCtx.UserFavorEntityRepository.BatchGetUserFavorEntity(l.ctx, in.Ids, in.UserId)

	// 构造响应
	response := make([]*coderhub.GetArticleResponse, 0)
	authorMap := make(map[int64]*coderhub.UserInfo)
	for _, author := range authors {
		authorMap[author.ID] = &coderhub.UserInfo{
			UserId:    author.ID,
			UserName:  author.UserName,
			Avatar:    author.Avatar.String,
			Email:     author.Email.String,
			Gender:    author.Gender,
			Age:       author.Age,
			Phone:     author.Phone.String,
			NickName:  author.NickName.String,
			IsAdmin:   author.IsAdmin,
			Status:    author.Status,
			CreatedAt: author.CreatedAt.Unix(),
			UpdatedAt: author.UpdatedAt.Unix(),
		}
	}

	for _, article := range articles {
		// 构造配图和封面图
		var images []*coderhub.Image
		for _, image := range contentImages.Relations {
			if image.EntityId == article.ID {
				images = append(images, &coderhub.Image{
					ImageId:      image.ImageId,
					Url:          image.Url,
					ThumbnailUrl: image.ThumbnailUrl,
				})
			}
		}

		var coverImage *coderhub.Image
		for _, image := range coverImages.Relations {
			if image.EntityId == article.ID {
				coverImage = &coderhub.Image{
					ImageId:      image.ImageId,
					Url:          image.Url,
					ThumbnailUrl: image.ThumbnailUrl,
				}
				break
			}
		}

		// 获取点赞、浏览和评论数据
		viewCount := articlePVsMap[article.ID]
		likeCount := likeCounts[article.ID]
		commentCount := commentCounts[article.ID]

		// 构造文章响应
		var tags []string
		if article.Tags != "" {
			tags = strings.Split(article.Tags, ",")
		}
		response = append(response, &coderhub.GetArticleResponse{
			Article: &coderhub.Article{
				Id:    article.ID,
				Type:  article.Type,
				Title: article.Title,
				Content: func() string {
					if article.Type == "article" {
						return article.Summary
					} else {
						return article.Content
					}
				}(),
				Summary:      article.Summary,
				Images:       images,
				CoverImage:   coverImage,
				AuthorId:     article.AuthorID,
				Tags:         tags,
				ViewCount:    viewCount,
				LikeCount:    likeCount,
				CommentCount: commentCount,
				Status:       article.Status,
				CreatedAt:    article.CreatedAt.Unix(),
				UpdatedAt:    article.UpdatedAt.Unix(),
				IsLicked:     isUserLiked[article.ID],
				IsFavorite:   isUserFavorite[article.ID],
			},
			Author: authorMap[article.AuthorID],
		})
	}
	fmt.Println("获取到的文章结果: ", response)
	return &coderhub.GetArticlesResponse{
		Articles: response,
	}, nil
}
