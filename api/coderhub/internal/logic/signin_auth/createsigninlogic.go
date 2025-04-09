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

type CreateSignInLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateSignInLogic 创建一个签到
func NewCreateSignInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSignInLogic {
	return &CreateSignInLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateSignInLogic) CreateSignIn(req *types.CreateSignReq) (resp *types.CreateSignResp, err error) {
	userID, err := utils.GetUserID(l.ctx)
	if err != nil {
		return l.errorResp(err)
	}
	signIn, err := l.svcCtx.SignInService.UserSignIn(l.ctx, &coderhub.UserSignInRequest{
		UserId: userID,
	})
	if err != nil {
		return l.errorResp(err)
	}
	var signInResp []types.UserSign
	for _, v := range signIn.Days {
		signInResp = append(signInResp, types.UserSign{
			IsSignin: v,
		})
	}
	return l.successResp(signInResp)
}

func (l *CreateSignInLogic) errorResp(err error) (resp *types.CreateSignResp, err1 error) {
	return &types.CreateSignResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: nil,
	}, nil
}

func (l *CreateSignInLogic) successResp(data []types.UserSign) (resp *types.CreateSignResp, err error) {
	return &types.CreateSignResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: data,
	}, nil
}
