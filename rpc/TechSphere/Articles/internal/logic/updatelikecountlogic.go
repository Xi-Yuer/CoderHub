package logic

import (
	"context"

	"coderhub/model"
	"coderhub/rpc/TechSphere/Articles/articles"
	"coderhub/rpc/TechSphere/Articles/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLikeCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLikeCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLikeCountLogic {
	return &UpdateLikeCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateLikeCountLogic) UpdateLikeCount(in *articles.UpdateLikeCountRequest) (*articles.UpdateLikeCountResponse, error) {
	// 更新文章点赞数
	articleRelationLike := model.ArticlesRelationLike{
		ArticleID: in.Id,
		UserID:    in.UserId,
	}
	isLike := l.svcCtx.ArticlesRelationLikeRepository.Get(l.ctx, &articleRelationLike)
	if isLike {
		// 取消点赞
		err := l.svcCtx.ArticlesRelationLikeRepository.Delete(l.ctx, &articleRelationLike)
		if err != nil {
			return nil, err
		}
	} else {
		// 点赞
		err := l.svcCtx.ArticlesRelationLikeRepository.Create(l.ctx, &articleRelationLike)
		if err != nil {
			return nil, err
		}
	}

	return &articles.UpdateLikeCountResponse{Success: true}, nil
}
