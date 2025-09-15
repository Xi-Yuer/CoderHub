package school_exp_public

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListSchoolExpLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewListSchoolExpLogic 获取学校经历列表
func NewListSchoolExpLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSchoolExpLogic {
	return &ListSchoolExpLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListSchoolExpLogic) ListSchoolExp(req *types.GetSchoolExpListReq) (resp *types.GetSchoolExpListResp, err error) {
	list, err := l.svcCtx.SchoolExpService.GetSchoolExpList(l.ctx, &coderhub.GetSchoolExpListRequest{
		Page:      req.Page,
		PageSize:  req.PageSize,
		Education: req.Education,
		School:    req.School,
		Major:     req.Major,
		WorkExp:   req.WorkExp,
	})
	if err != nil {
		return l.errorResp(err)
	}

	var schoolExpList []*types.SchoolExp
	for _, v := range list.SchoolExps {
		schoolExpList = append(schoolExpList, &types.SchoolExp{
			ID:        utils.Int2String(v.Id),
			Education: v.Education,
			School:    v.School,
			Major:     v.Major,
			WorkExp:   v.WorkExp,
			Content:   v.Content,
			UserId:    utils.Int2String(v.UserId),
			CreatedAt: v.CreateTime,
			UpdatedAt: v.UpdateTime,
		})
	}

	return l.successResp(&types.SchoolExpList{
		Total: list.Total,
		List:  schoolExpList,
	})
}
func (l *ListSchoolExpLogic) errorResp(err error) (resp *types.GetSchoolExpListResp, err1 error) {
	return &types.GetSchoolExpListResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: nil,
	}, nil
}

func (l *ListSchoolExpLogic) successResp(data *types.SchoolExpList) (resp *types.GetSchoolExpListResp, err error) {
	return &types.GetSchoolExpListResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: data,
	}, nil
}
