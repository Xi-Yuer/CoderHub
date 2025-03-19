package workexpservicelogic

import (
	"coderhub/model"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWorkExpListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetWorkExpListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWorkExpListLogic {
	return &GetWorkExpListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetWorkExpListLogic) GetWorkExpList(in *coderhub.GetWorkExpListRequest) (*coderhub.GetWorkExpListResponse, error) {
	list, i, err := l.svcCtx.WorkExpRepository.List(l.ctx, &model.WorkExp{
		Company:      in.Company,
		WorkDuration: in.WorkExp,
		Region:       in.Region,
		Position:     in.Position,
	}, int64(in.Page), int64(in.PageSize))
	if err != nil {
		return nil, err
	}

	var workExps []*coderhub.WorkExp
	for _, v := range list {
		workExps = append(workExps, &coderhub.WorkExp{
			Id:         v.ID,
			Company:    v.Company,
			Region:     v.Region,
			WorkExp:    v.WorkDuration,
			Content:    v.Content,
			UserId:     v.UserId,
			Position:   v.Position,
			CreateTime: v.CreateTime.Unix(),
			UpdateTime: v.UpdateTime.Unix(),
		})
	}

	return &coderhub.GetWorkExpListResponse{
		WorkExps: workExps,
		Total:    i,
	}, nil
}
