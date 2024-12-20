package userservicelogic

import (
	"coderhub/model"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"
	"coderhub/shared/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
)

type UpdateUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateUserInfo 更新用户信息
func (l *UpdateUserInfoLogic) UpdateUserInfo(in *coderhub.UpdateUserInfoRequest) (*coderhub.UpdateUserInfoResponse, error) {
	if userId, err := utils.GetUserMetaData(l.ctx); err != nil || userId != strconv.FormatInt(in.UserId, 10) {
		return nil, errors.New("非法操作")
	}
	UserInfo, err := NewGetUserInfoLogic(l.ctx, l.svcCtx).GetUserInfo(&coderhub.GetUserInfoRequest{
		UserId: in.UserId,
	})
	if err != nil {
		return nil, err
	}
	if UserInfo.UserId == 0 {
		return nil, errors.New("用户不存在")
	}

	fmt.Printf("in:%#v\n", in)
	oldUser, err := l.svcCtx.UserRepository.GetUserByID(UserInfo.UserId)
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.UserRepository.UpdateUser(&model.User{
		ID:       UserInfo.UserId,
		UserName: oldUser.UserName,
		Password: oldUser.Password,
		NickName: sql.NullString{String: in.Nickname, Valid: in.Nickname != ""},
		Phone:    sql.NullString{String: in.Phone, Valid: in.Phone != ""},
		Age:      in.Age,
		Gender:   in.Gender,
		Email:    sql.NullString{String: in.Email, Valid: in.Email != ""},
		Status:   oldUser.Status,
		IsAdmin:  oldUser.IsAdmin,
	}); err != nil {
		return nil, err
	}

	return &coderhub.UpdateUserInfoResponse{
		Success: true,
	}, nil
}
