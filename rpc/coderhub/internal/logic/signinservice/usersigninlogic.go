package signinservicelogic

import (
	"coderhub/conf"
	"context"
	"time"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserSignInLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserSignInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserSignInLogic {
	return &UserSignInLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserSignInLogic) UserSignIn(in *coderhub.UserSignInRequest) (*coderhub.SignInResponse, error) {
	err := l.svcCtx.UserSigninLogRepository.AddSigninLog(l.ctx, in.UserId, time.Now())
	if err != nil {
		return nil, err
	}
	pools, err := l.svcCtx.UserSigninLogRepository.GetUserSigninStatusInMonth(l.ctx, in.UserId, time.Now().Year(), time.Now().Month())
	if err != nil {
		return nil, err
	}
	// 每签到一次，用户等级加 conf.SignInLevel
	_ = l.svcCtx.UserRepository.IncrUserLevel(in.UserId, conf.SignInLevel)
	return &coderhub.SignInResponse{
		Success: true,
		Days:    pools,
	}, nil
}
