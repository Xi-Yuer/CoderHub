package articles_auth

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
	"context"
	"strconv"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除文章
func NewDeleteArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteArticleLogic {
	return &DeleteArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteArticleLogic) DeleteArticle(req *types.DeleteArticleReq) (resp *types.DeleteArticleResp, err error) {
	if err := utils.NewValidator().ArticleID(utils.String2Int(req.Id)).Check(); err != nil {
		return l.errorResp(err), nil
	}

	articleIdInt, err := strconv.ParseInt(req.Id, 10, 64)
	if err != nil {
		return l.errorResp(err), nil
	}

	if err := l.deleteArticle(utils.SetUserMetaData(l.ctx), articleIdInt); err != nil {
		return l.errorResp(err), nil
	}

	return l.successResp(), nil
}

func (l *DeleteArticleLogic) errorResp(err error) *types.DeleteArticleResp {
	return &types.DeleteArticleResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: false,
	}
}

func (l *DeleteArticleLogic) successResp() *types.DeleteArticleResp {
	return &types.DeleteArticleResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: true,
	}
}

func (l *DeleteArticleLogic) deleteArticle(ctx context.Context, articleId int64) error {
	userID, err := utils.GetUserID(l.ctx)
	if err != nil {
		return err
	}
	_, err = l.svcCtx.ArticlesService.DeleteArticle(ctx, &coderhub.DeleteArticleRequest{
		Id:     articleId,
		UserId: userID,
	})
	return err
}
