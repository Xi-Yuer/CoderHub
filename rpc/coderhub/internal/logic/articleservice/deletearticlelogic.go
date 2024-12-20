package articleservicelogic

import (
	"coderhub/conf"
	"coderhub/model"
	"coderhub/shared/utils"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteArticleLogic {
	return &DeleteArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteArticleLogic) DeleteArticle(in *coderhub.DeleteArticleRequest) (*coderhub.DeleteArticleResponse, error) {
	// 使用 defer 记录方法执行时间
	defer func(start time.Time) {
		l.Logger.Infof("DeleteArticle 执行耗时: %v", time.Since(start))
	}(time.Now())

	// 从 metadata 中获取 userId
	var (
		userId string
		err    error
	)
	if userId, err = utils.GetUserMetaData(l.ctx); err != nil {
		return nil, err
	}

	if userId != strconv.FormatInt(in.UserId, 10) {
		return nil, fmt.Errorf("非法操作")
	}

	// 参数校验
	if err := l.validateRequest(in); err != nil {
		return nil, err
	}

	// 检查文章是否存在
	article, err := l.getArticle(in.Id)
	if err != nil {
		return nil, err
	}

	// 记录操作日志
	l.Logger.Infof("开始删除文章, ID: %d, 标题: %s", article.ID, article.Title)

	// 权限校验
	if article.AuthorID != in.UserId {
		l.Logger.Errorf(conf.ErrValidationFailed)
		return nil, errors.New(conf.ErrValidationFailed)
	}

	// 执行软删除
	if err := l.deleteArticle(in.Id); err != nil {
		return nil, err
	}

	return &coderhub.DeleteArticleResponse{
		Success: true,
	}, nil
}

// validateRequest 验证删除请求的参数
func (l *DeleteArticleLogic) validateRequest(in *coderhub.DeleteArticleRequest) error {
	if err := utils.NewValidator().ArticleID(in.Id).Check(); err != nil {
		l.Logger.Errorf(conf.ErrValidationFailed, err)
		return fmt.Errorf(conf.ErrValidationFailed, err)
	}
	return nil
}

// getArticle 获取并验证文章是否存在
func (l *DeleteArticleLogic) getArticle(articleID int64) (*model.Articles, error) {
	article, err := l.svcCtx.ArticleRepository.GetArticleByID(articleID)
	if err != nil {
		l.Logger.Errorf(conf.ErrGetArticleFailed, err)
		return nil, fmt.Errorf(conf.ErrGetArticleFailed, err)
	}

	if article == nil {
		l.Logger.Errorf(conf.ErrArticleNotFound)
		return nil, errors.New(conf.ErrArticleNotFound)
	}

	return article, nil
}

// deleteArticle 执行文章删除操作
func (l *DeleteArticleLogic) deleteArticle(articleID int64) error {
	if err := l.svcCtx.ArticleRepository.DeleteArticle(articleID); err != nil {
		l.Logger.Errorf(conf.ErrDeleteArticleFailed, err)
		return fmt.Errorf(conf.ErrDeleteArticleFailed, err)
	}
	return nil
}
