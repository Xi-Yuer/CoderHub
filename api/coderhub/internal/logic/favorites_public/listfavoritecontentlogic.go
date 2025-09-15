package favorites_public

import (
	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListFavoriteContentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewListFavoriteContentLogic 获取收藏内容列表
func NewListFavoriteContentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFavoriteContentLogic {
	return &ListFavoriteContentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListFavoriteContentLogic) ListFavoriteContent(req *types.GetFavorListReq) (resp *types.GetFavorListResp, err error) {
	list, err := l.svcCtx.FavoriteContentService.GetFavorList(l.ctx, &coderhub.GetFavorListRequest{
		Page:          req.Page,
		PageSize:      req.PageSize,
		UserId:        utils.String2Int(req.UserId),
		EntityType:    req.EntityType,
		FavorFolderId: utils.String2Int(req.FavorFoldId),
	})
	if err != nil {
		return l.errorResp(err)
	}

	// 确保 list 不为 nil
	if list == nil {
		logx.Errorf("GetFavorList returned nil list")
		return l.successResp(types.FavorList{Total: 0, List: []*types.Favor{}})
	}

	// 确保 list.Favors 不为 nil
	if list.Favors == nil {
		list.Favors = []*coderhub.Favor{}
	}

	response := make([]*types.Favor, 0, len(list.Favors))

	for _, v := range list.Favors {
		if v == nil || v.EntityValue == nil || v.EntityValue.User == nil {
			continue
		}
		response = append(response, &types.Favor{
			ID:          utils.Int2String(v.Id),
			CreateUser:  utils.Int2String(v.UserId),
			FavorFoldId: utils.Int2String(v.FavorFolderId),
			EntityId:    utils.Int2String(v.EntityId),
			EntityValue: types.EntityPreviewValue{
				EntityId:   utils.Int2String(v.EntityValue.EntityId),
				Title:      v.EntityValue.Title,
				Content:    v.EntityValue.Content,
				EntityType: v.EntityValue.EntityType,
				CoverImage: v.EntityValue.CoverImage,
				UserInfo: types.UserInfo{
					Id:       utils.Int2String(v.EntityValue.User.UserId),
					Username: v.EntityValue.User.UserName,
					Nickname: v.EntityValue.User.NickName,
					Email:    v.EntityValue.User.Email,
					Phone:    v.EntityValue.User.Phone,
					Avatar:   v.EntityValue.User.Avatar,
					Gender:   v.EntityValue.User.Gender,
					Age:      v.EntityValue.User.Age,
					Status:   v.EntityValue.User.Status,
					IsAdmin:  v.EntityValue.User.IsAdmin,
					CreateAt: v.EntityValue.User.CreatedAt,
					UpdateAt: v.EntityValue.User.UpdatedAt,
				},
			},
			EntityType: v.EntityType,
			CreatedAt:  v.CreateTime,
		})
	}

	return l.successResp(types.FavorList{
		Total: list.Total,
		List:  response,
	})
}
func (l *ListFavoriteContentLogic) successResp(list types.FavorList) (*types.GetFavorListResp, error) {
	return &types.GetFavorListResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: list,
	}, nil
}

func (l *ListFavoriteContentLogic) errorResp(err error) (*types.GetFavorListResp, error) {
	return &types.GetFavorListResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: types.FavorList{},
	}, nil
}
