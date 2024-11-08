package logic

import (
	"coderhub/model"
	"coderhub/rpc/user/internal/svc"
	"coderhub/rpc/user/user"
	"coderhub/shared/bcryptUtil"
	"coderhub/shared/metaData"
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
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
func (l *ChangePasswordLogic) ChangePassword(in *user.ChangePasswordRequest) (*user.ChangePasswordResponse, error) {
	var (
		userId string
		err    error
	)
	// 从 metadata 中获取 userId
	if userId, err = metaData.GetMetaData(l.ctx, "userId"); err != nil {
		return nil, err
	}

	if userId != strconv.FormatInt(in.UserId, 10) {
		return nil, fmt.Errorf("非法操作")
	}

	var userInfo model.User
	l.svcCtx.SqlDB.First(&userInfo, "id = ?", in.UserId)

	// 验证旧密码是否正确
	if !bcryptUtil.CompareHashAndPassword(userInfo.Password, in.OldPassword) {
		return nil, errors.New("旧密码不正确")
	}

	// 对新密码进行哈希处理
	hashedNewPassword, err := bcryptUtil.PasswordHash(in.NewPassword)
	if err != nil {
		return nil, err
	}

	// 更新用户密码
	if tx := l.svcCtx.SqlDB.Model(&userInfo).Update("password", hashedNewPassword); tx.Error != nil {
		return nil, tx.Error
	}

	// 返回成功响应
	return &user.ChangePasswordResponse{
		Success: true,
	}, nil
}
