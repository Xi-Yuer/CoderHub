package favorservicelogic

import (
	"coderhub/model"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFavorListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFavorListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFavorListLogic {
	return &GetFavorListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetFavorList 获取收藏夹列表
func (l *GetFavorListLogic) GetFavorList(in *coderhub.GetFavorListRequest) (*coderhub.GetFavorListResponse, error) {
	list, total, err := l.svcCtx.UserFavorEntityRepository.GetList(l.ctx, &model.UserFavor{
		UserId:      in.UserId,
		FavorFoldId: in.FavorFolderId,
		EntityType:  in.EntityType,
	}, int(in.Page), int(in.PageSize))
	if err != nil {
		return nil, err
	}

	// 实体id
	ids := make([]int64, 0, len(list))
	for _, v := range list {
		ids = append(ids, v.EntityId)
	}

	// 获取收藏夹内容
	entityValue := make(map[int64]*coderhub.FavorPreview)

	// 获取收藏夹内容详情
	if in.EntityType == model.UserFavorRelationArticle {
		articles, err := l.svcCtx.ArticleRepository.BatchGetArticle(ids)
		if err != nil {
			return nil, err
		}
		for _, v := range articles {
			entityValue[v.ArticleID] = &coderhub.FavorPreview{
				EntityId:   v.ArticleID,
				Title:      v.Title,
				Content:    v.Summary,
				EntityType: model.UserFavorRelationArticle,
				CoverImage: v.CoverImage,
				User: &coderhub.UserInfo{
					UserId:   v.AuthID,
					UserName: v.AuthName,
					Avatar:   v.Avatar,
				},
			}
		}
	}

	// 获取收藏夹内容详情
	if in.EntityType == model.UserFavorRelationQuestion {
		banks, err := l.svcCtx.QuestionBankRepository.BatchGetQuestion(l.ctx, ids)
		if err != nil {
			return nil, err
		}
		for _, v := range banks {
			entityValue[v.ID] = &coderhub.FavorPreview{
				EntityId:   v.ID,
				Title:      v.Name,
				Content:    v.Description,
				EntityType: model.UserFavorRelationQuestion,
				CoverImage: v.CoverImage,
				User: &coderhub.UserInfo{
					UserId:   v.CreateUserID,
					UserName: v.Name,
					Avatar:   v.Avatar,
				},
			}
		}
	}

	favors := make([]*coderhub.Favor, 0, len(list))
	for _, v := range list {
		entity, ok := entityValue[v.EntityId]
		if !ok || entity == nil {
			break
		}
		favor := &coderhub.Favor{
			Id:            int64(v.ID),
			UserId:        v.UserId,
			FavorFolderId: v.FavorFoldId,
			EntityId:      v.EntityId,
			EntityType:    v.EntityType,
			CreateTime:    v.CreatedAt.Unix(),
			EntityValue:   entity,
		}
		favors = append(favors, favor)
	}
	fmt.Println("len(favors)", len(favors))

	return &coderhub.GetFavorListResponse{
		Favors: favors,
		Total:  total,
	}, nil
}
