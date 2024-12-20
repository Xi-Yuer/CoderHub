package userservicelogic

import (
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchGetUserByIDLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchGetUserByIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchGetUserByIDLogic {
	return &BatchGetUserByIDLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BatchGetUserByID 批量获取用户信息
func (l *BatchGetUserByIDLogic) BatchGetUserByID(in *coderhub.BatchGetUserByIDRequest) (*coderhub.BatchGetUserByIDResponse, error) {
	l.Logger.Info("BatchGetUserByID", in.UserIds)
	users, err := l.svcCtx.UserRepository.BatchGetUserByID(in.UserIds)
	if err != nil {
		return nil, err
	}

	l.Logger.Info("users", users)

	userInfos := make([]*coderhub.UserInfo, len(users))
	for i, val := range users {
		userInfos[i] = &coderhub.UserInfo{
			UserId:    val.ID,
			UserName:  val.UserName,
			Avatar:    val.Avatar.String,
			Email:     val.Email.String,
			Gender:    val.Gender,
			Age:       val.Age,
			Phone:     val.Phone.String,
			NickName:  val.NickName.String,
			IsAdmin:   val.IsAdmin,
			Status:    val.Status,
			CreatedAt: val.CreatedAt.Unix(),
			UpdatedAt: val.UpdatedAt.Unix(),
		}
	}

	return &coderhub.BatchGetUserByIDResponse{UserInfos: userInfos}, nil
}
