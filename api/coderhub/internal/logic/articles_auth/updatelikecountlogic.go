package articles_auth

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

type UpdateLikeCountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUpdateLikeCountLogic 更新文章点赞数
func NewUpdateLikeCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLikeCountLogic {
	return &UpdateLikeCountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLikeCountLogic) UpdateLikeCount(req *types.UpdateLikeCountReq) (resp *types.UpdateLikeCountResp, err error) {
	userId, err := utils.GetUserID(l.ctx)
	if err != nil {
		return l.errorResp(err), nil
	}

	if err := utils.NewValidator().ArticleID(utils.String2Int(req.Id)).Check(); err != nil {
		return l.errorResp(err), nil
	}

	if _, err := l.svcCtx.ArticlesService.UpdateLikeCount(utils.SetUserMetaData(l.ctx), &coderhub.UpdateLikeCountRequest{
		Id:     utils.String2Int(req.Id),
		UserId: userId,
	}); err != nil {
		return l.errorResp(err), nil
	}
	// 发送点赞消息
	if req.Trigger {
		err = l.SendMessage(l.ctx, utils.String2Int(req.Id), userId)
		if err != nil {
			return l.errorResp(err), nil
		}
	}
	return l.successResp(), nil
}

// SendMessage 发送点赞消息
func (l *UpdateLikeCountLogic) SendMessage(ctx context.Context, entityId, userId int64) error {
	article, err := l.svcCtx.ArticlesService.GetArticle(l.ctx, &coderhub.GetArticleRequest{
		Id:     entityId,
		UserId: 0,
	})
	if err != nil {
		fmt.Printf("获取文章失败: %v", err)
		return err
	}
	userInfo, err := l.svcCtx.UserService.GetUserInfo(l.ctx, &coderhub.GetUserInfoRequest{
		UserId:        article.Article.AuthorId,
		RequestUserId: 0,
	})
	if err != nil {
		fmt.Printf("获取用户信息失败：%s\n", err.Error())
		return err
	}
	var content string
	var title string
	if article.Article.Title != "" {
		title = article.Article.Title
	} else {
		title = "沸点"
	}
	if article.Article.Type == model.ArticleType {
		content = fmt.Sprintf("用户 <a className=\"font-bold inline-block mx-2\" href=\"/user/%s\" target=\"_blank\">%s</a> 点赞了你的文章 <a className=\"font-bold inline-block mx-2\" href=\"/post/%s\" target=\"_blank\">《%s》</a>", utils.Int2String(userInfo.UserId), userInfo.UserName, utils.Int2String(article.Article.Id), title)
	}
	if article.Article.Type == model.MicroPostType {
		content = fmt.Sprintf("用户 <a className=\"font-bold inline-block mx-2\" href=\"/user/%s\" target=\"_blank\">%s</a> 点赞了你的帖子", utils.Int2String(userInfo.UserId), userInfo.UserName)
	}
	_, err = l.svcCtx.MessageService.CreateMessage(ctx, &coderhub.CreateMessageRequest{
		SenderId:   userId,
		ReceiverId: article.Author.UserId,
		Type:       model.MessageLike,
		EntityId:   article.Article.Id,
		Content:    content,
	})
	if err != nil {
		fmt.Printf("发送点赞消息失败：%s\n", err.Error())
		return err
	}
	return nil
}

func (l *UpdateLikeCountLogic) errorResp(err error) *types.UpdateLikeCountResp {
	return &types.UpdateLikeCountResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: false,
	}
}

func (l *UpdateLikeCountLogic) successResp() *types.UpdateLikeCountResp {
	return &types.UpdateLikeCountResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: "success",
		},
		Data: true,
	}
}
