package academic_public

import (
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAcademicNavigatorLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取学术导航
func NewGetAcademicNavigatorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAcademicNavigatorLogic {
	return &GetAcademicNavigatorLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAcademicNavigatorLogic) GetAcademicNavigator(req *types.GetAcademicNavigatorReq) (resp *types.GetAcademicNavigatorResp, err error) {
	var respAcademicNavigator *coderhub.GetAcademicNavigatorResponse
	respAcademicNavigator, err = l.svcCtx.AcademicService.GetAcademicNavigator(l.ctx, &coderhub.GetAcademicNavigatorRequest{
		UserId:    req.UserId,
		Education: req.Education,
		Major:     req.Major,
		School:    req.School,
		WorkExp:   req.WorkExp,
		Content:   req.Content,
		Page:      req.Page,
		PageSize:  req.PageSize,
	})
	if err != nil {
		return l.errorResp(err)
	}

	academicNavigator := make([]types.AcademicNavigator, len(respAcademicNavigator.AcademicNavigator))
	for i, v := range respAcademicNavigator.AcademicNavigator {
		academicNavigator[i] = types.AcademicNavigator{
			Id:        v.Id,
			UserId:    v.UserId,
			Education: v.Education,
			Content:   v.Content,
			Major:     v.Major,
			School:    v.School,
			WorkExp:   v.WorkExp,
			LikeCount: v.LikeCount,
		}
	}

	return l.successResp(academicNavigator, respAcademicNavigator.Total)
}

func (l *GetAcademicNavigatorLogic) errorResp(err error) (*types.GetAcademicNavigatorResp, error) {
	return &types.GetAcademicNavigatorResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: nil,
	}, nil
}

func (l *GetAcademicNavigatorLogic) successResp(academicNavigators []types.AcademicNavigator, total int64) (*types.GetAcademicNavigatorResp, error) {
	return &types.GetAcademicNavigatorResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: &types.AcademicList{
			Total: total,
			List:  academicNavigators,
		},
	}, nil
}
