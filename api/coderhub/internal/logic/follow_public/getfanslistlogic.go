package follow_public

import (
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFansListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取粉丝列表
func NewGetFansListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFansListLogic {
	return &GetFansListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFansListLogic) GetFansList(req *types.GetFansListReq) (resp *types.GetFansListResp, err error) {
	userFansResp, err := l.svcCtx.UserFollowService.GetUserFans(l.ctx, &coderhub.GetUserFansReq{
		FollowedId: req.UserId,
		Page:       int32(req.Page),
		PageSize:   int32(req.PageSize),
	})
	if err != nil {
		return l.errorResp(err)
	}

	return l.successResp(userFansResp)
}

func (l *GetFansListLogic) successResp(userFansResp *coderhub.GetUserFansResp) (*types.GetFansListResp, error) {
	userFansList := make([]types.UserInfo, 0, len(userFansResp.UserFans))
	for _, userFan := range userFansResp.UserFans {
		userFansList = append(userFansList, types.UserInfo{
			Id:       userFan.UserId,
			Username: userFan.UserName,
			Nickname: userFan.NickName,
			Email:    userFan.Email,
			Phone:    userFan.Phone,
			Avatar:   userFan.Avatar,
			Gender:   userFan.Gender,
			Age:      userFan.Age,
			Status:   userFan.Status,
			IsAdmin:  userFan.IsAdmin,
			CreateAt: userFan.CreatedAt,
			UpdateAt: userFan.UpdatedAt,
		})
	}

	return &types.GetFansListResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: types.FollowList{
			List:  userFansList,
			Total: userFansResp.Total,
		},
	}, nil
}

func (l *GetFansListLogic) errorResp(err error) (*types.GetFansListResp, error) {
	return &types.GetFansListResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: types.FollowList{
			List:  nil,
			Total: 0,
		},
	}, nil
}
