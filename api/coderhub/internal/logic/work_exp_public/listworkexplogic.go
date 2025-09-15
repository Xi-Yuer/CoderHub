package work_exp_public

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListWorkExpLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewListWorkExpLogic 获取工作经历列表
func NewListWorkExpLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListWorkExpLogic {
	return &ListWorkExpLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListWorkExpLogic) ListWorkExp(req *types.GetWorkExpListReq) (resp *types.GetWorkExpListResp, err error) {
	list, err := l.svcCtx.WorkExpService.GetWorkExpList(l.ctx, &coderhub.GetWorkExpListRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Company:  req.Company,
		WorkExp:  req.WorkExp,
		Position: req.Position,
		Region:   req.Region,
	})
	if err != nil {
		return l.errorResp(err)
	}

	var workExpList []*types.WorkExp
	for _, v := range list.WorkExps {
		workExpList = append(workExpList, &types.WorkExp{
			ID:        utils.Int2String(v.Id),
			WorkExp:   v.WorkExp,
			Company:   v.Company,
			Region:    v.Region,
			Position:  v.Position,
			Content:   v.Content,
			UserID:    utils.Int2String(v.UserId),
			CreatedAt: v.CreateTime,
			UpdatedAt: v.UpdateTime,
		})
	}

	return l.successResp(&types.WorkExpList{
		Total: list.Total,
		List:  workExpList,
	})
}
func (l *ListWorkExpLogic) errorResp(err error) (resp *types.GetWorkExpListResp, err1 error) {
	return &types.GetWorkExpListResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: nil,
	}, nil
}

func (l *ListWorkExpLogic) successResp(data *types.WorkExpList) (resp *types.GetWorkExpListResp, err error) {
	return &types.GetWorkExpListResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: data,
	}, nil
}
