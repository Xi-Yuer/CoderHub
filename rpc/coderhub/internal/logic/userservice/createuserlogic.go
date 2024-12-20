package userservicelogic

import (
	"coderhub/model"
	"coderhub/shared/security"
	"coderhub/shared/utils"
	"context"
	"errors"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateUser 创建用户
func (l *CreateUserLogic) CreateUser(in *coderhub.CreateUserRequest) (*coderhub.CreateUserResponse, error) {
	if err := utils.NewValidator().Username(in.Username).Password(in.PasswordHash).Check(); err != nil {
		return nil, err
	}

	exists, _ := NewCheckUserExistsLogic(l.ctx, l.svcCtx).CheckUserExists(&coderhub.CheckUserExistsRequest{Username: in.Username})
	if exists.Exists {
		return nil, errors.New("用户已存在")
	}

	ID := utils.GenID()
	Password, _ := security.PasswordHash(in.PasswordHash)
	if err := l.svcCtx.UserRepository.CreateUser(&model.User{
		ID:       ID,
		UserName: in.Username,
		Password: Password,
	}); err != nil {
		return nil, err
	}

	return &coderhub.CreateUserResponse{
		UserId: ID,
	}, nil
}
