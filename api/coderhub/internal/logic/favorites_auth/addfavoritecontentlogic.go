package favorites_auth

import (
	"coderhub/conf"
	"coderhub/model"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
	"context"
	"fmt"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddFavoriteContentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewAddFavoriteContentLogic 添加收藏内容
func NewAddFavoriteContentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFavoriteContentLogic {
	return &AddFavoriteContentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddFavoriteContentLogic) AddFavoriteContent(req *types.CreateFavorReq) (resp *types.CreateFavorResp, err error) {
	userID, err := utils.GetUserID(l.ctx)
	if err != nil {
		return l.errorResp(err)
	}
	response, err := l.svcCtx.FavoriteContentService.CreateFavor(l.ctx, &coderhub.CreateFavorRequest{
		UserId:        userID,
		FavorFolderId: utils.String2Int(req.FoldId),
		EntityId:      utils.String2Int(req.EntityId),
		EntityType:    req.EntityType,
	})
	if err != nil {
		return l.errorResp(err)
	}
	err = l.SendMessage(l.ctx, utils.String2Int(req.EntityId), userID, req.EntityType)
	if err != nil {
		return l.errorResp(err)
	}
	return l.successResp(response)
}

// SendMessage 发送收藏消息
func (l *AddFavoriteContentLogic) SendMessage(ctx context.Context, entityId, userId int64, entityType string) (err error) {
	if entityType == model.ArticleType {
		article, err := l.svcCtx.ArticlesService.GetArticle(ctx, &coderhub.GetArticleRequest{
			Id:     entityId,
			UserId: 0,
		})
		if err != nil {
			return err
		}
		userInfo, err := l.svcCtx.UserService.GetUserInfo(l.ctx, &coderhub.GetUserInfoRequest{
			UserId:        userId,
			RequestUserId: 0,
		})
		if err != nil {
			return err
		}
		var title string
		var _type string
		if article.Article.Title != "" {
			title = article.Article.Title
			_type = "文章"
		} else {
			title = "沸点"
			_type = ""
		}
		_, err = l.svcCtx.MessageService.CreateMessage(ctx, &coderhub.CreateMessageRequest{
			SenderId:   userId,
			ReceiverId: article.Author.UserId,
			Type:       model.MessageFavorite,
			EntityId:   entityId,
			Content:    fmt.Sprintf("用户 <a className=\"font-bold inline-block mx-2\" href=\"/user/%s\" target=\"_blank\">%s</a> 收藏了你的 %s <a className=\"font-bold inline-block mx-2\" href=\"/post/%s\" target=\"_blank\">《%s》</a>", utils.Int2String(userInfo.UserId), userInfo.UserName, _type, utils.Int2String(article.Article.Id), title),
		})
		if err != nil {
			fmt.Printf("发送收藏消息失败：%s\n", err.Error())
			return err
		}
	}
	return nil
}

func (l *AddFavoriteContentLogic) successResp(createFavorResp *coderhub.CreateFavorResponse) (*types.CreateFavorResp, error) {
	return &types.CreateFavorResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: createFavorResp.Success,
	}, nil
}

func (l *AddFavoriteContentLogic) errorResp(err error) (*types.CreateFavorResp, error) {
	return &types.CreateFavorResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: false,
	}, nil
}
