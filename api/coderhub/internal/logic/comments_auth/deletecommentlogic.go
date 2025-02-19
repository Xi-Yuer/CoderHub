package comments_auth

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/client/commentservice"
	"coderhub/shared/utils"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewDeleteCommentLogic 删除评论
func NewDeleteCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCommentLogic {
	return &DeleteCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteCommentLogic) DeleteComment(req *types.DeleteCommentReq) (resp *types.DeleteCommentResp, err error) {
	// 权限校验
	userID, err := utils.GetUserID(l.ctx)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return l.errorResp(err)
	}

	_, err = l.svcCtx.CommentService.DeleteComment(utils.SetUserMetaData(l.ctx), &commentservice.DeleteCommentRequest{
		CommentId: utils.String2Int(req.CommentId),
		UserId:    userID,
	})
	if err != nil {
		return l.errorResp(err)
	}

	return l.successResp()
}

func (l *DeleteCommentLogic) successResp() (*types.DeleteCommentResp, error) {
	return &types.DeleteCommentResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: true,
	}, nil
}

func (l *DeleteCommentLogic) errorResp(err error) (*types.DeleteCommentResp, error) {
	return &types.DeleteCommentResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
	}, nil
}
