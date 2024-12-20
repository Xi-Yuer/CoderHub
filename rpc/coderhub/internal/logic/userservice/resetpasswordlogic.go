package userservicelogic

import (
	"coderhub/shared/utils"
	"context"
	"fmt"
	"strconv"
	"time"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ResetPassword 重置密码, 通过邮箱发送重置密码链接
func (l *ResetPasswordLogic) ResetPassword(in *coderhub.ResetPasswordRequest) (*coderhub.ResetPasswordResponse, error) {
	// 1. 检查用户是否存在
	userInfo, err := l.svcCtx.UserRepository.FindOneByEmail(in.Email)
	if err != nil {
		return nil, err
	}
	if userInfo == nil {
		return nil, fmt.Errorf("用户不存在")
	}
	if userInfo.Email.String == "" {
		return nil, fmt.Errorf("用户邮箱不存在")
	}
	// 2. 生成重置密码链接
	token := utils.GenID()
	// 3. 将token存入redis, 过期时间为10分钟
	err = l.svcCtx.RedisDB.SetWithTTL(fmt.Sprintf("reset_password:%s", in.Email), strconv.FormatInt(token, 10), 10*time.Minute)
	if err != nil {
		return nil, err
	}
	link := fmt.Sprintf("http://localhost/reset-password?email=%s&token=%d", userInfo.Email.String, token)
	// 3. 发送重置密码链接
	err = l.svcCtx.GoMail.SendWithHTML(userInfo.Email.String, "邮箱密码重置确认", link)
	if err != nil {
		return nil, err
	}
	// 4. 返回重置密码链接
	return &coderhub.ResetPasswordResponse{
		Success: true,
	}, nil
}
