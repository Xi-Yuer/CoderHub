package articles_auth

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLikeCountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新文章点赞数
func NewUpdateLikeCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLikeCountLogic {
	return &UpdateLikeCountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLikeCountLogic) UpdateLikeCount(req *types.UpdateLikeCountReq) (resp *types.UpdateLikeCountResp, err error) {
	userId, err := utils.GetUserID(l.ctx)
	if err != nil {
		return l.errorResp(err), nil
	}
	ctx := utils.SetUserMetaData(l.ctx)

	if err := utils.NewValidator().ArticleID(req.Id).Check(); err != nil {
		return l.errorResp(err), nil
	}

	if _, err := l.svcCtx.ArticlesService.UpdateLikeCount(ctx, &coderhub.UpdateLikeCountRequest{
		Id:     req.Id,
		UserId: userId,
	}); err != nil {
		return l.errorResp(err), nil
	}

	return l.successResp(), nil
}

func (l *UpdateLikeCountLogic) errorResp(err error) *types.UpdateLikeCountResp {
	return &types.UpdateLikeCountResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: false,
	}
}

func (l *UpdateLikeCountLogic) successResp() *types.UpdateLikeCountResp {
	return &types.UpdateLikeCountResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: "success",
		},
		Data: true,
	}
}
