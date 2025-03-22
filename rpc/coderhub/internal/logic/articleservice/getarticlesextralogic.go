package articleservicelogic

import (
	"context"
	"fmt"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetArticlesExtraLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetArticlesExtraLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticlesExtraLogic {
	return &GetArticlesExtraLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetArticlesExtraLogic) GetArticlesExtra(in *coderhub.GetArticleRequest) (*coderhub.ArticleAdditionalInfo, error) {
	// 获取文章点赞数
	likeCount, err := l.svcCtx.ArticlesRelationLikeRepository.List(l.ctx, in.Id)
	if err != nil {
		l.Logger.Errorf("获取文章点赞数失败: %v", err)
		return nil, fmt.Errorf("获取文章点赞数失败: %v", err)
	}

	// 获取文章浏览量
	articlePV, err := l.svcCtx.ArticlePVRepository.GetArticlePVByArticleID(in.Id)
	if err != nil {
		l.Logger.Errorf("获取文章浏览量失败: %v", err)
		return nil, fmt.Errorf("获取文章浏览量失败: %v", err)
	}

	// 获取文章评论数
	commentCount, err := l.svcCtx.CommentRepository.CountByArticleID(l.ctx, in.Id)
	if err != nil {
		l.Logger.Errorf("获取文章评论数失败: %v", err)
		return nil, fmt.Errorf("获取文章评论数失败: %v", err)
	}

	// 获取文章收藏数量
	favoriteCount, err := l.svcCtx.UserFavorEntityRepository.GetEntityFavorCount(l.ctx, in.Id, "article")

	// 获取文章是否被用户点赞
	isUserLiked, err := l.svcCtx.ArticlesRelationLikeRepository.BatchArticlesHasBeenUserLiked(l.ctx, []int64{in.Id}, in.UserId)

	// 获取文章是否被用户收藏
	isUserFavorite, err := l.svcCtx.UserFavorEntityRepository.BatchGetUserFavorEntity(l.ctx, []int64{in.Id}, in.UserId)

	return &coderhub.ArticleAdditionalInfo{
		Id:            in.Id,
		ViewCount:     articlePV.Count,
		LikeCount:     likeCount,
		CommentCount:  commentCount,
		IsLicked:      isUserLiked[in.Id],
		IsFavorite:    isUserFavorite[in.Id],
		FavoriteCount: int32(favoriteCount),
	}, nil
}
