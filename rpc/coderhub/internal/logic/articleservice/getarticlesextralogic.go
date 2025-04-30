package articleservicelogic

import (
	"context"
	"fmt"

	"coderhub/model"
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
	type result struct {
		author         int64
		likeCount      int64
		articlePV      *model.ArticlePV
		commentCount   int64
		favoriteCount  int64
		isUserLiked    map[int64]bool
		isUserFavorite map[int64]bool
		err            error
	}

	resChan := make(chan result, 1)

	go func() {
		var (
			author         int64
			likeCount      int64
			articlePV      *model.ArticlePV
			commentCount   int64
			favoriteCount  int64
			isUserLiked    map[int64]bool
			isUserFavorite map[int64]bool
			errs           []error
		)

		// 使用 goroutine 并行获取各项信息

		authorCh := make(chan struct {
			author int64
			err    error
		})

		go func() {
			article, err := l.svcCtx.ArticleRepository.GetArticleByID(in.Id)
			authorCh <- struct {
				author int64
				err    error
			}{article.AuthorID, err}
		}()

		likeCountCh := make(chan struct {
			count int64
			err   error
		})

		go func() {
			count, err := l.svcCtx.ArticlesRelationLikeRepository.List(l.ctx, in.Id)
			likeCountCh <- struct {
				count int64
				err   error
			}{count, err}
		}()

		articlePVCh := make(chan struct {
			pv  *model.ArticlePV
			err error
		})
		go func() {
			pv, err := l.svcCtx.ArticlePVRepository.GetArticlePVByArticleID(in.Id)
			articlePVCh <- struct {
				pv  *model.ArticlePV
				err error
			}{pv, err}
		}()

		commentCountCh := make(chan struct {
			count int64
			err   error
		})
		go func() {
			count, err := l.svcCtx.CommentRepository.CountByArticleID(l.ctx, in.Id)
			commentCountCh <- struct {
				count int64
				err   error
			}{count, err}
		}()

		favoriteCountCh := make(chan struct {
			count int64
			err   error
		})
		go func() {
			count, err := l.svcCtx.UserFavorEntityRepository.GetEntityFavorCount(l.ctx, in.Id, "article")
			favoriteCountCh <- struct {
				count int64
				err   error
			}{count, err}
		}()

		isUserLikedCh := make(chan struct {
			liked map[int64]bool
			err   error
		})
		go func() {
			liked, err := l.svcCtx.ArticlesRelationLikeRepository.BatchArticlesHasBeenUserLiked(l.ctx, []int64{in.Id}, in.UserId)
			isUserLikedCh <- struct {
				liked map[int64]bool
				err   error
			}{liked, err}
		}()

		isUserFavoriteCh := make(chan struct {
			favorite map[int64]bool
			err      error
		})
		go func() {
			favorite, err := l.svcCtx.UserFavorEntityRepository.BatchGetUserFavorEntity(l.ctx, []int64{in.Id}, in.UserId)
			isUserFavoriteCh <- struct {
				favorite map[int64]bool
				err      error
			}{favorite, err}
		}()

		// 收集结果
		authorRes := <-authorCh
		if authorRes.err != nil {
			l.Logger.Errorf("获取文章作者失败: %v", authorRes.err)
			errs = append(errs, authorRes.err)
		} else {
			author = authorRes.author
		}

		likeRes := <-likeCountCh
		if likeRes.err != nil {
			l.Logger.Errorf("获取文章点赞数失败: %v", likeRes.err)
			errs = append(errs, likeRes.err)
		} else {
			likeCount = likeRes.count
		}

		articlePVRes := <-articlePVCh
		if articlePVRes.err != nil {
			l.Logger.Errorf("获取文章浏览量失败: %v", articlePVRes.err)
			errs = append(errs, articlePVRes.err)
		} else {
			articlePV = articlePVRes.pv
		}

		commentRes := <-commentCountCh
		if commentRes.err != nil {
			l.Logger.Errorf("获取文章评论数失败: %v", commentRes.err)
			errs = append(errs, commentRes.err)
		} else {
			commentCount = commentRes.count
		}

		favoriteRes := <-favoriteCountCh
		if favoriteRes.err != nil {
			l.Logger.Errorf("获取文章收藏数量失败: %v", favoriteRes.err)
			errs = append(errs, favoriteRes.err)
		} else {
			favoriteCount = favoriteRes.count
		}

		isUserLikedRes := <-isUserLikedCh
		if isUserLikedRes.err != nil {
			l.Logger.Errorf("获取文章是否被用户点赞失败: %v", isUserLikedRes.err)
			errs = append(errs, isUserLikedRes.err)
		} else {
			isUserLiked = isUserLikedRes.liked
		}

		isUserFavoriteRes := <-isUserFavoriteCh
		if isUserFavoriteRes.err != nil {
			l.Logger.Errorf("获取文章是否被用户收藏失败: %v", isUserFavoriteRes.err)
			errs = append(errs, isUserFavoriteRes.err)
		} else {
			isUserFavorite = isUserFavoriteRes.favorite
		}

		if len(errs) > 0 {
			resChan <- result{err: fmt.Errorf("获取文章额外信息失败: %v", errs)}
			return
		}

		resChan <- result{
			author:         author,
			likeCount:      likeCount,
			articlePV:      articlePV,
			commentCount:   commentCount,
			favoriteCount:  favoriteCount,
			isUserLiked:    isUserLiked,
			isUserFavorite: isUserFavorite,
		}
	}()

	res := <-resChan
	if res.err != nil {
		return nil, res.err
	}

	return &coderhub.ArticleAdditionalInfo{
		Id:            in.Id,
		ViewCount:     res.articlePV.Count,
		LikeCount:     res.likeCount,
		CommentCount:  res.commentCount,
		IsLicked:      res.isUserLiked[in.Id],
		IsFavorite:    res.isUserFavorite[in.Id],
		FavoriteCount: int32(res.favoriteCount),
		AuthorId:      res.author,
	}, nil
}
