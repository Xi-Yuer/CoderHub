package comments_auth

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCommentLikeCountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUpdateCommentLikeCountLogic 更新评论点赞数
func NewUpdateCommentLikeCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCommentLikeCountLogic {
	return &UpdateCommentLikeCountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCommentLikeCountLogic) UpdateCommentLikeCount(req *types.UpdateCommentLikeCountReq) (resp *types.UpdateCommentLikeCountResp, err error) {
	userId, err := utils.GetUserID(l.ctx)
	if err != nil {
		return l.errorResp(err), nil
	}
	ctx := utils.SetUserMetaData(l.ctx)

	if _, err := l.svcCtx.CommentService.UpdateCommentLikeCount(ctx, &coderhub.UpdateCommentLikeCountRequest{
		CommentId: req.CommentId,
		UserId:    userId,
	}); err != nil {
		return l.errorResp(err), nil
	}

	return l.successResp(), nil
}

func (l *UpdateCommentLikeCountLogic) errorResp(err error) *types.UpdateCommentLikeCountResp {
	return &types.UpdateCommentLikeCountResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: false,
	}
}

func (l *UpdateCommentLikeCountLogic) successResp() *types.UpdateCommentLikeCountResp {
	return &types.UpdateCommentLikeCountResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: "success",
		},
		Data: true,
	}
}
