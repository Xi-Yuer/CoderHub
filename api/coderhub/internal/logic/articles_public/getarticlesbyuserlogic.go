package articles_public

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
	"context"
	"fmt"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetArticlesByUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetArticlesByUserLogic 获取某个用户的文章列表
func NewGetArticlesByUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticlesByUserLogic {
	return &GetArticlesByUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetArticlesByUserLogic) GetArticlesByUser(req *types.GetArticlesByUserReq) (resp *types.GetArticlesByUserResp, err error) {
	// 获取作者文章列表ID
	articles, err := l.svcCtx.ArticlesService.ListArticleIDsByAuthor(l.ctx, &coderhub.ListAuthorArticlesRequest{
		AuthorId: utils.String2Int(req.AuthorID),
		Page:     int64(req.Page),
		PageSize: int64(req.PageSize),
		Type:     req.Type,
	})
	if err != nil {
		return l.errorResp(err)
	}

	fmt.Println("articlesIDs: ", articles.Ids)
	if len(articles.Ids) == 0 {
		return &types.GetArticlesByUserResp{
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
		UserId: utils.String2Int(req.RequestUserID),
	})
	if err != nil {
		return l.errorResp(err)
	}

	if len(response.Articles) == 0 {
		return &types.GetArticlesByUserResp{
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

	return &types.GetArticlesByUserResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: types.GetArticleResponse{
			List:  list,
			Total: int32(articles.Total),
		},
	}, nil
}
func (l *GetArticlesByUserLogic) errorResp(err error) (*types.GetArticlesByUserResp, error) {
	return &types.GetArticlesByUserResp{
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

func (l *GetArticlesByUserLogic) convertToArticleType(article *coderhub.GetArticleResponse) *types.GetArticle {
	if article == nil {
		fmt.Println("article is nil")
		return nil
	}
	if article.Article == nil {
		fmt.Println("article.Article is nil")
		return nil
	}
	if article.Author == nil {
		fmt.Println("article.Author is nil")
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
