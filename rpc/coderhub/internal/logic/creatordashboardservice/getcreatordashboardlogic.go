package creatordashboardservicelogic

import (
	"coderhub/model"
	"context"
	"sync"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCreatorDashBoardLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCreatorDashBoardLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCreatorDashBoardLogic {
	return &GetCreatorDashBoardLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCreatorDashBoardLogic) GetCreatorDashBoard(in *coderhub.GetCreatorDashBoardRequest) (*coderhub.GetCreatorDashBoardResponse, error) {
	var wg sync.WaitGroup
	var articleCount int64
	var microPostCount int64
	var likeCount int64
	var commentCount int64
	var articleFavorCount int64
	var followerCount int64
	var birthday int64
	var articlePV int64
	var err error
	wg.Add(6)

	// 获取用户文章数量
	go func() {
		defer wg.Done()
		articleCount, err = l.svcCtx.ArticleRepository.GetUserArticleCount(in.UserId)
	}()
	// 获取用户点赞数量
	go func() {
		defer wg.Done()
	}()
	// 获取用户关注数量
	go func() {
		defer wg.Done()
		followerCount, err = l.svcCtx.UserFollowRepository.GetUserFansCount(in.UserId)
	}()
	// 获取用户沸点数量
	go func() {
		defer wg.Done()
		microPostCount, err = l.svcCtx.ArticleRepository.GetUserMicroPostCount(in.UserId)
	}()
	// 获取用户生日
	go func() {
		defer wg.Done()
		birthday, err = l.svcCtx.UserRepository.GetUserBirthday(in.UserId)
	}()
	go func() {
		defer wg.Done()
		// 获取用户所有的文章ID
		ids, _ := l.svcCtx.ArticleRepository.GetUserAllArticleIDS(in.UserId)
		commentCount, err = l.svcCtx.CommentRepository.GetEntityCommentCount(l.ctx, ids)
		articleFavorCount, err = l.svcCtx.UserFavorEntityRepository.GetEntitiesFavorCount(l.ctx, ids, model.ArticleType)
		articlePV, err = l.svcCtx.ArticlePVRepository.GetAllArticlePVByArticleIDs(ids)
		likeCount, err = l.svcCtx.ArticlesRelationLikeRepository.GetArticleLikeCount(l.ctx, ids)
	}()

	wg.Wait()

	if err != nil {
		return nil, err
	}

	return &coderhub.GetCreatorDashBoardResponse{
		ArticleCount:         articleCount,
		MicroPostCount:       microPostCount,
		LikeCount:            likeCount,
		CommentCount:         commentCount,
		FollowerCount:        followerCount,
		ArticleViewCount:     articlePV,
		ArticleFavoriteCount: articleFavorCount,
		CreatorBirthday:      birthday,
	}, nil
}
