package follow_auth

import (
	"coderhub/model"
	"context"
	"fmt"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewFollowUserLogic 关注用户
func NewFollowUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowUserLogic {
	return &FollowUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FollowUserLogic) FollowUser(req *types.FollowUserReq) (resp *types.FollowUserResp, err error) {
	UserID, err := utils.GetUserID(l.ctx)
	if err != nil {
		return l.errorResp(err)
	}
	_, err = l.svcCtx.UserFollowService.CreateUserFollow(l.ctx, &coderhub.CreateUserFollowReq{
		FollowerId: UserID,
		FollowedId: utils.String2Int(req.FollowUserId),
	})
	if err != nil {
		return l.errorResp(err)
	}

	err = l.SendMessage(l.ctx, utils.String2Int(req.FollowUserId), UserID)
	if err != nil {
		return l.errorResp(err)
	}

	return l.successResp()
}

// SendMessage 发送关注消息
func (l *FollowUserLogic) SendMessage(ctx context.Context, entityId, userId int64) error {
	userInfo, err := l.svcCtx.UserService.GetUserInfo(l.ctx, &coderhub.GetUserInfoRequest{
		UserId:        userId,
		RequestUserId: 0,
	})
	if err != nil {
		fmt.Printf("获取用户信息失败：%s\n", err.Error())
		return err
	}
	_, err = l.svcCtx.MessageService.CreateMessage(ctx, &coderhub.CreateMessageRequest{
		SenderId:   userId,
		ReceiverId: entityId,
		Type:       model.MessageFollow,
		EntityId:   entityId,
		Content:    fmt.Sprintf("用户 <a className=\"font-bold inline-block mx-2\" href=\"/user/%s\" target=\"_blank\">%s</a> 关注了你", utils.Int2String(userInfo.UserId), userInfo.UserName),
	})
	if err != nil {
		fmt.Printf("发送关注消息失败：%s\n", err.Error())
		return err
	}
	return nil
}

func (l *FollowUserLogic) successResp() (*types.FollowUserResp, error) {
	return &types.FollowUserResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: true,
	}, nil
}

func (l *FollowUserLogic) errorResp(err error) (*types.FollowUserResp, error) {
	return &types.FollowUserResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
	}, nil
}
