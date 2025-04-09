package signin_auth

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListSignInLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewListSignInLogic 获取签到列表
func NewListSignInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSignInLogic {
	return &ListSignInLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListSignInLogic) ListSignIn(req *types.GetUserSignReq) (resp *types.GetUserSignResp, err error) {
	userID, err := utils.GetUserID(l.ctx)
	if err != nil {
		return l.errorResp(err)
	}

	signInMonth, err := l.svcCtx.SignInService.GetUserSignInMonth(l.ctx, &coderhub.GetUserSignInMonthRequest{
		UserId: userID,
		Year:   int32(req.Year),
		Month:  int32(req.Month),
	})
	if err != nil {
		return l.errorResp(err)
	}

	var signInResp []types.UserSign
	for _, v := range signInMonth.Days {
		signInResp = append(signInResp, types.UserSign{
			IsSignin: v,
		})
	}
	return l.successResp(signInResp)
}

func (l *ListSignInLogic) errorResp(err error) (resp *types.GetUserSignResp, err1 error) {
	return &types.GetUserSignResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: nil,
	}, nil
}

func (l *ListSignInLogic) successResp(data []types.UserSign) (resp *types.GetUserSignResp, err error) {
	return &types.GetUserSignResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: data,
	}, nil
}
