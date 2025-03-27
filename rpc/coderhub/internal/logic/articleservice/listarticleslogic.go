package articleservicelogic

import (
	"coderhub/model"
	"coderhub/rpc/coderhub/coderhub"
	imagerelationservicelogic "coderhub/rpc/coderhub/internal/logic/imagerelationservice"
	"coderhub/rpc/coderhub/internal/svc"
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
	"sync"
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

	// 批量获取文章
	articles, err := l.svcCtx.ArticleRepository.GetArticlesByIDs(in.Ids)
	if err != nil {
		l.Logger.Errorf("批量获取文章失败: %v", err)
		return nil, fmt.Errorf("获取文章失败: %v", err)
	}
	if len(articles) == 0 {
		return &coderhub.GetArticlesResponse{Articles: nil}, nil
	}

	// 使用 Goroutines 并行查询
	var wg sync.WaitGroup
	var contentImages, coverImages *coderhub.BatchGetImagesByEntityResponse
	var likeCounts map[int64]int64
	var articlePVs map[int64]int64
	var commentCounts map[int64]int64
	var userLiked map[int64]bool
	var userFavored map[int64]bool
	var authors []*model.User

	wg.Add(3) // 添加 3 个 Goroutine

	// 获取文章配图和封面图
	go func() {
		defer wg.Done()
		batchGetImageService := imagerelationservicelogic.NewBatchGetImagesByEntityLogic(l.ctx, l.svcCtx)
		contentImagesResp, err := batchGetImageService.BatchGetImagesByEntity(&coderhub.BatchGetImagesByEntityRequest{
			EntityIds:  in.Ids,
			EntityType: model.ImageRelationArticleContent,
		})
		if err != nil {
			l.Logger.Errorf("获取文章配图失败: %v", err)
		}
		contentImages = contentImagesResp

		coverImagesResp, err := batchGetImageService.BatchGetImagesByEntity(&coderhub.BatchGetImagesByEntityRequest{
			EntityIds:  in.Ids,
			EntityType: model.ImageRelationArticleCover,
		})
		if err != nil {
			l.Logger.Errorf("获取文章封面失败: %v", err)
		}
		coverImages = coverImagesResp
	}()

	// 获取点赞数、浏览量和评论数
	go func() {
		defer wg.Done()
		likeCounts, _ = l.svcCtx.ArticlesRelationLikeRepository.BatchList(l.ctx, in.Ids)
		articlePV, _ := l.svcCtx.ArticlePVRepository.GetArticlePVsByArticleIDs(in.Ids)
		userLiked, _ = l.svcCtx.ArticlesRelationLikeRepository.BatchArticlesHasBeenUserLiked(l.ctx, in.Ids, in.UserId)
		userFavored, _ = l.svcCtx.UserFavorEntityRepository.BatchGetUserFavorEntity(l.ctx, in.Ids, in.UserId)
		articlePVs = make(map[int64]int64)
		for _, pv := range articlePV {
			articlePVs[pv.ArticleID] = pv.Count
		}
		commentCounts, _ = l.svcCtx.CommentRepository.BatchCountByArticleIDs(l.ctx, in.Ids)
	}()

	// 获取作者信息
	go func() {
		defer wg.Done()
		authorIDs := make([]int64, len(articles))
		for i, article := range articles {
			authorIDs[i] = article.AuthorID
		}
		authorsResp, err := l.svcCtx.UserRepository.BatchGetUserByID(authorIDs)
		if err != nil {
			l.Logger.Errorf("获取作者信息失败: %v", err)
		}
		authors = authorsResp
	}()

	// 等待所有 Goroutine 执行完
	wg.Wait()

	// 处理数据
	response := make([]*coderhub.GetArticleResponse, len(articles))
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

	// 构建响应数据
	for i, article := range articles {
		// 配图和封面图处理
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

		// 获取点赞、浏览量和评论数
		viewCount := articlePVs[article.ID]
		likeCount := likeCounts[article.ID]
		commentCount := commentCounts[article.ID]
		userLiked := userLiked[article.ID]
		userFavored := userFavored[article.ID]

		// 构建文章响应
		var tags []string
		if article.Tags != "" {
			tags = strings.Split(article.Tags, ",")
		}
		response[i] = &coderhub.GetArticleResponse{
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
				CategoryId:   article.CategoryID,
				ViewCount:    viewCount,
				LikeCount:    likeCount,
				IsFavorite:   userFavored,
				IsLicked:     userLiked,
				CommentCount: commentCount,
				Status:       article.Status,
				CreatedAt:    article.CreatedAt.Unix(),
				UpdatedAt:    article.UpdatedAt.Unix(),
			},
			Author: authorMap[article.AuthorID],
		}
	}

	return &coderhub.GetArticlesResponse{Articles: response}, nil
}
