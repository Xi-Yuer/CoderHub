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

type GetArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetArticleLogic 获取文章
func NewGetArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticleLogic {
	return &GetArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetArticleLogic) GetArticle(req *types.GetArticleReq) (resp *types.GetArticleResp, err error) {
	if err := utils.NewValidator().ArticleID(utils.String2Int(req.Id)).Check(); err != nil {
		return l.errorResp(err), nil
	}
	article, err := l.getArticle(utils.String2Int(req.Id), utils.String2Int(req.RequestUserID))
	if err != nil {
		return l.errorResp(err), nil
	}

	go l.incrementViewCount(utils.String2Int(req.Id))

	return l.successResp(article), nil
}

func (l *GetArticleLogic) errorResp(err error) *types.GetArticleResp {
	return &types.GetArticleResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: nil,
	}
}

func (l *GetArticleLogic) successResp(article *coderhub.GetArticleResponse) *types.GetArticleResp {
	return &types.GetArticleResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: l.convertToArticleType(article),
	}
}

func (l *GetArticleLogic) getArticle(articleId, userID int64) (*coderhub.GetArticleResponse, error) {
	return l.svcCtx.ArticlesService.GetArticle(l.ctx, &coderhub.GetArticleRequest{
		Id:     articleId,
		UserId: userID,
	})
}

func (l *GetArticleLogic) convertToArticleType(article *coderhub.GetArticleResponse) *types.GetArticle {
	if article == nil {
		return nil
	}

	// 将图片模型转换为图片URL
	imageUrls := make([]string, 0)
	if article.Article.Images != nil {
		for _, image := range article.Article.Images {
			if image != nil {
				imageUrls = append(imageUrls, image.Url)
			}
		}
	}
	l.Logger.Info("API: 获取文章配图成功, 配图数量:", len(imageUrls))

	coverImageUrl := ""
	if article.Article.CoverImage != nil {
		coverImageUrl = article.Article.CoverImage.Url
	}
	l.Logger.Info("API: 获取文章封面成功, 封面URL:", coverImageUrl)

	return &types.GetArticle{
		Article: &types.Article{
			Id:           utils.Int2String(article.Article.Id),
			Type:         article.Article.Type,
			Title:        article.Article.Title,
			Content:      article.Article.Content,
			Summary:      article.Article.Summary,
			ImageUrls:    imageUrls,
			CoverImage:   &coverImageUrl,
			AuthorId:     utils.Int2String(article.Article.AuthorId),
			Tags:         article.Article.Tags,
			ViewCount:    article.Article.ViewCount,
			LikeCount:    article.Article.LikeCount,
			IsLiked:      article.Article.IsLicked,
			IsFavorited:  article.Article.IsFavorite,
			CommentCount: article.Article.CommentCount,
			Status:       article.Article.Status,
			CreatedAt:    article.Article.CreatedAt,
			UpdatedAt:    article.Article.UpdatedAt,
		},
		Author: &types.UserInfo{
			Id:       utils.Int2String(article.Author.UserId),
			Username: article.Author.UserName,
			Nickname: article.Author.NickName,
			Email:    article.Author.Email,
			Phone:    article.Author.Phone,
			Avatar:   article.Author.Avatar,
			Gender:   article.Author.Gender,
			Age:      article.Author.Age,
			Status:   article.Author.Status,
			IsAdmin:  article.Author.IsAdmin,
			CreateAt: article.Author.CreatedAt,
			UpdateAt: article.Author.UpdatedAt,
		},
	}
}

func (l *GetArticleLogic) incrementViewCount(articleId int64) {
	// TODO: 实现浏览量更新逻辑
	// 可以使用 Redis 进行计数，定期同步到数据库
}
