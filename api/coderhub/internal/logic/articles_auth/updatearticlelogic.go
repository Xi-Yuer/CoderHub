package articles_auth

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
	"context"
	"errors"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新文章
func NewUpdateArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateArticleLogic {
	return &UpdateArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateArticleLogic) UpdateArticle(req *types.UpdateArticleReq) (resp *types.UpdateArticleResp, err error) {
	userId, err := utils.GetUserID(l.ctx)
	if err != nil {
		return l.errorResp(err), nil
	}
	ctx := utils.SetUserMetaData(l.ctx) // 设置元数据
	// 1. 参数验证
	if err := utils.NewValidator().ArticleID(req.Id).Check(); err != nil {
		return l.errorResp(err), nil
	}

	// 3. 获取文章信息
	article, err := l.getArticle(req.Id)
	if err != nil {
		return l.errorResp(err), nil
	}

	// 4. 验证文章作者权限
	if article.Article.AuthorId != userId {
		return l.errorResp(errors.New("非法操作")), nil
	}

	// 5. 更新文章逻辑
	if _, err := l.updateArticle(ctx, req); err != nil {
		return l.errorResp(err), nil
	}

	return l.successResp(), nil
}

// 统一的错误响应处理
func (l *UpdateArticleLogic) errorResp(err error) *types.UpdateArticleResp {
	return &types.UpdateArticleResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: false,
	}
}

// 统一的成功响应处理
func (l *UpdateArticleLogic) successResp() *types.UpdateArticleResp {
	return &types.UpdateArticleResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: true,
	}
}

// 获取文章信息
func (l *UpdateArticleLogic) getArticle(articleId int64) (*coderhub.GetArticleResponse, error) {
	return l.svcCtx.ArticlesService.GetArticle(l.ctx, &coderhub.GetArticleRequest{
		Id: articleId,
	})
}

// 更新文章
func (l *UpdateArticleLogic) updateArticle(ctx context.Context, req *types.UpdateArticleReq) (*coderhub.UpdateArticleResponse, error) {
	return l.svcCtx.ArticlesService.UpdateArticle(ctx, &coderhub.UpdateArticleRequest{
		Id:           req.Id,
		Title:        req.Title,
		Content:      req.Content,
		Summary:      l.generateSummary(req.Summary, req.Content),
		ImageIds:     req.ImageIds,
		CoverImageId: req.CoverImageID,
		Tags:         req.Tags,
		Status:       req.Status,
	})
}

func (l *UpdateArticleLogic) generateSummary(summary, content string) string {
	if summary == "" {
		// TODO: 实现文章摘要生成逻辑
		// 可以截取前N个字符或使用AI生成摘要
		const summaryLength = 200
		if len(content) > summaryLength {
			return content[:summaryLength] + "..."
		}
		return content
	} else {
		return summary
	}
}
