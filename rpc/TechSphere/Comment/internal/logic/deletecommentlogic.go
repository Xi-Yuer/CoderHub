package logic

import (
	"context"

	"coderhub/rpc/TechSphere/Comment/comment"
	"coderhub/rpc/TechSphere/Comment/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCommentLogic {
	return &DeleteCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeleteComment 删除评论
func (l *DeleteCommentLogic) DeleteComment(in *comment.DeleteCommentRequest) (*comment.DeleteCommentResponse, error) {
	if err := l.svcCtx.CommentRepository.Delete(l.ctx, in.CommentId); err != nil {
		return nil, err
	}
	return &comment.DeleteCommentResponse{
		Success: true,
	}, nil
}