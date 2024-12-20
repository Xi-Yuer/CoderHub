package userservicelogic

import (
	"coderhub/model"
	"coderhub/shared/security"
	"coderhub/shared/utils"
	"context"
	"errors"
	"fmt"
	"strconv"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ChangePassword 修改密码
func (l *ChangePasswordLogic) ChangePassword(in *coderhub.ChangePasswordRequest) (*coderhub.ChangePasswordResponse, error) {
	if err := utils.NewValidator().Password(in.OldPassword).Password(in.NewPassword).Check(); err != nil {
		return nil, err
	}

	var (
		userId string
		err    error
	)
	// 从 metadata 中获取 userId
	if userId, err = utils.GetUserMetaData(l.ctx); err != nil {
		return nil, err
	}

	if userId != strconv.FormatInt(in.UserId, 10) {
		return nil, fmt.Errorf("非法操作")
	}

	userInfo, err := NewGetUserInfoLogic(l.ctx, l.svcCtx).GetUserInfo(&coderhub.GetUserInfoRequest{UserId: in.UserId})
	if err != nil {
		return nil, err
	}

	// 验证旧密码是否正确
	if !security.CompareHashAndPassword(userInfo.Password, in.OldPassword) {
		return nil, errors.New("旧密码不正确")
	}

	// 对新密码进行哈希处理
	hashedNewPassword, err := security.PasswordHash(in.NewPassword)
	if err != nil {
		return nil, err
	}

	// 更新用户密码
	if err := l.svcCtx.UserRepository.UpdateUser(&model.User{ID: in.UserId, Password: hashedNewPassword}); err != nil {
		return nil, err
	}

	// 返回成功响应
	return &coderhub.ChangePasswordResponse{
		Success: true,
	}, nil
}
