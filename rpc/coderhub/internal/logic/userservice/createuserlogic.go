package userservicelogic

import (
	"coderhub/model"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"
	"coderhub/shared/security"
	"coderhub/shared/utils"
	"context"
	"errors"

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

	// 创建用户成功后，给新用户发送欢迎系统消息
	_ = l.svcCtx.MessageRepository.Create(l.ctx, &model.Message{
		SenderID:   0,
		ReceiverID: ID,
		Type:       model.MessageSystem,
		EntityID:   0,
		Content:    "亲爱的朋友，非常欢迎您加入 CoderHub 社区！在这里，您可以尽情分享您的编程知识、经验和项目成果，与广大编程爱好者交流互动，共同提升技术水平。",
		IsRead:     false,
	})

	return &coderhub.CreateUserResponse{
		UserId: ID,
	}, nil
}
