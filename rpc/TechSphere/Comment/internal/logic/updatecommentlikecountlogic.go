package logic

import (
	"context"

	"coderhub/model"
	"coderhub/rpc/TechSphere/Comment/comment"
	"coderhub/rpc/TechSphere/Comment/internal/svc"

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
func (l *UpdateCommentLikeCountLogic) UpdateCommentLikeCount(in *comment.UpdateCommentLikeCountRequest) (*comment.UpdateCommentLikeCountResponse, error) {
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
	} else {
		// 点赞
		err := l.svcCtx.CommentRelationLikeRepository.Create(l.ctx, &commentRelationLike)
		if err != nil {
			return nil, err
		}
	}
	return &comment.UpdateCommentLikeCountResponse{
		Success: true,
	}, nil
}