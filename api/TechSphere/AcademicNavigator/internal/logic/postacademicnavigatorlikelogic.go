package logic

import (
	"context"

	"coderhub/api/TechSphere/AcademicNavigator/internal/svc"
	"coderhub/api/TechSphere/AcademicNavigator/internal/types"
	"coderhub/conf"
	"coderhub/rpc/TechSphere/AcademicNavigator/academic_navigator"
	"coderhub/shared/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type PostAcademicNavigatorLikeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 点赞学术导航
func NewPostAcademicNavigatorLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostAcademicNavigatorLikeLogic {
	return &PostAcademicNavigatorLikeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PostAcademicNavigatorLikeLogic) PostAcademicNavigatorLike(req *types.PostAcademicNavigatorLikeReq) (resp *types.PostAcademicNavigatorLikeResp, err error) {
	// 获取用户ID
	userId, err := utils.GetUserID(l.ctx)
	if err != nil {
		return l.errorResp(err)
	}
	_, err = l.svcCtx.AcademicNavigatorService.LikeAcademicNavigator(l.ctx, &academic_navigator.LikeAcademicNavigatorRequest{
		Id:     req.Id,
		UserId: userId,
	})
	if err != nil {
		return l.errorResp(err)
	}

	return l.successResp()
}

func (l *PostAcademicNavigatorLikeLogic) successResp() (*types.PostAcademicNavigatorLikeResp, error) {
	return &types.PostAcademicNavigatorLikeResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: true,
	}, nil
}

func (l *PostAcademicNavigatorLikeLogic) errorResp(err error) (*types.PostAcademicNavigatorLikeResp, error) {
	return &types.PostAcademicNavigatorLikeResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: false,
	}, nil
}