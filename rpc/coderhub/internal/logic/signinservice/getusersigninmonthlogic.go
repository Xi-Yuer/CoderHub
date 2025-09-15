package signinservicelogic

import (
	"context"
	"fmt"
	"time"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserSignInMonthLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserSignInMonthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserSignInMonthLogic {
	return &GetUserSignInMonthLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserSignInMonthLogic) GetUserSignInMonth(in *coderhub.GetUserSignInMonthRequest) (*coderhub.SignInResponse, error) {
	pools, err := l.svcCtx.UserSigninLogRepository.GetUserSigninStatusInMonth(l.ctx, in.UserId, int(in.Year), time.Month(in.Month))
	if err != nil {
		return nil, err
	}

	fmt.Println("pools", pools)

	return &coderhub.SignInResponse{
		Success: true,
		Days:    pools,
	}, nil
}
