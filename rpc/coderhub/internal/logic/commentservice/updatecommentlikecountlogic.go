package commentservicelogic

import (
	"coderhub/conf"
	"coderhub/model"
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCommentLikeCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateCommentLikeCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCommentLikeCountLogic {
	return &UpdateCommentLikeCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateCommentLikeCount 更新评论点赞数
func (l *UpdateCommentLikeCountLogic) UpdateCommentLikeCount(in *coderhub.UpdateCommentLikeCountRequest) (*coderhub.UpdateCommentLikeCountResponse, error) {
	// 更新文章点赞数
	commentRelationLike := model.CommentRelationLike{
		CommentID: in.CommentId,
		UserID:    in.UserId,
	}
	// 判断是否点赞
	isLike := l.svcCtx.CommentRelationLikeRepository.Get(l.ctx, &commentRelationLike)
	if isLike {
		// 取消点赞
		err := l.svcCtx.CommentRelationLikeRepository.Delete(l.ctx, &commentRelationLike)
		if err != nil {
			return nil, err
		}
		go func() {
			_ = l.svcCtx.UserRepository.IncrUserLevel(in.UserId, conf.DisCommentLikeLevel)
		}()
	} else {
		// 点赞
		err := l.svcCtx.CommentRelationLikeRepository.Create(l.ctx, &commentRelationLike)
		if err != nil {
			return nil, err
		}
		// 增加用户经验
		go func() {
			_ = l.svcCtx.UserRepository.IncrUserLevel(in.UserId, conf.CommentLikeLevel)
		}()
	}
	return &coderhub.UpdateCommentLikeCountResponse{
		Success: true,
	}, nil
}
