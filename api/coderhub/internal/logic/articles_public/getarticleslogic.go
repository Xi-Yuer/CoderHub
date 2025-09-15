package articles_public

import (
	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetArticlesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetArticlesLogic 获取文章列表
func NewGetArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticlesLogic {
	return &GetArticlesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetArticlesLogic) GetArticles(req *types.GetArticlesReq) (resp *types.GetArticlesResp, err error) {
	// 获取推荐文章列表ID
	articles, err := l.svcCtx.ArticlesService.ListRecommendedArticles(l.ctx, &coderhub.ListRecommendedArticlesRequest{
		Type:       req.Type,
		Page:       int64(req.Page),
		PageSize:   int64(req.PageSize),
		CategoryId: utils.String2Int(req.CategoryID),
	})
	if err != nil {
		return l.errorResp(err)
	}

	if len(articles.Ids) == 0 {
		return &types.GetArticlesResp{
			Response: types.Response{
				Code:    conf.HttpCode.HttpStatusOK,
				Message: "No recommended articles found",
			},
			Data: types.GetArticleResponse{
				List:  nil,
				Total: 0,
			},
		}, nil
	}

	// 获取文章列表详情
	response, err := l.svcCtx.ArticlesService.ListArticles(l.ctx, &coderhub.GetArticlesRequest{
		Ids:    articles.Ids,
		UserId: utils.String2Int(req.UserID),
	})
	if err != nil {
		return l.errorResp(err)
	}

	if len(response.Articles) == 0 {
		return &types.GetArticlesResp{
			Response: types.Response{
				Code:    conf.HttpCode.HttpStatusOK,
				Message: "No articles found for given IDs",
			},
			Data: types.GetArticleResponse{
				List:  nil,
				Total: 0,
			},
		}, nil
	}

	// 构建文章列表
	var list []*types.GetArticle
	for _, article := range response.Articles {
		fmt.Printf("article: %#v\n", article)
		if converted := l.convertToArticleType(article); converted != nil {
			list = append(list, converted)
		}
	}

	return &types.GetArticlesResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: types.GetArticleResponse{
			List:  list,
			Total: 0,
		},
	}, nil
}

func (l *GetArticlesLogic) errorResp(err error) (*types.GetArticlesResp, error) {
	return &types.GetArticlesResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: types.GetArticleResponse{
			List:  nil,
			Total: 0,
		},
	}, nil
}

func (l *GetArticlesLogic) convertToArticleType(article *coderhub.GetArticleResponse) *types.GetArticle {
	if article == nil {
		return nil
	}
	if article.Article == nil {
		return nil
	}
	if article.Author == nil {
		return nil
	}

	fmt.Printf("article.Article.Images: %#v\n", article.Article.Images)
	fmt.Printf("article.Article.CoverImage: %#v\n", article.Article.CoverImage)

	images := make([]string, len(article.Article.Images))
	for i, image := range article.Article.Images {
		if image != nil {
			images[i] = image.Url
		} else {
			fmt.Printf("article.Article.Images[%d] is nil\n", i)
		}
	}
	fmt.Printf("article.Article.IsLicked ==> %v\n", article.Article.IsLicked)
	return &types.GetArticle{
		Article: &types.Article{
			Id:        utils.Int2String(article.Article.Id),
			Type:      article.Article.Type,
			Title:     article.Article.Title,
			Content:   article.Article.Content,
			Summary:   article.Article.Summary,
			ImageUrls: images,
			CoverImage: func(img *coderhub.Image) *string {
				if img == nil {
					return nil
				}
				return &img.Url
			}(article.Article.CoverImage),
			AuthorId:     utils.Int2String(article.Author.UserId),
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
			Level:    article.Author.Level,
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
