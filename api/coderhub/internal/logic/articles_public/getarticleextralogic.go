package articles_public

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetArticleExtraLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetArticleExtraLogic 获取文章附加信息
func NewGetArticleExtraLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticleExtraLogic {
	return &GetArticleExtraLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetArticleExtraLogic) GetArticleExtra(req *types.GetArticleReq) (resp *types.GetArticleExtraResp, err error) {
	extra, err := l.svcCtx.ArticlesService.GetArticlesExtra(l.ctx, &coderhub.GetArticleRequest{
		Id:     utils.String2Int(req.Id),
		UserId: utils.String2Int(req.RequestUserID),
	})
	if err != nil {
		return l.errorResp(err)
	}
	return l.successResp(extra)
}

func (l *GetArticleExtraLogic) successResp(data *coderhub.ArticleAdditionalInfo) (*types.GetArticleExtraResp, error) {
	return &types.GetArticleExtraResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: &types.ArticleExtra{
			ID:            utils.Int2String(data.Id),
			ViewCount:     data.ViewCount,
			LikeCount:     data.LikeCount,
			CommentCount:  data.CommentCount,
			IsLiked:       data.IsLicked,
			IsFavorited:   data.IsFavorite,
			FavoriteCount: data.FavoriteCount,
		},
	}, nil
}

func (l *GetArticleExtraLogic) errorResp(err error) (*types.GetArticleExtraResp, error) {
	return &types.GetArticleExtraResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
	}, err
}
