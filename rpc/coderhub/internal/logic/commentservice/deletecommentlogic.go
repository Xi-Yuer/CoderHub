package commentservicelogic

import (
	imagerelationservicelogic "coderhub/rpc/coderhub/internal/logic/imagerelationservice"
	"coderhub/shared/utils"
	"context"
	"fmt"
	"strconv"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

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
func (l *DeleteCommentLogic) DeleteComment(in *coderhub.DeleteCommentRequest) (*coderhub.DeleteCommentResponse, error) {
	// 权限校验
	var (
		userId string
		err    error
	)
	if userId, err = utils.GetUserMetaData(l.ctx); err != nil {
		return nil, err
	}

	if userId != strconv.FormatInt(in.UserId, 10) {
		return nil, fmt.Errorf("非法操作")
	}

	if err := l.svcCtx.CommentRepository.Delete(l.ctx, in.CommentId); err != nil {
		return nil, err
	}
	// 删除图片关联
	batchDeleteRelationService := imagerelationservicelogic.NewBatchDeleteRelationLogic(l.ctx, l.svcCtx)
	_, err = batchDeleteRelationService.BatchDeleteRelation(&coderhub.BatchDeleteRelationRequest{
		Ids: []int64{in.CommentId},
	})
	if err != nil {
		return nil, err
	}
	return &coderhub.DeleteCommentResponse{
		Success: true,
	}, nil
}
