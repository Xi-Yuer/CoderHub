package schoolexpservicelogic

import (
	"coderhub/model"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSchoolExpListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSchoolExpListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSchoolExpListLogic {
	return &GetSchoolExpListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSchoolExpListLogic) GetSchoolExpList(in *coderhub.GetSchoolExpListRequest) (*coderhub.GetSchoolExpListResponse, error) {
	list, i, err := l.svcCtx.SchoolExpRepository.List(l.ctx, &model.SchoolExp{
		Education: in.Education,
		Major:     in.Major,
		School:    in.School,
		WorkExp:   in.WorkExp,
	}, int64(in.Page), int64(in.PageSize))
	if err != nil {
		return nil, err
	}
	var schoolExps []*coderhub.SchoolExp
	for _, v := range list {
		schoolExps = append(schoolExps, &coderhub.SchoolExp{
			Id:         v.ID,
			Education:  v.Education,
			School:     v.School,
			Major:      v.Major,
			WorkExp:    v.WorkExp,
			Content:    v.Content,
			UserId:     v.UserId,
			CreateTime: v.CreateTime.Unix(),
			UpdateTime: v.UpdateTime.Unix(),
		})
	}

	return &coderhub.GetSchoolExpListResponse{
		SchoolExps: schoolExps,
		Total:      i,
	}, nil
}
