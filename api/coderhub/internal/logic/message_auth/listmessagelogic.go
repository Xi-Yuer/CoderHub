package message_auth

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewListMessageLogic 获取消息列表
func NewListMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMessageLogic {
	return &ListMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMessageLogic) ListMessage(req *types.GetMessageListReq) (resp *types.GetMessageListResp, err error) {
	list, err := l.svcCtx.MessageService.GetMessageList(l.ctx, &coderhub.GetMessageListRequest{
		UserId:   utils.String2Int(req.RequestUserId),
		Type:     int32(req.Type),
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
	})
	if err != nil {
		return l.errorResp(err)
	}

	// 去重用户 ID
	userIdSet := make(map[int64]struct{})
	for _, v := range list.List {
		userIdSet[v.SenderId] = struct{}{}
	}

	userIdList := make([]int64, 0, len(userIdSet))
	for id := range userIdSet {
		userIdList = append(userIdList, id)
	}

	batchGetUserByID, err := l.svcCtx.UserService.BatchGetUserByID(l.ctx, &coderhub.BatchGetUserByIDRequest{
		UserIds: userIdList,
	})
	if err != nil {
		return nil, err
	}

	// 构建用户信息映射
	userInfoMap := make(map[int64]*coderhub.UserInfo)
	for _, user := range batchGetUserByID.UserInfos {
		userInfoMap[user.UserId] = user
	}

	messageList := make([]*types.Message, 0, len(list.List))
	for _, v := range list.List {
		info, ok := userInfoMap[v.SenderId]
		if !ok {
			info = &coderhub.UserInfo{}
		}

		messageList = append(messageList, &types.Message{
			ID:         utils.Int2String(v.Id),
			SenderID:   utils.Int2String(v.SenderId),
			ReceiverID: utils.Int2String(v.ReceiverId),
			Type:       int64(v.Type),
			EntityID:   utils.Int2String(v.EntityId),
			Content:    v.Content,
			IsRead:     v.IsRead,
			SenderInfo: &types.UserInfo{
				Id:          utils.Int2String(info.UserId),
				Username:    info.UserName,
				Nickname:    info.NickName,
				Email:       info.Email,
				Phone:       info.Phone,
				Avatar:      info.Avatar,
				Gender:      info.Gender,
				Age:         info.Age,
				Status:      info.Status,
				IsAdmin:     info.IsAdmin,
				CreateAt:    info.CreatedAt,
				UpdateAt:    info.UpdatedAt,
				FollowCount: info.FollowCount,
				FansCount:   info.FollowerCount,
				IsFollowed:  info.IsFollowed,
			},
			CreatedAt: v.CreateTime,
			UpdatedAt: v.CreateTime,
		})
	}

	return l.successResp(&types.MessageList{
		List:  messageList,
		Total: int64(list.Total),
	})
}

func (l *ListMessageLogic) errorResp(err error) (resp *types.GetMessageListResp, err1 error) {
	return &types.GetMessageListResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: nil,
	}, nil
}

func (l *ListMessageLogic) successResp(data *types.MessageList) (resp *types.GetMessageListResp, err error) {
	return &types.GetMessageListResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: data,
	}, nil
}
