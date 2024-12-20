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

// 获取文章
func NewGetArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticleLogic {
	return &GetArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetArticleLogic) GetArticle(req *types.GetArticleReq) (resp *types.GetArticleResp, err error) {
	if err := utils.NewValidator().ArticleID(req.Id).Check(); err != nil {
		return l.errorResp(err), nil
	}

	article, err := l.getArticle(req.Id)
	if err != nil {
		return l.errorResp(err), nil
	}

	go l.incrementViewCount(req.Id)

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
		Data: l.convertToArticleType(article.Article),
	}
}

func (l *GetArticleLogic) getArticle(articleId int64) (*coderhub.GetArticleResponse, error) {
	return l.svcCtx.ArticlesService.GetArticle(l.ctx, &coderhub.GetArticleRequest{
		Id: articleId,
	})
}

func (l *GetArticleLogic) convertToArticleType(article *coderhub.Article) *types.Article {
	if article == nil {
		return nil
	}

	// 将图片模型转换为图片URL
	imageUrls := make([]string, 0)
	if article.Images != nil {
		for _, image := range article.Images {
			if image != nil {
				imageUrls = append(imageUrls, image.Url)
			}
		}
	}
	l.Logger.Info("API: 获取文章配图成功, 配图数量:", len(imageUrls))

	coverImageUrl := ""
	if article.CoverImage != nil {
		coverImageUrl = article.CoverImage.Url
	}
	l.Logger.Info("API: 获取文章封面成功, 封面URL:", coverImageUrl)

	return &types.Article{
		Id:           article.Id,
		Type:         article.Type,
		Title:        article.Title,
		Content:      article.Content,
		Summary:      article.Summary,
		CoverImage:   &coverImageUrl,
		ImageUrls:    imageUrls,
		AuthorId:     article.AuthorId,
		Tags:         article.Tags,
		ViewCount:    article.ViewCount,
		LikeCount:    article.LikeCount,
		CommentCount: article.CommentCount,
		Status:       article.Status,
		CreatedAt:    article.CreatedAt,
		UpdatedAt:    article.UpdatedAt,
	}
}

func (l *GetArticleLogic) incrementViewCount(articleId int64) {
	// TODO: 实现浏览量更新逻辑
	// 可以使用 Redis 进行计数，定期同步到数据库
}
