package svc

import (
	"coderhub/api/coderhub/internal/config"
	"coderhub/rpc/coderhub/coderhub"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config              config.Config
	UserService         coderhub.UserServiceClient
	ImageAuthService    coderhub.ImageServiceClient
	ArticlesService     coderhub.ArticleServiceClient
	AcademicService     coderhub.AcademicNavigatorServiceClient
	UserFollowService   coderhub.UserFollowServiceClient
	ImagesService       coderhub.ImageServiceClient
	CommentService      coderhub.CommentServiceClient
	QuestionBankService coderhub.QuestionServiceClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:              c,
		UserService:         coderhub.NewUserServiceClient(zrpc.MustNewClient(c.UserService).Conn()),
		ImageAuthService:    coderhub.NewImageServiceClient(zrpc.MustNewClient(c.ImageAuthService).Conn()),
		ArticlesService:     coderhub.NewArticleServiceClient(zrpc.MustNewClient(c.ArticlesService).Conn()),
		AcademicService:     coderhub.NewAcademicNavigatorServiceClient(zrpc.MustNewClient(c.AcademicService).Conn()),
		UserFollowService:   coderhub.NewUserFollowServiceClient(zrpc.MustNewClient(c.UserFollowService).Conn()),
		ImagesService:       coderhub.NewImageServiceClient(zrpc.MustNewClient(c.ImagesService).Conn()),
		CommentService:      coderhub.NewCommentServiceClient(zrpc.MustNewClient(c.CommentService).Conn()),
		QuestionBankService: coderhub.NewQuestionServiceClient(zrpc.MustNewClient(c.QuestionBankService).Conn()),
	}
}
