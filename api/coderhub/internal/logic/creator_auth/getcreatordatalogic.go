package creator_auth

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCreatorDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetCreatorDataLogic 获取创作者相关数据信息
func NewGetCreatorDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCreatorDataLogic {
	return &GetCreatorDataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCreatorDataLogic) GetCreatorData(req *types.GetCreatorDashboardReq) (resp *types.GetCreatorDashboardResp, err error) {
	userID, err := utils.GetUserID(l.ctx)
	if err != nil {
		return l.errorResp(err)
	}

	response, err := l.svcCtx.CreatorDashBoardService.GetCreatorDashBoard(l.ctx, &coderhub.GetCreatorDashBoardRequest{
		UserId: userID,
	})
	if err != nil {
		return l.errorResp(err)
	}

	return l.successResp(&types.CreatorDashboard{
		ArticleCount:      response.ArticleCount,
		MicroPostCount:    response.MicroPostCount,
		LikeCount:         response.LikeCount,
		CommentCount:      response.CommentCount,
		ArticleFavorCount: response.ArticleFavoriteCount,
		FollowerCount:     response.FollowerCount,
		Birthday:          response.CreatorBirthday,
		ArticlePV:         response.ArticleViewCount,
	})
}

func (l *GetCreatorDataLogic) successResp(data *types.CreatorDashboard) (*types.GetCreatorDashboardResp, error) {
	return &types.GetCreatorDashboardResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: data,
	}, nil
}

func (l *GetCreatorDataLogic) errorResp(err error) (*types.GetCreatorDashboardResp, error) {
	return &types.GetCreatorDashboardResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
	}, nil
}
